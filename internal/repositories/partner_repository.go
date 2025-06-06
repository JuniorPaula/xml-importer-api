package repositories

import (
	"importerapi/internal/models"

	"gorm.io/gorm"
)

type PartnerRepo struct {
	db *gorm.DB
}

func NewPartnerRepo(db *gorm.DB) *PartnerRepo {
	return &PartnerRepo{
		db: db,
	}
}

func (r *PartnerRepo) WithTx(tx *gorm.DB) PartnerRepository {
	return &PartnerRepo{
		db: tx,
	}
}

// FirstOrCreatePartner attempts to find a partner by ID, and if not found,
// creates a new partner with the provided details.
func (r *PartnerRepo) FirstOrCreatePartner(partner *models.Partner) error {
	return r.db.Where("id = ?", partner.ID).FirstOrCreate(partner).Error
}
