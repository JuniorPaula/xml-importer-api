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

type InvoiceItemFilter struct {
	InvoiceID  string
	ProductID  string
	CreditType string
	OrderBy    string
	OrderDir   string
	Page       int
	PageSize   int
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

// FindAll retrieves invoice items based on the provided filter criteria.
func (r *InvoiceItemRepo) FindAll(filter InvoiceItemFilter) ([]models.InvoiceItem, int64, error) {
	var items []models.InvoiceItem
	var total int64

	query := r.db.Preload("Invoice.Customer").Model(&models.InvoiceItem{}).Preload("Invoice.Partner").Preload("Product").Preload("Invoice")

	if filter.InvoiceID != "" {
		query = query.Where("invoice_id = ?", filter.InvoiceID)
	}
	if filter.ProductID != "" {
		query = query.Where("product_id = ?", filter.ProductID)
	}
	if filter.CreditType != "" {
		query = query.Where("credit_type = ?", filter.CreditType)
	}
	if filter.OrderBy != "" {
		order := "asc"
		if filter.OrderDir == "desc" {
			order = "desc"
		}
		validOrdersFields := map[string]bool{
			"meter_name":     true,
			"meter_category": true,
		}
		if validOrdersFields[filter.OrderBy] {
			query = query.Order(filter.OrderBy + " " + order)
		}
	}

	query.Count(&total)

	offset := (filter.Page - 1) * filter.PageSize
	err := query.Offset(offset).Limit(filter.PageSize).Find(&items).Error

	return items, total, err
}

// FindByID retrieves an invoice item by its ID.
func (r *InvoiceItemRepo) FindByID(id int) (*models.InvoiceItem, error) {
	var item models.InvoiceItem
	err := r.db.Preload("Invoice.Customer").Where("id = ?", id).Preload("Invoice.Partner").Preload("Product").First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetSummary retrieves a summary of invoice items
func (r *InvoiceItemRepo) GetSummary() (*models.Summary, error) {
	var totalInvoices int64
	var totalItems int64
	var totalBilling float64
	var creditDistribution []struct {
		CreditType string
		Count      int64
	}
	var topProducts []models.ProductTotal
	var categoryTotals []models.CategoryTotal

	if err := r.db.Model(&models.Invoice{}).Count(&totalInvoices).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.InvoiceItem{}).Count(&totalItems).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.InvoiceItem{}).Select("SUM(total_price)").Scan(&totalBilling).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.InvoiceItem{}).
		Select("credit_type, COUNT(*) as count").
		Group("credit_type").
		Scan(&creditDistribution).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.InvoiceItem{}).
		Select("product_id, products.name AS product_name, SUM(total_price) as total").
		Joins("JOIN products ON products.id = invoice_items.product_id").
		Group("product_id, products.name").
		Order("total DESC").
		Limit(5).
		Scan(&topProducts).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&models.InvoiceItem{}).
		Select("meter_category, SUM(total_price) as total").
		Group("meter_category").
		Order("total DESC").
		Scan(&categoryTotals).Error; err != nil {
		return nil, err
	}

	// Convert credit distribution to a map for easier access
	creditMap := make(map[string]int64)
	for _, row := range creditDistribution {
		creditMap[row.CreditType] = row.Count
	}

	return &models.Summary{
		TotalInvoices:          totalInvoices,
		TotalInvoiceItems:      totalItems,
		TotalBilling:           totalBilling,
		CreditTypeDistribution: creditMap,
		TopProducts:            topProducts,
		MeterCategoryTotals:    categoryTotals,
	}, nil
}
