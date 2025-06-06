package repositories

import (
	"importerapi/internal/models"

	"gorm.io/gorm"
)

type CustomerRepo struct {
	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) *CustomerRepo {
	return &CustomerRepo{
		db: db,
	}
}

// WithTx returns a new CustomerRepository that uses the provided transaction.
func (r *CustomerRepo) WithTx(tx *gorm.DB) CustomerRepository {
	return &CustomerRepo{
		db: tx,
	}
}

// FirstOrCreateCustomer attempts to find a customer by ID, and if not found,
// creates a new customer with the provided details.
func (r *CustomerRepo) FirstOrCreateCustomer(customer *models.Customer) error {
	return r.db.Where("id = ?", customer.ID).FirstOrCreate(customer).Error
}
