package domain

import (
	"app/xonvera-core/internal/adapters/dto"
	"time"
)

type Invoice struct {
	ID        int64
	AuthorID  uint
	Issuer    string
	Customer  string
	IssueDate string
	DueDate   time.Time
	Note      string
	Status    string
	Timestamp
}

type InvoiceItem struct {
	ID          uint
	InvoiceID   int64
	Description string
	Qty         int
	Price       int
	Total       int
	Timestamp
}

func (Invoice) TableName() string {
	return "app.invoices"
}

func (InvoiceItem) TableName() string {
	return "app.invoice_items"
}

func (i *Invoice) Response() dto.InvoiceResponse {
	return dto.InvoiceResponse{
		ID:        i.ID,
		Issuer:    i.Issuer,
		Customer:  i.Customer,
		IssueDate: i.IssueDate,
		DueDate:   i.DueDate,
		Note:      i.Note,
		Status:    i.Status,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}
