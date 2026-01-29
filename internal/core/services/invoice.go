package services

import (
	"context"
	"errors"
	"time"

	"app/xonvera-core/internal/adapters/dto"
	"app/xonvera-core/internal/core/domain"
	portRepository "app/xonvera-core/internal/core/ports/repository"
	portService "app/xonvera-core/internal/core/ports/service"
	"app/xonvera-core/internal/infrastructure/logger"

	"go.uber.org/zap"
)

type invoiceService struct {
	repo portRepository.InvoiceRepository
	tx   portRepository.TxRepository
}

func NewInvoiceService(invoiceRepo portRepository.InvoiceRepository, tx portRepository.TxRepository) portService.InvoiceService {
	return &invoiceService{
		repo: invoiceRepo,
		tx:   tx,
	}
}

func (s *invoiceService) Create(ctx context.Context, req dto.CreateInvoiceRequest) error {
	tx, err := s.tx.Begin()
	if err != nil {
		logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback()

	// Generate invoice ID
	invoiceID, err := s.repo.GenerateInvoiceID(ctx, tx)
	if err != nil {
		logger.Error("failed to generate invoice ID", zap.Error(err))
		return err
	}

	t := time.Now()

	data := domain.Invoice{
		ID:        invoiceID,
		Issuer:    req.Issuer,
		Customer:  req.Customer,
		IssueDate: req.IssueDate,
		Note:      req.Note,
		CreatedAt: t,
		UpdatedAt: t,
	}

	// Create invoice
	err = s.repo.Create(ctx, tx, &data)
	if err != nil {
		logger.Error("failed to create invoice", zap.Error(err))
		return err
	}

	//create invoice items
	var items = make([]domain.InvoiceItem, 0, len(req.Items))
	for _, v := range req.Items {
		total := v.Qty * v.Price
		item := domain.InvoiceItem{
			InvoiceID:   invoiceID,
			Description: v.Description,
			Qty:         v.Qty,
			Price:       v.Price,
			Total:       total,
			CreatedAt:   t,
			UpdatedAt:   t,
		}
		items = append(items, item)
	}
	err = s.repo.CreateItem(ctx, tx, items)
	if err != nil {
		logger.Error("failed to create invoice items", zap.Error(err))
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *invoiceService) GetInvoiceByID(ctx context.Context, id int64) (*domain.Invoice, error) {
	invoice, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("failed to get invoice by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return invoice, nil
}

func (s *invoiceService) GetAllInvoices(ctx context.Context, limit, offset int) ([]domain.Invoice, error) {
	if limit <= 0 {
		limit = 20 // default limit
	}
	if limit > 100 {
		return nil, errors.New("limit cannot exceed 100")
	}

	invoices, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		logger.Error("failed to get all invoices", zap.Error(err))
		return nil, err
	}
	return invoices, nil
}

func (s *invoiceService) GetInvoiceItems(ctx context.Context, invoiceID int64) ([]domain.InvoiceItem, error) {
	items, err := s.repo.GetItems(ctx, invoiceID)
	if err != nil {
		logger.Error("failed to get invoice items", zap.Int64("invoice_id", invoiceID), zap.Error(err))
		return nil, err
	}
	return items, nil
}
