package repositories

import (
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
