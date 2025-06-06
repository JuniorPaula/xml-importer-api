package models

import "time"

type Invoice struct {
	ID               string    `gorm:"primaryKey" json:"id"` // InvoiceNumber
	PartnerID        string    `json:"partner_id"`
	CustomerID       string    `json:"customer_id"`
	ExchangeRate     float64   `json:"exchange_rate"`
	ExchangeRateDate time.Time `json:"exchange_rate_date"`
	ChargeStartDate  time.Time `json:"charge_start_date"`
	ChargeEndDate    time.Time `json:"charge_end_date"`
	Customer         Customer  `gorm:"foreignKey:CustomerID" json:"customer"`
	Partner          Partner   `gorm:"foreignKey:PartnerID" json:"partner"`
}
