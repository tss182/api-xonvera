package domain

import (
	"time"

	"gorm.io/gorm"
)

type DiscountType string

const (
	DiscountTypePercentage DiscountType = "PERCENTAGE"
	DiscountTypeAmount     DiscountType = "AMOUNT"
)

type Package struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null"`
	Price        int            `json:"price" gorm:"not null"`
	DiscountType DiscountType   `json:"discount_type" gorm:"not null;default:'PERCENTAGE'"`
	Discount     int            `json:"discount" gorm:"not null;default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Package) TableName() string {
	return "catalog.packages"
}
