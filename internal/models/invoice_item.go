package models

type InvoiceItem struct {
	ID             uint `gorm:"primaryKey"`
	InvoiceID      string
	ProductID      string
	SkuName        string
	MeterName      string
	MeterCategory  string
	MeterSubCat    string
	MeterRegion    string
	Unit           string
	Quantity       float64
	UnitPrice      float64
	TotalPrice     float64 // BillingPreTaxTotal
	CreditType     string
	CreditPerc     float64
	EffectivePrice float64

	Invoice Invoice `gorm:"foreignKey:InvoiceID"`
	Product Product `gorm:"foreignKey:ProductID"`
}
