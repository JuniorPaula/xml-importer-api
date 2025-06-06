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

type ProductFilter struct {
	Name        string
	PublisherID string
	Page        int
	PageSize    int
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

// FindAll retrieves products based on the provided filter criteria.
func (r *ProductRepo) FindAll(filter ProductFilter) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.PublisherID != "" {
		query = query.Where("publisher_id = ?", filter.PublisherID)
	}

	query.Count(&total)

	offset := (filter.Page - 1) * filter.PageSize
	err := query.Offset(offset).Limit(filter.PageSize).Find(&products).Error

	return products, total, err
}

// FindByID retrieves a product by its ID.
func (r *ProductRepo) FindByID(id string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}
