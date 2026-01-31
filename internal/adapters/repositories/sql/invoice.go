package repositoriesSql

import (
	"context"
	"fmt"
	"time"

	"app/xonvera-core/internal/adapters/dto"
	"app/xonvera-core/internal/core/domain"
	portRepository "app/xonvera-core/internal/core/ports/repository"

	"gorm.io/gorm"
)

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) portRepository.InvoiceRepository {
	return &invoiceRepository{db: db}
}

// GenerateInvoiceID generates an invoice ID with format: 2YYYYMMDDUUUUUSSSS
// 2 => invoice id prefix
// YYYYMMDD => year, month, date
// UUUUU => user id (5 digits, zero-padded)
// SSSS => per-user daily suffix (1-9999)
func (r *invoiceRepository) GenerateInvoiceID(ctx context.Context, tx portRepository.Transaction, userID uint, date time.Time) (int64, error) {

	// Format: 2 + YYYYMMDD + UUUUU + SSSS
	prefix := "2"
	dateStr := date.Format("20060102") // YYYYMMDD
	dateOnly := date.Format("2006-01-02")

	// Get per-user daily sequence value atomically
	var suffix int64
	err := txDb(tx, r.db).WithContext(ctx).Raw(
		`INSERT INTO app.invoice_user_daily_seq (user_id, day, counter)
		 VALUES (?, ?::date, 1)
		 ON CONFLICT (user_id, day)
		 DO UPDATE SET counter = app.invoice_user_daily_seq.counter + 1
		 RETURNING counter`,
		userID, dateOnly,
	).Scan(&suffix).Error
	if err != nil {
		return 0, err
	}
	if suffix > 9999 {
		return 0, fmt.Errorf("invoice suffix overflow for user %d on %s", userID, dateStr)
	}

	// Format the complete invoice ID
	invoiceIDStr := fmt.Sprintf("%s%s%04d", prefix, dateStr, suffix)

	var invoiceID int64
	if _, err := fmt.Sscanf(invoiceIDStr, "%d", &invoiceID); err != nil {
		return 0, err
	}

	return invoiceID, nil
}

func (r *invoiceRepository) Get(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	query := r.db.WithContext(ctx).Order("created_at DESC")

	if req.Limit > 0 {
		query = query.Limit(int(req.Limit))
	}
	if req.Offset > 0 {
		query = query.Offset(int(req.Offset))
	}

	query = query.Model(&domain.Invoice{}).Where("author_id = ?", req.UserID)

	//get count
	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return nil, err
	}

	var resp dto.PaginationResponse
	resp.Meta = dto.PaginationMetaResponse{
		Page:      req.Page,
		Limit:     req.Limit,
		TotalData: uint64(count),
		TotalPage: GetTotalPage(count, req.Limit),
	}

	query = query.Limit(int(req.Limit)).Offset(int(req.Offset))

	var data []domain.Invoice
	err = query.Scan(&data).Error
	if err != nil {
		return nil, err
	}

	resp.Data = make([]any, len(data))
	for i, v := range data {
		resp.Data[i] = v.Response()
	}

	return &resp, nil
}

func (r *invoiceRepository) GetByID(ctx context.Context, id int64) (*domain.Invoice, error) {
	var invoice domain.Invoice
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&invoice).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("404:not found invoice")
		}
		return nil, err
	}
	return &invoice, nil
}


func (r *invoiceRepository) GetItems(ctx context.Context, invoiceID []int64) ([]domain.InvoiceItem, error) {
	var items []domain.InvoiceItem
	err := r.db.WithContext(ctx).Where("invoice_id IN ?", invoiceID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *invoiceRepository) Create(ctx context.Context, tx portRepository.Transaction, data *domain.Invoice) error {
	return txDb(tx, r.db).WithContext(ctx).Create(data).Error
}

func (r *invoiceRepository) Update(ctx context.Context, tx portRepository.Transaction, data *domain.Invoice) error {
	updates := map[string]interface{}{
		"issuer":     data.Issuer,
		"customer":   data.Customer,
		"issue_date": data.IssueDate,
		"due_date":   data.DueDate,
		"note":       data.Note,
		"updated_at": data.UpdatedAt,
	}

	return txDb(tx, r.db).
		WithContext(ctx).
		Model(&domain.Invoice{}).
		Where("id = ? AND author_id = ?", data.ID, data.AuthorID).
		Updates(updates).
		Error
}

func (r *invoiceRepository) CreateItem(ctx context.Context, tx portRepository.Transaction, data []domain.InvoiceItem) error {
	return txDb(tx, r.db).WithContext(ctx).CreateInBatches(data, 100).Error
}

func (r *invoiceRepository) DeleteItemsByInvoiceID(ctx context.Context, tx portRepository.Transaction, invoiceID int64) error {
	return txDb(tx, r.db).WithContext(ctx).Where("invoice_id = ?", invoiceID).Delete(&domain.InvoiceItem{}).Error
}
