package domain

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	ID          int64          `json:"id" gorm:"primaryKey"`
	AddTo       string         `json:"add_to" gorm:"not null"`
	InvoiceFor  string         `json:"invoice_for" gorm:"not null"`
	InvoiceFrom string         `json:"invoice_from" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type InvoiceItem struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	InvoiceID   int64          `json:"invoice_id" gorm:"not null;index"`
	Description string         `json:"description" gorm:"not null"`
	Quantity    int            `json:"quantity" gorm:"not null;default:1"`
	UnitPrice   int            `json:"unit_price" gorm:"not null"`
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
