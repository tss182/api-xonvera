package repositoriesSql

import (
	"context"
	"fmt"
	"time"

	"app/xonvera-core/internal/core/domain"
	"app/xonvera-core/internal/core/ports/repository"

	"gorm.io/gorm"
)

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) portRepository.InvoiceRepository {
	return &invoiceRepository{db: db}
}

// GenerateInvoiceID generates an invoice ID with format: 2YYYYMMDD00001
// 2 => invoice id prefix
// YYYYMMDD => year, month, date
// 00001 => suffix auto increment
// TODO: Implement daily sequence reset to ensure suffix stays within 00001-99999
func (r *invoiceRepository) GenerateInvoiceID(ctx context.Context) (int64, error) {
	now := time.Now()
	
	// Format: 2 + YYYYMMDD + SSSSS
	prefix := "2"
	dateStr := now.Format("20060102") // YYYYMMDD
	
	// Get next sequence value
	var suffix int
	err := r.db.WithContext(ctx).Raw("SELECT nextval('billing.invoice_suffix_seq')").Scan(&suffix).Error
	if err != nil {
		return 0, err
	}
	
	// Format the complete invoice ID
	invoiceIDStr := fmt.Sprintf("%s%s%05d", prefix, dateStr, suffix)
	
	var invoiceID int64
	fmt.Sscanf(invoiceIDStr, "%d", &invoiceID)
	
	return invoiceID, nil
}

func (r *invoiceRepository) Create(ctx context.Context, invoice *domain.Invoice, items []domain.InvoiceItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create invoice
		if err := tx.Create(invoice).Error; err != nil {
			return err
		}
		
		// Create invoice items if any
		if len(items) > 0 {
			for i := range items {
				items[i].InvoiceID = invoice.ID
			}
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}

func (r *invoiceRepository) GetByID(ctx context.Context, id int64) (*domain.Invoice, error) {
	var invoice domain.Invoice
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Invoice, error) {
	var invoices []domain.Invoice
	query := r.db.WithContext(ctx).Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Find(&invoices).Error
	if err != nil {
		return nil, err
	}
	return invoices, nil
}

func (r *invoiceRepository) GetItems(ctx context.Context, invoiceID int64) ([]domain.InvoiceItem, error) {
	var items []domain.InvoiceItem
	err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
