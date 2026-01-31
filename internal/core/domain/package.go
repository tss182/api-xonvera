package domain

type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeAmount     DiscountType = "amount"
)

type Package struct {
	ID           string
	Name         string
	Price        int
	DiscountType DiscountType
	Discount     int
	Duration     string
	Timestamp
}

func (Package) TableName() string {
	return "app.packages"
}
