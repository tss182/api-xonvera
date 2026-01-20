package domain

import (
	"time"

	"gorm.io/gorm"
)

type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeAmount     DiscountType = "amount"
)

type Package struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null"`
	Price        int            `json:"price" gorm:"not null"`
	DiscountType DiscountType   `json:"discount_type" gorm:"not null;default:'percentage'"`
	Discount     int            `json:"discount" gorm:"not null;default:0"`
	Duration     string         `json:"duration" gorm:"not null;default:'1d'"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Package) TableName() string {
	return "catalog.packages"
}
