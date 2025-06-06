package repositories

import (
	"importerapi/internal/models"

	"gorm.io/gorm"
)

type InvoiceItemRepo struct {
	db *gorm.DB
}

func NewInvoiceItemRepo(db *gorm.DB) *InvoiceItemRepo {
	return &InvoiceItemRepo{
		db: db,
	}
}

type InvoiceItemFilter struct {
	InvoiceID  string
	ProductID  string
	CreditType string
	OrderBy    string
	OrderDir   string
	Page       int
	PageSize   int
}

// WithTx returns a new InvoiceItemRepository that uses the provided transaction.
func (r *InvoiceItemRepo) WithTx(tx *gorm.DB) InvoiceItemRepository {
	return &InvoiceItemRepo{
		db: tx,
	}
}

// FirstOrCreateInvoiceItem attempts to find an invoice item by ID, and if not found,
// creates a new invoice item with the provided details.
func (r *InvoiceItemRepo) FirstOrCreateInvoiceItem(invoiceItem *models.InvoiceItem) error {
	return r.db.Where("id = ?", invoiceItem.ID).FirstOrCreate(invoiceItem).Error
}

// FindAll retrieves invoice items based on the provided filter criteria.
func (r *InvoiceItemRepo) FindAll(filter InvoiceItemFilter) ([]models.InvoiceItem, int64, error) {
	var items []models.InvoiceItem
	var total int64

	query := r.db.Model(&models.InvoiceItem{}).Preload("Product").Preload("Invoice")

	if filter.InvoiceID != "" {
		query = query.Where("invoice_id = ?", filter.InvoiceID)
	}
	if filter.ProductID != "" {
		query = query.Where("product_id = ?", filter.ProductID)
	}
	if filter.CreditType != "" {
		query = query.Where("credit_type = ?", filter.CreditType)
	}
	if filter.OrderBy != "" {
		order := "asc"
		if filter.OrderDir == "desc" {
			order = "desc"
		}
		validOrdersFields := map[string]bool{
			"meter_name":     true,
			"meter_category": true,
		}
		if validOrdersFields[filter.OrderBy] {
			query = query.Order(filter.OrderBy + " " + order)
		}
	}

	query.Count(&total)

	offset := (filter.Page - 1) * filter.PageSize
	err := query.Offset(offset).Limit(filter.PageSize).Find(&items).Error

	return items, total, err
}
