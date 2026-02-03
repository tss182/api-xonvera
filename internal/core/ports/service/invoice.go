package portService

import (
	"context"

	"app/xonvera-core/internal/adapters/dto"
)

type InvoiceService interface {
	Get(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginationResponse, error)
	GetByID(ctx context.Context, invoiceID int64, userID uint) (*dto.InvoiceResponse, error)
	Create(ctx context.Context, req *dto.InvoiceRequest) error
	Update(ctx context.Context, req *dto.InvoiceRequest) error
}
