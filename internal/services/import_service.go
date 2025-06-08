package services

import (
	"fmt"
	"importerapi/internal/models"
	"importerapi/internal/repositories"
	"importerapi/internal/util"
	"time"

	"gorm.io/gorm"
)

type ImportService struct {
	DB               *gorm.DB
	CustomerRepo     repositories.CustomerRepository
	PartnerRepo      repositories.PartnerRepository
	ProductRepo      repositories.ProductRepository
	InvoiceRepo      repositories.InvoiceRepository
	InvoiceItemRepo  repositories.InvoiceItemRepository
	ImportStatusRepo repositories.ImportStatusRepository
}

func NewImportService(db *gorm.DB) *ImportService {
	return &ImportService{
		DB:               db,
		CustomerRepo:     repositories.NewCustomerRepo(db),
		PartnerRepo:      repositories.NewPartnerRepo(db),
		ProductRepo:      repositories.NewProductRepo(db),
		InvoiceRepo:      repositories.NewInvoiceRepo(db),
		InvoiceItemRepo:  repositories.NewInvoiceItemRepo(db),
		ImportStatusRepo: repositories.NewImportStatusRepo(db),
	}
}

func (s *ImportService) ImportFromXML(records []util.ExcelRecord) error {
	tx := s.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, record := range records {
		if err := s.importRecord(tx, record); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *ImportService) importRecord(tx *gorm.DB, record util.ExcelRecord) error {
	partner := models.Partner{
		ID:         record.ParterID,
		Name:       record.PartnerName,
		MpnID:      record.MpnID,
		Tier2MpnID: record.Tier2MpnID,
	}
	err := s.PartnerRepo.WithTx(tx).FirstOrCreatePartner(&partner)
	if err != nil {
		return fmt.Errorf("failed to create or find partner: %w", err)
	}

	customer := models.Customer{
		ID:      record.CustomerID,
		Name:    record.CustomerName,
		Domain:  record.CustomerDomainName,
		Country: record.CustomerCountry,
	}
	err = s.CustomerRepo.WithTx(tx).FirstOrCreateCustomer(&customer)
	if err != nil {
		return fmt.Errorf("failed to create or find customer: %w", err)
	}

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
	err = s.ProductRepo.WithTx(tx).FirstOrCreateProduct(&product)
	if err != nil {
		return fmt.Errorf("failed to create or find product: %w", err)
	}

	invoice := models.Invoice{
		ID:               record.InvoiceNumber,
		PartnerID:        record.ParterID,
		CustomerID:       record.CustomerID,
		ChargeStartDate:  s.parseDate(record.ChargeStartDate),
		ChargeEndDate:    s.parseDate(record.ChargeEndDate),
		ExchangeRateDate: s.parseDate(record.PCToBCExchangeRateDate),
		ExchangeRate:     float64(record.PCToBCExchangeRate),
	}
	err = s.InvoiceRepo.WithTx(tx).FirstOrCreateInvoice(&invoice)
	if err != nil {
		return fmt.Errorf("failed to create or find invoice: %w", err)
	}

	item := models.InvoiceItem{
		InvoiceID:                     record.InvoiceNumber,
		ProductID:                     record.ProductID,
		MeterID:                       record.MeterID,
		MeterName:                     record.MeterName,
		MeterType:                     record.MeterType,
		MeterCategory:                 record.MeterCategory,
		MeterSubCategory:              record.MeterSubCategory,
		MeterRegion:                   record.MeterRegion,
		ResourceURI:                   record.ResourceURI,
		Quantity:                      record.Quantity,
		UnitPrice:                     record.UnitPrice,
		TotalPrice:                    record.BillingPreTaxTotal,
		EffectiveUnitPrice:            record.EffectiveUnitPrice,
		Unit:                          record.Unit,
		UnitType:                      record.UnitType,
		ChargeType:                    record.ChargeType,
		BillingCurrency:               record.BillingCurrency,
		PricingCurrency:               record.PricingCurrency,
		ServiceInfo1:                  record.ServiceInfo1,
		ServiceInfo2:                  record.ServiceInfo2,
		CreditType:                    record.CreditType,
		CreditPercentage:              record.CreditPercentage,
		PartnerEarnedCreditPercentage: record.PartnerEarnedCreditPercentage,
		Tags:                          record.Tags,
		AdditionalInfo:                record.AdditionalInfo,
	}
	err = s.InvoiceItemRepo.WithTx(tx).FirstOrCreateInvoiceItem(&item)
	if err != nil {
		return fmt.Errorf("failed to create or find invoice item: %w", err)
	}
	return nil
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
	err := s.ImportStatusRepo.UpdateStatus(importID, status)
	if err != nil {
		return fmt.Errorf("failed to update import status: %w", err)
	}
	return nil
}
