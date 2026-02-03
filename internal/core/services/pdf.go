package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"app/xonvera-core/internal/adapters/dto"
	portService "app/xonvera-core/internal/core/ports/service"
	"app/xonvera-core/internal/infrastructure/logger"

	"go.uber.org/zap"
)

type pdfService struct {
	pdfDir string
}

// NewPDFService creates a new PDF service instance
// pdfDir should be the path where PDFs will be stored (e.g., "assets/pdf")
func NewPDFService(pdfDir string) (portService.PDFService, error) {
	// Create the PDF directory if it doesn't exist
	if err := os.MkdirAll(pdfDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create PDF directory: %w", err)
	}

	return &pdfService{
		pdfDir: pdfDir,
	}, nil
}

// getPDFFileName generates a unique filename for an invoice PDF
func (s *pdfService) getPDFFileName(invoiceID int64) string {
	return fmt.Sprintf("invoice_%d.pdf", invoiceID)
}

// getPDFPath returns the full path for an invoice PDF
func (s *pdfService) getPDFPath(invoiceID int64) string {
	return filepath.Join(s.pdfDir, s.getPDFFileName(invoiceID))
}

// GenerateInvoicePDF generates a PDF from invoice data and returns the PDF bytes
// This uses a simple text-based PDF format for demonstration
func (s *pdfService) GenerateInvoicePDF(ctx context.Context, invoice *dto.InvoiceResponse, items []dto.InvoiceItemDTO) ([]byte, error) {
	if invoice == nil {
		return nil, fmt.Errorf("invoice data is required")
	}

	// For production, consider using github.com/jung-kurt/gofpdf or similar
	// This is a simplified implementation
	pdfContent := s.generateSimplePDF(invoice, items)
	return pdfContent, nil
}

// generateSimplePDF generates a simple PDF content as bytes
func (s *pdfService) generateSimplePDF(invoice *dto.InvoiceResponse, items []dto.InvoiceItemDTO) []byte {
	var buffer bytes.Buffer

	// PDF Header
	buffer.WriteString("%PDF-1.4\n")

	// PDF Objects
	buffer.WriteString("1 0 obj\n")
	buffer.WriteString("<< /Type /Catalog /Pages 2 0 R >>\n")
	buffer.WriteString("endobj\n")

	buffer.WriteString("2 0 obj\n")
	buffer.WriteString("<< /Type /Pages /Kids [3 0 R] /Count 1 >>\n")
	buffer.WriteString("endobj\n")

	// Content stream
	contentStream := s.buildContentStream(invoice, items)
	contentLength := len(contentStream)

	buffer.WriteString("3 0 obj\n")
	buffer.WriteString("<< /Type /Page /Parent 2 0 R /Resources << /Font << /F1 4 0 R >> >> /MediaBox [0 0 612 792] /Contents 5 0 R >>\n")
	buffer.WriteString("endobj\n")

	buffer.WriteString("4 0 obj\n")
	buffer.WriteString("<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\n")
	buffer.WriteString("endobj\n")

	buffer.WriteString("5 0 obj\n")
	buffer.WriteString(fmt.Sprintf("<< /Length %d >>\n", contentLength))
	buffer.WriteString("stream\n")
	buffer.WriteString(contentStream)
	buffer.WriteString("\nendstream\n")
	buffer.WriteString("endobj\n")

	// xref table
	xrefOffset := buffer.Len()
	buffer.WriteString("xref\n")
	buffer.WriteString("0 6\n")
	buffer.WriteString("0000000000 65535 f \n")
	buffer.WriteString("0000000009 00000 n \n")
	buffer.WriteString("0000000058 00000 n \n")
	buffer.WriteString("0000000115 00000 n \n")
	buffer.WriteString("0000000273 00000 n \n")
	buffer.WriteString(fmt.Sprintf("%010d 00000 n \n", xrefOffset-contentLength-50))

	// Trailer
	buffer.WriteString("trailer\n")
	buffer.WriteString("<< /Size 6 /Root 1 0 R >>\n")
	buffer.WriteString("startxref\n")
	buffer.WriteString(fmt.Sprintf("%d\n", xrefOffset))
	buffer.WriteString("%%EOF\n")

	return buffer.Bytes()
}

// buildContentStream builds the text content for the PDF
func (s *pdfService) buildContentStream(invoice *dto.InvoiceResponse, items []dto.InvoiceItemDTO) string {
	var content bytes.Buffer

	// PDF text positioning commands
	content.WriteString("BT\n")
	content.WriteString("/F1 14 Tf\n")
	content.WriteString("50 750 Td\n")
	content.WriteString("(INVOICE) Tj\n")
	content.WriteString("ET\n")

	// Invoice details
	yPos := 700.0
	lineHeight := 20.0

	details := []string{
		fmt.Sprintf("Invoice ID: %d", invoice.ID),
		fmt.Sprintf("Issue Date: %s", invoice.IssueDate),
		fmt.Sprintf("Due Date: %s", invoice.DueDate.Format(time.DateOnly)),
		fmt.Sprintf("From: %s", invoice.Issuer),
		fmt.Sprintf("To: %s", invoice.Customer),
	}

	for _, detail := range details {
		content.WriteString("BT\n")
		content.WriteString("/F1 10 Tf\n")
		content.WriteString(fmt.Sprintf("50 %.0f Td\n", yPos))
		content.WriteString(fmt.Sprintf("(%s) Tj\n", detail))
		content.WriteString("ET\n")
		yPos -= lineHeight
	}

	// Items table header
	yPos -= 20
	content.WriteString("BT\n")
	content.WriteString("/F1 10 Tf\n")
	content.WriteString(fmt.Sprintf("50 %.0f Td\n", yPos))
	content.WriteString("(Description | Qty | Price | Total) Tj\n")
	content.WriteString("ET\n")

	// Items
	yPos -= lineHeight
	totalAmount := 0
	for _, item := range items {
		itemLine := fmt.Sprintf("%s | %d | %d | %d", item.Description, item.Qty, item.Price, item.Total)
		content.WriteString("BT\n")
		content.WriteString("/F1 9 Tf\n")
		content.WriteString(fmt.Sprintf("50 %.0f Td\n", yPos))
		content.WriteString(fmt.Sprintf("(%s) Tj\n", itemLine))
		content.WriteString("ET\n")
		yPos -= lineHeight
		totalAmount += item.Total
	}

	// Total
	yPos -= 10
	totalLine := fmt.Sprintf("Total Amount: %d", totalAmount)
	content.WriteString("BT\n")
	content.WriteString("/F1 12 Tf\n")
	content.WriteString(fmt.Sprintf("50 %.0f Td\n", yPos))
	content.WriteString(fmt.Sprintf("(%s) Tj\n", totalLine))
	content.WriteString("ET\n")

	// Note
	if invoice.Note != "" {
		yPos -= 20
		content.WriteString("BT\n")
		content.WriteString("/F1 9 Tf\n")
		content.WriteString(fmt.Sprintf("50 %.0f Td\n", yPos))
		content.WriteString(fmt.Sprintf("(Note: %s) Tj\n", invoice.Note))
		content.WriteString("ET\n")
	}

	// Generated timestamp
	content.WriteString("BT\n")
	content.WriteString("/F1 8 Tf\n")
	content.WriteString("50 20 Td\n")
	content.WriteString(fmt.Sprintf("(Generated: %s) Tj\n", time.Now().Format(time.DateTime)))
	content.WriteString("ET\n")

	return content.String()
}

// SaveInvoicePDF saves PDF bytes to the filesystem
func (s *pdfService) SaveInvoicePDF(ctx context.Context, invoiceID int64, pdfData []byte) (string, error) {
	if len(pdfData) == 0 {
		return "", fmt.Errorf("PDF data is empty")
	}

	pdfPath := s.getPDFPath(invoiceID)

	if err := os.WriteFile(pdfPath, pdfData, 0644); err != nil {
		logger.StdContextError(ctx, "failed to save invoice PDF", zap.Error(err), zap.String("path", pdfPath))
		return "", fmt.Errorf("failed to save invoice PDF: %w", err)
	}

	logger.StdContextInfo(ctx, "invoice PDF saved successfully", zap.String("path", pdfPath), zap.Int64("invoice_id", invoiceID))
	return pdfPath, nil
}

// GetInvoicePDF retrieves an existing invoice PDF by invoice ID
func (s *pdfService) GetInvoicePDF(ctx context.Context, invoiceID int64) ([]byte, error) {
	pdfPath := s.getPDFPath(invoiceID)

	pdfData, err := os.ReadFile(pdfPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("invoice PDF not found: %d", invoiceID)
		}
		logger.StdContextError(ctx, "failed to read invoice PDF", zap.Error(err), zap.String("path", pdfPath))
		return nil, fmt.Errorf("failed to read invoice PDF: %w", err)
	}

	return pdfData, nil
}

// GetInvoicePDFStream retrieves an existing invoice PDF as a stream
func (s *pdfService) GetInvoicePDFStream(ctx context.Context, invoiceID int64) (io.ReadCloser, error) {
	pdfPath := s.getPDFPath(invoiceID)

	file, err := os.Open(pdfPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("invoice PDF not found: %d", invoiceID)
		}
		logger.StdContextError(ctx, "failed to open invoice PDF", zap.Error(err), zap.String("path", pdfPath))
		return nil, fmt.Errorf("failed to open invoice PDF: %w", err)
	}

	return file, nil
}

// PDFExists checks if a PDF file exists for the given invoice ID
func (s *pdfService) PDFExists(ctx context.Context, invoiceID int64) bool {
	pdfPath := s.getPDFPath(invoiceID)
	_, err := os.Stat(pdfPath)
	return err == nil
}

// DeleteInvoicePDF removes a PDF file for the given invoice ID
func (s *pdfService) DeleteInvoicePDF(ctx context.Context, invoiceID int64) error {
	pdfPath := s.getPDFPath(invoiceID)

	if err := os.Remove(pdfPath); err != nil {
		if !os.IsNotExist(err) {
			logger.StdContextError(ctx, "failed to delete invoice PDF", zap.Error(err), zap.String("path", pdfPath))
			return fmt.Errorf("failed to delete invoice PDF: %w", err)
		}
		// File doesn't exist, which is fine
	}

	return nil
}

// Verify that pdfService implements PDFService interface
var _ portService.PDFService = (*pdfService)(nil)
