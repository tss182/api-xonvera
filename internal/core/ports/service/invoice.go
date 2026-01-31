package portService

import (
	"context"

	"app/xonvera-core/internal/adapters/dto"
)

type InvoiceService interface {
	Get(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginationResponse, error)
	Create(ctx context.Context, req *dto.InvoiceRequest) error
	Update(ctx context.Context, req *dto.InvoiceRequest) error
}
