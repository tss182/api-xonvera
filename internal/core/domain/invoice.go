package domain

import (
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

func (i *Invoice) Response(items []InvoiceItem) InvoiceResponse {
	var itemResponses []InvoiceItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, item.Response())
	}
	return InvoiceResponse{
		ID:        i.ID,
		Issuer:    i.Issuer,
		Customer:  i.Customer,
		IssueDate: i.IssueDate,
		DueDate:   i.DueDate,
		Note:      i.Note,
		Items:     itemResponses,
		Status:    i.Status,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

func (i *InvoiceItem) Response() InvoiceItemResponse {
	return InvoiceItemResponse{
		ID:          i.ID,
		Description: i.Description,
		Qty:         i.Qty,
		Price:       i.Price,
		Total:       i.Total,
		CreatedAt:   i.CreatedAt,
	}
}
