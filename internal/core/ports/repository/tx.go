package portRepository

import (
	"gorm.io/gorm"
)

type TxRepository interface {
	Begin() (*gorm.DB, error)
}
