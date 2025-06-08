package repositories

import (
	"importerapi/internal/models"

	"gorm.io/gorm"
)

type ImportStatusRepo struct {
	DB *gorm.DB
}

func NewImportStatusRepo(db *gorm.DB) *ImportStatusRepo {
	return &ImportStatusRepo{
		DB: db,
	}
}

func (r *ImportStatusRepo) Create(data *models.ImportStatus) error {
	return r.DB.Create(data).Error
}

func (r *ImportStatusRepo) FindByImportID(importID string) (*models.ImportStatus, error) {
	var status models.ImportStatus
	err := r.DB.Where("import_id = ?", importID).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *ImportStatusRepo) UpdateStatus(importID string, status string) error {
	var importStatus models.ImportStatus
	err := r.DB.Where("import_id = ?", importID).First(&importStatus).Error
	if err != nil {
		return err
	}

	importStatus.Status = status
	return r.DB.Save(&importStatus).Error
}
