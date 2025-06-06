package repositories

import (
	"importerapi/internal/models"

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
