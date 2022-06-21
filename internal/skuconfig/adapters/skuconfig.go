package adapters

import (
	"time"
)

type SKUConfigModel struct {
	UUID          string `gorm:"primaryKey"`
	Package       string `gorm:"index"`
	CountryCode   string `gorm:"index"`
	PercentileMin uint
	PercentileMax uint
	SKU           string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (SKUConfigModel) TableName() string {
	return "sku_configs"
}
