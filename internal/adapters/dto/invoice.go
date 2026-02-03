package dto

import "time"

type InvoicePageReq struct {
	PaginationRequest
}

// InvoiceItemRequest represents invoice item input
type InvoiceItemRequest struct {
	Description string `json:"description" validate:"required,min=1,max=500"`
	Qty         int    `json:"qty" validate:"required,min=1"`
	Price       int    `json:"price" validate:"required,min=0"`
}

// InvoiceItemDTO represents invoice item data for PDF generation
type InvoiceItemDTO struct {
	ID          uint   `json:"id"`
	InvoiceID   int64  `json:"invoice_id"`
	Description string `json:"description"`
	Qty         int    `json:"qty"`
	Price       int    `json:"price"`
	Total       int    `json:"total"`
}

// CreateInvoiceRequest represents invoice creation input
type InvoiceRequest struct {
	ID        int64                `json:"id" validate:"required"`
	Issuer    string               `json:"issuer" validate:"required,min=1,max=200"`
	Customer  string               `json:"customer" validate:"required,min=1,max=200"`
	IssueDate string               `json:"issue_date" validate:"required,datetime=2006-01-02"`
	DueDate   string               `json:"due_date" validate:"required,datetime=2006-01-02 15:04:05"`
	Note      string               `json:"note" validate:"max=1000"`
	Items     []InvoiceItemRequest `json:"items" validate:"required,min=1,dive"`
	UserID    uint                 `json:"-"`
}

// InvoiceItemResponse represents invoice item output
type InvoiceItemResponse struct {
	ID          uint   `json:"id"`
	InvoiceID   int64  `json:"invoice_id"`
	Description string `json:"description"`
	Qty         int    `json:"qty"`
	Price       int    `json:"price"`
	Total       int    `json:"total"`
	CreatedAt   string `json:"created_at"`
}

// InvoiceResponse represents invoice output
type InvoiceResponse struct {
	ID        int64                 `json:"id"`
	Customer  string                `json:"customer"`
	Issuer    string                `json:"issuer"`
	IssueDate string                `json:"issue_date"`
	DueDate   time.Time             `json:"due_date"`
	Note      string                `json:"note"`
	Status    string                `json:"status"`
	Items     []InvoiceItemResponse `json:"items,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// InvoiceListResponse represents list of invoices output
type InvoiceListResponse struct {
	Invoices []InvoiceResponse `json:"invoices"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}
