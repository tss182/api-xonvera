package portRepository

import (
	"context"
	"time"

	"app/xonvera-core/internal/adapters/dto"
	"app/xonvera-core/internal/core/domain"
)

type InvoiceRepository interface {
	Get(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginationResponse, error)
	GenerateInvoiceID(ctx context.Context, tx Transaction, userID uint, date time.Time) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Invoice, error)
	GetItems(ctx context.Context, invoiceID []int64) ([]domain.InvoiceItem, error)
	GetItemsByInvoiceID(ctx context.Context, invoiceID int64) ([]domain.InvoiceItem, error)
	Create(ctx context.Context, tx Transaction, data *domain.Invoice) error
	CreateItem(ctx context.Context, tx Transaction, data []domain.InvoiceItem) error
	Update(ctx context.Context, tx Transaction, data *domain.Invoice) error
	DeleteItemsByInvoiceID(ctx context.Context, tx Transaction, invoiceID int64) error
}
