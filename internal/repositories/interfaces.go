package repositories

import (
	"importerapi/internal/models"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	WithTx(tx *gorm.DB) CustomerRepository
	FirstOrCreateCustomer(customer *models.Customer) error
	FindAll(filter CustomerFilter) ([]models.Customer, int64, error)
}

type PartnerRepository interface {
	WithTx(tx *gorm.DB) PartnerRepository
	FirstOrCreatePartner(partner *models.Partner) error
}

type ProductRepository interface {
	WithTx(tx *gorm.DB) ProductRepository
	FirstOrCreateProduct(product *models.Product) error
}

type InvoiceRepository interface {
	WithTx(tx *gorm.DB) InvoiceRepository
	FirstOrCreateInvoice(invoice *models.Invoice) error
}

type InvoiceItemRepository interface {
	WithTx(tx *gorm.DB) InvoiceItemRepository
	FirstOrCreateInvoiceItem(item *models.InvoiceItem) error
}
