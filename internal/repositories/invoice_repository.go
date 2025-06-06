package repositories

import (
	"importerapi/internal/models"
	"time"

	"gorm.io/gorm"
)

type InvoiceRepo struct {
	db *gorm.DB
}

func NewInvoiceRepo(db *gorm.DB) *InvoiceRepo {
	return &InvoiceRepo{
		db: db,
	}
}

type InvoiceFilter struct {
	CustomerID string
	Month      int
	Year       int
	Page       int
	PageSize   int
}

// WithTx returns a new InvoiceRepository that uses the provided transaction.
func (r *InvoiceRepo) WithTx(tx *gorm.DB) InvoiceRepository {
	return &InvoiceRepo{
		db: tx,
	}
}

// FirstOrCreateInvoice attempts to find an invoice by ID, and if not found,
// creates a new invoice with the provided details.
func (r *InvoiceRepo) FirstOrCreateInvoice(invoice *models.Invoice) error {
	return r.db.Where("id = ?", invoice.ID).FirstOrCreate(invoice).Error
}

// FindAll retrieves invoices based on the provided filter criteria.
func (r *InvoiceRepo) FindAll(filter InvoiceFilter) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	query := r.db.Model(&models.Invoice{}).Preload("Customer").Preload("Partner")

	if filter.CustomerID != "" {
		query = query.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.Month > 0 && filter.Year > 0 {
		start := time.Date(filter.Year, time.Month(filter.Month), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, -1)
		query = query.Where("charge_start_date BETWEEN ? AND ?", start, end)
	}

	query.Count(&total)

	offset := (filter.Page - 1) * filter.PageSize
	err := query.Order("charge_start_date DESC").Offset(offset).Limit(filter.PageSize).Find(&invoices).Error

	return invoices, total, err
}
