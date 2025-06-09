package services

import (
	"fmt"
	"importerapi/internal/models"
	"importerapi/internal/repositories"
	"importerapi/internal/util"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ImportService struct {
	DB               *gorm.DB
	ImportStatusRepo repositories.ImportStatusRepository
}

func NewImportService(db *gorm.DB) *ImportService {
	return &ImportService{
		DB:               db,
		ImportStatusRepo: repositories.NewImportStatusRepo(db),
	}
}

func (s *ImportService) ImportFromXML(records []util.ExcelRecord) error {
	tx := s.DB.Session(&gorm.Session{
		SkipHooks:            true,
		FullSaveAssociations: false,
	}).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()

	// cache to control the number of calls to the database
	partnerCache := map[string]bool{}
	customerCache := map[string]bool{}
	productCache := map[string]bool{}
	invoiceCache := map[string]bool{}

	for _, record := range records {
		if err := s.importDependencies(tx, record, partnerCache, customerCache, productCache, invoiceCache); err != nil {
			tx.Rollback()
			return err
		}
	}

	invoiceItems := s.BuildInvoiceItems(records)
	if err := tx.CreateInBatches(invoiceItems, 1000).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to batch insert invoice items: %w", err)
	}

	return tx.Commit().Error
}

func (s *ImportService) importDependencies(
	tx *gorm.DB, record util.ExcelRecord,
	partnerCache, customerCache, productCache, invoiceCache map[string]bool,
) error {
	if !partnerCache[record.ParterID] {
		partner := models.Partner{
			ID:         record.ParterID,
			Name:       record.PartnerName,
			MpnID:      record.MpnID,
			Tier2MpnID: record.Tier2MpnID,
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&partner).Error; err != nil {
			return fmt.Errorf("failed to create partner: %w", err)
		}
		partnerCache[record.ParterID] = true

	}

	if !customerCache[record.CustomerID] {
		customer := models.Customer{
			ID:      record.CustomerID,
			Name:    record.CustomerName,
			Domain:  record.CustomerDomainName,
			Country: record.CustomerCountry,
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&customer).Error; err != nil {
			return fmt.Errorf("failed to create customer: %w", err)
		}
		customerCache[record.CustomerID] = true
	}

	if !productCache[record.ProductID] {
		product := models.Product{
			ID:              record.ProductID,
			Name:            record.ProductName,
			SkuID:           record.SKUID,
			SkuName:         record.SKUName,
			AvailabilityID:  record.AvailabilityID,
			PublisherName:   record.PublisherName,
			PublisherID:     record.PublisherID,
			EntitlementID:   record.EntitlementId,
			EntitlementDesc: record.EntitlementDescription,
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&product).Error; err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}
		productCache[record.ProductID] = true
	}

	if !invoiceCache[record.InvoiceNumber] {
		invoice := models.Invoice{
			ID:               record.InvoiceNumber,
			PartnerID:        record.ParterID,
			CustomerID:       record.CustomerID,
			ChargeStartDate:  s.parseDate(record.ChargeStartDate),
			ChargeEndDate:    s.parseDate(record.ChargeEndDate),
			ExchangeRateDate: s.parseDate(record.PCToBCExchangeRateDate),
			ExchangeRate:     float64(record.PCToBCExchangeRate),
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&invoice).Error; err != nil {
			return fmt.Errorf("failed to create invoice: %w", err)
		}
		invoiceCache[record.InvoiceNumber] = true
	}

	return nil
}

func (s *ImportService) BuildInvoiceItems(records []util.ExcelRecord) []models.InvoiceItem {
	items := make([]models.InvoiceItem, 0, len(records))
	for i, r := range records {
		if i%10000 == 0 && i != 0 {
			fmt.Printf("Progressing: %d records processed\n", i)
		}

		items = append(items, models.InvoiceItem{
			InvoiceID:                     r.InvoiceNumber,
			ProductID:                     r.ProductID,
			MeterID:                       r.MeterID,
			MeterName:                     r.MeterName,
			MeterType:                     r.MeterType,
			MeterCategory:                 r.MeterCategory,
			MeterSubCategory:              r.MeterSubCategory,
			MeterRegion:                   r.MeterRegion,
			ResourceURI:                   r.ResourceURI,
			Quantity:                      r.Quantity,
			UnitPrice:                     r.UnitPrice,
			TotalPrice:                    r.BillingPreTaxTotal,
			EffectiveUnitPrice:            r.EffectiveUnitPrice,
			Unit:                          r.Unit,
			UnitType:                      r.UnitType,
			ChargeType:                    r.ChargeType,
			BillingCurrency:               r.BillingCurrency,
			PricingCurrency:               r.PricingCurrency,
			ServiceInfo1:                  r.ServiceInfo1,
			ServiceInfo2:                  r.ServiceInfo2,
			CreditType:                    r.CreditType,
			CreditPercentage:              r.CreditPercentage,
			PartnerEarnedCreditPercentage: r.PartnerEarnedCreditPercentage,
			Tags:                          r.Tags,
			AdditionalInfo:                r.AdditionalInfo,
		})
	}
	return items
}

func (s *ImportService) parseDate(dataStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dataStr)
	return t
}

func (s *ImportService) GetImportStatus(importID string) (*models.ImportStatus, error) {
	status, err := s.ImportStatusRepo.FindByImportID(importID)
	if err != nil {
		return nil, fmt.Errorf("failed to find import status: %w", err)
	}
	if status == nil {
		return nil, fmt.Errorf("import status not found for import ID: %s", importID)
	}
	return status, nil
}

func (s *ImportService) UpdateImportStatus(importID string, status string) error {
	if err := s.ImportStatusRepo.UpdateStatus(importID, status); err != nil {
		return fmt.Errorf("failed to update import status: %w", err)
	}
	return nil
}
