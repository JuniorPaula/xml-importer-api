package repositories

import (
	"importerapi/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}

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
	FindAll(filter ProductFilter) ([]models.Product, int64, error)
	FindByID(id string) (*models.Product, error)
}

type InvoiceRepository interface {
	WithTx(tx *gorm.DB) InvoiceRepository
	FirstOrCreateInvoice(invoice *models.Invoice) error
	FindAll(filter InvoiceFilter) ([]models.Invoice, int64, error)
}

type InvoiceItemRepository interface {
	WithTx(tx *gorm.DB) InvoiceItemRepository
	FirstOrCreateInvoiceItem(item *models.InvoiceItem) error
	FindAll(filter InvoiceItemFilter) ([]models.InvoiceItem, int64, error)
	FindByID(id int) (*models.InvoiceItem, error)
	GetSummary() (*models.Summary, error)
}
