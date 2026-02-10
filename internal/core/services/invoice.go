package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"app/xonvera-core/internal/adapters/dto"
	"app/xonvera-core/internal/core/domain"
	portRepository "app/xonvera-core/internal/core/ports/repository"
	portService "app/xonvera-core/internal/core/ports/service"
	"app/xonvera-core/internal/infrastructure/logger"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
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

func (s *invoiceService) Get(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginationResponse, error) {
	res, err := s.repo.Get(ctx, req)
	if err != nil {
		logger.StdContextError(ctx, "failed to get all invoices", zap.Error(err))
		return nil, err
	}
	return res, nil
}

// GetByID retrieves a single invoice with its items by invoice ID
func (s *invoiceService) GetByID(ctx context.Context, invoiceID int64, userID uint) (*dto.InvoiceResponse, error) {
	invoice, err := s.repo.GetByID(ctx, invoiceID)
	if err != nil {
		logger.StdContextError(ctx, "failed to get invoice by ID", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, err
	}

	// Verify invoice belongs to user
	if invoice.AuthorID != userID {
		logger.StdContextWarn(ctx, "unauthorized invoice access", zap.Int64("invoice_id", invoiceID), zap.Uint("user_id", userID))
		return nil, fmt.Errorf("404:not found invoice")
	}

	// Fetch invoice items
	items, err := s.repo.GetItemsByInvoiceID(ctx, invoiceID)
	if err != nil {
		logger.StdContextError(ctx, "failed to get invoice items", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, err
	}

	// Convert to response DTOs
	itemDTOs := make([]dto.InvoiceItemResponse, len(items))
	for i, item := range items {
		itemDTOs[i] = dto.InvoiceItemResponse{
			ID:          item.ID,
			InvoiceID:   item.InvoiceID,
			Description: item.Description,
			Qty:         item.Qty,
			Price:       item.Price,
			Total:       item.Total,
			CreatedAt:   item.CreatedAt.Format(time.DateTime),
		}
	}

	response := &dto.InvoiceResponse{
		ID:        invoice.ID,
		Customer:  invoice.Customer,
		Issuer:    invoice.Issuer,
		IssueDate: invoice.IssueDate,
		DueDate:   invoice.DueDate,
		Note:      invoice.Note,
		Status:    invoice.Status,
		Items:     itemDTOs,
		CreatedAt: invoice.CreatedAt,
		UpdatedAt: invoice.UpdatedAt,
	}

	return response, nil
}

func (s *invoiceService) Create(ctx context.Context, req *dto.InvoiceRequest) error {
	tx, err := s.tx.Begin()
	if err != nil {
		logger.StdContextError(ctx, "failed to begin transaction", zap.Error(err))
		return err
	}

	issueDate, err := time.ParseInLocation("2006-01-02", req.IssueDate, time.Local)
	if err != nil {
		logger.StdContextError(ctx, "failed to parse issue date", zap.Error(err))
		return err
	}

	dueDate, err := time.ParseInLocation("2006-01-02 15:04:05", req.DueDate, time.Local)
	if err != nil {
		logger.StdContextError(ctx, "failed to parse due date", zap.Error(err))
		return err
	}

	// Ensure rollback on error, commit will override this
	defer tx.Rollback()

	// Generate invoice ID
	invoiceID, err := s.repo.GenerateInvoiceID(ctx, tx, req.UserID, issueDate)
	if err != nil {
		tx.Rollback()
		logger.StdContextError(ctx, "failed to generate invoice ID", zap.Error(err))
		return err
	}

	t := time.Now()

	data := domain.Invoice{
		ID:        invoiceID,
		Issuer:    req.Issuer,
		Customer:  req.Customer,
		IssueDate: issueDate.Format(time.DateOnly),
		DueDate:   dueDate,
		Note:      req.Note,
		AuthorID:  req.UserID,
		Status:    "unpaid",
		Timestamp: domain.Timestamp{CreatedAt: t, UpdatedAt: t},
	}

	// Create invoice
	if err = s.repo.Create(ctx, tx, &data); err != nil {
		tx.Rollback()
		logger.StdContextError(ctx, "failed to create invoice", zap.Error(err))
		return err
	}

	// Create invoice items
	items := make([]domain.InvoiceItem, 0, len(req.Items))
	for i, v := range req.Items {
		total := v.Qty * v.Price
		item := domain.InvoiceItem{
			ID:          uint(i + 1),
			InvoiceID:   invoiceID,
			Description: v.Description,
			Qty:         v.Qty,
			Price:       v.Price,
			Total:       total,
			Timestamp:   domain.Timestamp{CreatedAt: t, UpdatedAt: t},
		}
		items = append(items, item)
	}

	if err = s.repo.CreateItem(ctx, tx, items); err != nil {
		tx.Rollback()
		logger.StdContextError(ctx, "failed to create invoice items", zap.Error(err))
		return err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		logger.StdContextError(ctx, "failed to commit transaction", zap.Error(err))
		return err
	}

	logger.StdContextInfo(ctx, "invoice created successfully", zap.Int64("invoice_id", invoiceID))
	return nil
}

func (s *invoiceService) Update(ctx context.Context, req *dto.InvoiceRequest) error {
	tx, err := s.tx.Begin()
	if err != nil {
		logger.StdContextError(ctx, "failed to begin transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback()

	// Ensure invoice exists and belongs to user
	inv, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if inv.AuthorID != req.UserID {
		return fmt.Errorf("404:not found invoice")
	}

	issueDate, err := time.Parse(time.DateOnly, req.IssueDate)
	if err != nil {
		logger.StdContextError(ctx, "failed to parse issue date", zap.Error(err))
		return err
	}

	dueDate, err := time.Parse("2006-01-02 15:04:05", req.DueDate)
	if err != nil {
		logger.StdContextError(ctx, "failed to parse due date", zap.Error(err))
		return err
	}

	updatedAt := time.Now()

	data := domain.Invoice{
		ID:        req.ID,
		Issuer:    req.Issuer,
		Customer:  req.Customer,
		IssueDate: issueDate.Format("2006-01-02"),
		DueDate:   dueDate,
		Note:      req.Note,
		AuthorID:  req.UserID,
		Timestamp: domain.Timestamp{UpdatedAt: updatedAt},
	}

	if err = s.repo.Update(ctx, tx, &data); err != nil {
		logger.StdContextError(ctx, "failed to update invoice", zap.Error(err))
		return err
	}

	if err = s.repo.DeleteItemsByInvoiceID(ctx, tx, req.ID); err != nil {
		logger.StdContextError(ctx, "failed to delete invoice items", zap.Error(err))
		return err
	}

	items := make([]domain.InvoiceItem, 0, len(req.Items))
	for i, v := range req.Items {
		total := v.Qty * v.Price
		item := domain.InvoiceItem{
			ID:          uint(i + 1),
			InvoiceID:   req.ID,
			Description: v.Description,
			Qty:         v.Qty,
			Price:       v.Price,
			Total:       total,
			Timestamp:   domain.Timestamp{CreatedAt: updatedAt, UpdatedAt: updatedAt},
		}
		items = append(items, item)
	}

	if err = s.repo.CreateItem(ctx, tx, items); err != nil {
		logger.StdContextError(ctx, "failed to create invoice items", zap.Error(err))
		return err
	}

	if err = tx.Commit(); err != nil {
		logger.StdContextError(ctx, "failed to commit transaction", zap.Error(err))
		return err
	}

	//remove pdf file if exists, so it will be regenerated on next request
	filePdf := fmt.Sprintf("assets/pdf/invoice_%d.pdf", req.ID)
	if _, err := os.Stat(filePdf); err == nil {
		err = os.Remove(filePdf)
		if err != nil {
			logger.StdContextError(ctx, "failed to remove existing pdf file", zap.Error(err), zap.Int64("invoice_id", req.ID))
		} else {
			logger.StdContextInfo(ctx, "existing pdf file removed", zap.Int64("invoice_id", req.ID))
		}
	}

	logger.StdContextInfo(ctx, "invoice updated successfully", zap.Int64("invoice_id", req.ID))
	return nil
}

func (s *invoiceService) GetPDF(ctx context.Context, invoiceID int64, userID uint) ([]byte, error) {
	// Ensure invoice exists and belongs to user
	inv, err := s.repo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	if inv.AuthorID != userID {
		return nil, fmt.Errorf("404:not found invoice")
	}

	filePdf := fmt.Sprintf("assets/pdf/invoice_%d.pdf", invoiceID)

	// Check if PDF already exists
	if _, err := os.Stat(filePdf); err == nil {
		// PDF exists, read it and return
		pdfBytes, err := os.ReadFile(filePdf)
		if err != nil {
			logger.StdContextError(ctx, "failed to read existing pdf", zap.Error(err), zap.Int64("invoice_id", invoiceID))
			return nil, err
		}
		return pdfBytes, nil
	}

	// PDF doesn't exist, generate it
	// Fetch invoice items
	items, err := s.repo.GetItemsByInvoiceID(ctx, invoiceID)
	if err != nil {
		logger.StdContextError(ctx, "failed to get invoice items", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, err
	}

	fmt.Println("items", items)

	m := s.generatePDF()
	doc, err := m.Generate()
	if err != nil {
		logger.StdContextError(ctx, "failed to generate pdf", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, err
	}

	// Get PDF as binary data
	pdfBytes := doc.GetBytes()

	// Save pdf to file for future use
	err = doc.Save(filePdf)
	if err != nil {
		logger.StdContextError(ctx, "failed to save pdf to file", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return nil, err
	}

	logger.StdContextInfo(ctx, "pdf generated successfully", zap.Int64("invoice_id", invoiceID), zap.Int("size_bytes", len(pdfBytes)))
	return pdfBytes, nil
}

func (s *invoiceService) generatePDF() core.Maroto {
	var intro = "This is an example of how to create a PDF using Maroto in Go. " +
		"Maroto is a simple and powerful library that allows you to create PDF documents " +
		"with ease. You can add text, images, tables, and more to your PDF files."
	cfg := config.NewBuilder().
		WithDebug(true).
		Build()

	m := maroto.New(cfg)

	m.AddAutoRow(
		image.NewFromFileCol(5, "docs/assets/images/biplane.jpg"),
		text.NewCol(7, intro),
	)

	m.AddAutoRow(
		image.NewFromFileCol(5, "docs/assets/images/biplane.jpg"),
		text.NewCol(7, intro, props.Text{
			Size: 13,
		}),
	)

	m.AddAutoRow(
		image.NewFromFileCol(5, "docs/assets/images/biplane.jpg"),
		text.NewCol(7, intro, props.Text{
			Size:   13,
			Top:    8,
			Bottom: 9,
		}),
	)

	m.AddAutoRow(
		code.NewBarCol(4, "code"),
		text.NewCol(8, intro),
	)

	m.AddAutoRow(
		code.NewMatrixCol(3, "code"),
		text.NewCol(9, intro),
	)

	m.AddAutoRow(
		code.NewQrCol(2, "code"),
		text.NewCol(10, intro),
	)

	return m
}
