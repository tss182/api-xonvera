package portService

import (
	"context"

	"app/xonvera-core/internal/adapters/dto"
	"app/xonvera-core/internal/core/domain"
)

type InvoiceService interface {
	Create(ctx context.Context, req dto.CreateInvoiceRequest) error
	GetInvoiceByID(ctx context.Context, id int64) (*domain.Invoice, error)
	GetAllInvoices(ctx context.Context, limit, offset int) ([]domain.Invoice, error)
	GetInvoiceItems(ctx context.Context, invoiceID int64) ([]domain.InvoiceItem, error)
}
