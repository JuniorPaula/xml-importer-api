package models

type InvoiceItem struct {
	ID                            uint `gorm:"primaryKey"`
	InvoiceID                     string
	ProductID                     string
	MeterID                       string
	MeterName                     string
	MeterType                     string
	MeterCategory                 string
	MeterSubCategory              string
	MeterRegion                   string
	ResourceURI                   string
	Quantity                      float64
	UnitPrice                     float64
	TotalPrice                    float64
	EffectiveUnitPrice            float64
	Unit                          string
	UnitType                      string
	ChargeType                    string
	BillingCurrency               string
	PricingCurrency               string
	ServiceInfo1                  string
	ServiceInfo2                  string
	CreditType                    string
	CreditPercentage              int
	PartnerEarnedCreditPercentage int
	Tags                          string
	AdditionalInfo                string

	Invoice Invoice `gorm:"foreignKey:InvoiceID"`
	Product Product `gorm:"foreignKey:ProductID"`
}
