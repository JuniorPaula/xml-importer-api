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

type CustomerFilter struct {
	Name     string
	Country  string
	Page     int
	PageSize int
}

// FindAll retrieves customers based on the provided filter criteria.
func (r *CustomerRepo) FindAll(filter CustomerFilter) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	query := r.db.Model(&models.Customer{})

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.Country != "" {
		query = query.Where("country = ?", filter.Country)
	}
	query.Count(&total)

	offset := (filter.Page - 1) * filter.PageSize
	err := query.Offset(offset).Limit(filter.PageSize).Find(&customers).Error

	return customers, total, err
}
