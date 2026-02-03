package portService

import (
	"context"
	"io"

	"app/xonvera-core/internal/adapters/dto"
)

// PDFService defines the contract for PDF generation and management
type PDFService interface {
	// GenerateInvoicePDF generates a PDF from invoice data and returns the PDF bytes
	GenerateInvoicePDF(ctx context.Context, invoice *dto.InvoiceResponse, items []dto.InvoiceItemDTO) ([]byte, error)

	// GetInvoicePDF retrieves an existing invoice PDF by invoice ID
	GetInvoicePDF(ctx context.Context, invoiceID int64) ([]byte, error)

	// GetInvoicePDFStream retrieves an existing invoice PDF as a stream
	GetInvoicePDFStream(ctx context.Context, invoiceID int64) (io.ReadCloser, error)

	// SaveInvoicePDF saves PDF bytes to the filesystem and returns the file path
	SaveInvoicePDF(ctx context.Context, invoiceID int64, pdfData []byte) (string, error)

	// PDFExists checks if a PDF file exists for the given invoice ID
	PDFExists(ctx context.Context, invoiceID int64) bool

	// DeleteInvoicePDF removes a PDF file for the given invoice ID
	DeleteInvoicePDF(ctx context.Context, invoiceID int64) error
}
