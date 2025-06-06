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
