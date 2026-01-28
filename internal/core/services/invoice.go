package services

import (
	"context"
	"errors"

	"app/xonvera-core/internal/core/domain"
	portRepository "app/xonvera-core/internal/core/ports/repository"
	portService "app/xonvera-core/internal/core/ports/service"
	"app/xonvera-core/internal/infrastructure/logger"

	"go.uber.org/zap"
)

type invoiceService struct {
	invoiceRepo portRepository.InvoiceRepository
}

func NewInvoiceService(
	invoiceRepo portRepository.InvoiceRepository,
) portService.InvoiceService {
	return &invoiceService{
		invoiceRepo: invoiceRepo,
	}
}

func (s *invoiceService) CreateInvoice(ctx context.Context, invoice *domain.Invoice, items []domain.InvoiceItem) error {
	// Generate invoice ID
	invoiceID, err := s.invoiceRepo.GenerateInvoiceID(ctx)
	if err != nil {
		logger.Error("failed to generate invoice ID", zap.Error(err))
		return err
	}
	
	invoice.ID = invoiceID
	
	// Create invoice with items
	if err := s.invoiceRepo.Create(ctx, invoice, items); err != nil {
		logger.Error("failed to create invoice", zap.Error(err))
		return err
	}
	
	return nil
}

func (s *invoiceService) GetInvoiceByID(ctx context.Context, id int64) (*domain.Invoice, error) {
	invoice, err := s.invoiceRepo.GetByID(ctx, id)
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
	
	invoices, err := s.invoiceRepo.GetAll(ctx, limit, offset)
	if err != nil {
		logger.Error("failed to get all invoices", zap.Error(err))
		return nil, err
	}
	return invoices, nil
}

func (s *invoiceService) GetInvoiceItems(ctx context.Context, invoiceID int64) ([]domain.InvoiceItem, error) {
	items, err := s.invoiceRepo.GetItems(ctx, invoiceID)
	if err != nil {
		logger.Error("failed to get invoice items", zap.Int64("invoice_id", invoiceID), zap.Error(err))
		return nil, err
	}
	return items, nil
}
