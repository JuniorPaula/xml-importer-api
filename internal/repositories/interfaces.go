package repositories

import (
	"importerapi/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	FindByID(id int) (*models.User, error)
	Create(user *models.User) error
}

type CustomerRepository interface {
	FindAll(filter CustomerFilter) ([]models.Customer, int64, error)
}

type PartnerRepository interface {
	WithTx(tx *gorm.DB) PartnerRepository
}

type ProductRepository interface {
	FindAll(filter ProductFilter) ([]models.Product, int64, error)
	FindByID(id string) (*models.Product, error)
}

type InvoiceRepository interface {
	FindAll(filter InvoiceFilter) ([]models.Invoice, int64, error)
}

type InvoiceItemRepository interface {
	FindAll(filter InvoiceItemFilter) ([]models.InvoiceItem, int64, error)
	FindByID(id int) (*models.InvoiceItem, error)
	GetSummary() (*models.Summary, error)
}

type ImportStatusRepository interface {
	Create(data *models.ImportStatus) error
	FindByImportID(importID string) (*models.ImportStatus, error)
	UpdateStatus(importID string, status string) error
}
