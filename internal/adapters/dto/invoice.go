package dto

// InvoiceItemRequest represents invoice item input
type InvoiceItemRequest struct {
	Description string `json:"description" validate:"required,min=1,max=500"`
	Quantity    int    `json:"quantity" validate:"required,min=1"`
	UnitPrice   int    `json:"unit_price" validate:"required,min=0"`
}

// CreateInvoiceRequest represents invoice creation input
type CreateInvoiceRequest struct {
	AddTo       string               `json:"add_to" validate:"required,min=1,max=500"`
	InvoiceFor  string               `json:"invoice_for" validate:"required,min=1,max=500"`
	InvoiceFrom string               `json:"invoice_from" validate:"required,min=1,max=500"`
	Items       []InvoiceItemRequest `json:"items" validate:"required,min=1,dive"`
}

// InvoiceItemResponse represents invoice item output
type InvoiceItemResponse struct {
	ID          uint   `json:"id"`
	InvoiceID   int64  `json:"invoice_id"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	UnitPrice   int    `json:"unit_price"`
	Total       int    `json:"total"`
	CreatedAt   string `json:"created_at"`
}

// InvoiceResponse represents invoice output
type InvoiceResponse struct {
	ID          int64                 `json:"id"`
	AddTo       string                `json:"add_to"`
	InvoiceFor  string                `json:"invoice_for"`
	InvoiceFrom string                `json:"invoice_from"`
	Items       []InvoiceItemResponse `json:"items,omitempty"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
}

// InvoiceListResponse represents list of invoices output
type InvoiceListResponse struct {
	Invoices []InvoiceResponse `json:"invoices"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}
