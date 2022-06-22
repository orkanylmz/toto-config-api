package adapters

import (
	"time"
)

type SKUConfigModel struct {
	ID            string `gorm:"primaryKey"`
	Package       string `gorm:"index"`
	CountryCode   string `gorm:"index"`
	PercentileMin uint
	PercentileMax uint
	SKU           string
	CreatedAt     time.Time `gorm:"type:timestamp"`
	UpdatedAt     time.Time `gorm:"type:timestamp"`
}

func (SKUConfigModel) TableName() string {
	return "sku_configs"
}
