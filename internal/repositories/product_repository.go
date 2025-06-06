package repositories

import (
	"importerapi/internal/models"

	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

// WithTx returns a new ProductRepository that uses the provided transaction.
func (r *ProductRepo) WithTx(tx *gorm.DB) ProductRepository {
	return &ProductRepo{
		db: tx,
	}
}

// FirstOrCreateProduct attempts to find a product by ID, and if not found,
// creates a new product with the provided details.
func (r *ProductRepo) FirstOrCreateProduct(product *models.Product) error {
	return r.db.Where("id = ?", product.ID).FirstOrCreate(product).Error
}
