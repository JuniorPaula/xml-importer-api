package models

type Summary struct {
	TotalInvoices          int64            `json:"total_invoices"`
	TotalInvoiceItems      int64            `json:"total_invoice_items"`
	TotalBilling           float64          `json:"total_billing"`
	CreditTypeDistribution map[string]int64 `json:"credit_type_distribution"`
	TopProducts            []ProductTotal   `json:"top_products"`
	MeterCategoryTotals    []CategoryTotal  `json:"meter_category_totals"`
}

type ProductTotal struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Total       float64 `json:"total"`
}

type CategoryTotal struct {
	MeterCategory string  `json:"meter_category"`
	Total         float64 `json:"total"`
}
