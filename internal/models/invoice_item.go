package models

type InvoiceItem struct {
	ID                            uint    `gorm:"primaryKey" json:"id"`
	InvoiceID                     string  `json:"invoice_id"`
	ProductID                     string  `json:"product_id"`
	MeterID                       string  `json:"meter_id"`
	MeterName                     string  `json:"meter_name"`
	MeterType                     string  `json:"meter_type"`
	MeterCategory                 string  `json:"meter_category"`
	MeterSubCategory              string  `json:"meter_sub_category"`
	MeterRegion                   string  `json:"meter_region"`
	ResourceURI                   string  `json:"resource_uri"`
	Quantity                      float64 `json:"quantity"`
	UnitPrice                     float64 `json:"unit_price"`
	TotalPrice                    float64 `json:"total_price"`
	EffectiveUnitPrice            float64 `json:"effective_unit_price"`
	Unit                          string  `json:"unit"`
	UnitType                      string  `json:"unit_type"`
	ChargeType                    string  `json:"charge_type"`
	BillingCurrency               string  `json:"billing_currency"`
	PricingCurrency               string  `json:"pricing_currency"`
	ServiceInfo1                  string  `json:"service_info_1"`
	ServiceInfo2                  string  `json:"service_info_2"`
	CreditType                    string  `json:"credit_type"`
	CreditPercentage              int     `json:"credit_percentage"`
	PartnerEarnedCreditPercentage int     `json:"partner_earned_credit_percentage"`
	Tags                          string  `json:"tags"`
	AdditionalInfo                string  `json:"additional_info"`

	Invoice Invoice `gorm:"foreignKey:InvoiceID" json:"invoice"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
