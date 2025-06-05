package models

import "time"

type Invoice struct {
	ID               string `gorm:"primaryKey"` // InvoiceNumber
	PartnerID        string
	CustomerID       string
	ExchangeRate     float64
	ExchangeRateDate time.Time
	ChargeStartDate  time.Time
	ChargeEndDate    time.Time
	Customer         Customer `gorm:"foreignKey:CustomerID"`
	Partner          Partner  `gorm:"foreignKey:PartnerID"`
}
