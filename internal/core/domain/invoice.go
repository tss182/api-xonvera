package domain

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Issuer    string    `json:"issuer" gorm:"not null"`
	Customer  string    `json:"customer" gorm:"not null"`
	IssueDate string    `json:"issue_date" gorm:"not null"`
	Note      string    `json:"note" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-" gorm:"index"`
}

type InvoiceItem struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	InvoiceID   int64          `json:"invoice_id" gorm:"not null;index"`
	Description string         `json:"description" gorm:"not null"`
	Qty         int            `json:"qty" gorm:"not null;default:1"`
	Price       int            `json:"price" gorm:"not null"`
	Total       int            `json:"total" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Invoice) TableName() string {
	return "billing.invoices"
}

func (InvoiceItem) TableName() string {
	return "billing.invoice_items"
}
