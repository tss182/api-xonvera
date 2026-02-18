package portService

import (
	"context"

	"app/xonvera-core/internal/core/domain"
)

type InvoiceService interface {
	Get(ctx context.Context, req *domain.PaginationRequest) (*domain.PaginationResponse, error)
	GetByID(ctx context.Context, invoiceID int64, userID uint) (*domain.InvoiceResponse, error)
	Create(ctx context.Context, req *domain.InvoiceRequest) error
	Update(ctx context.Context, req *domain.InvoiceRequest) error
	GetPDF(ctx context.Context, invoiceID int64, userID uint) ([]byte, error)
}
