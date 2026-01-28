package portRepository

import (
	"context"

	"app/xonvera-core/internal/core/domain"
)

type InvoiceRepository interface {
	Create(ctx context.Context, invoice *domain.Invoice, items []domain.InvoiceItem) error
	GetByID(ctx context.Context, id int64) (*domain.Invoice, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Invoice, error)
	GetItems(ctx context.Context, invoiceID int64) ([]domain.InvoiceItem, error)
	GenerateInvoiceID(ctx context.Context) (int64, error)
}
