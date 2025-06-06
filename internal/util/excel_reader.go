package util

import (
	"fmt"
	"io"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type ExcelRecord struct {
	ParterID                      string
	PartnerName                   string
	CustomerID                    string
	CustomerName                  string
	CustomerDomainName            string
	CustomerCountry               string
	MpnID                         int
	Tier2MpnID                    int
	InvoiceNumber                 string
	ProductID                     string
	SKUID                         string
	AvailabilityID                string
	SKUName                       string
	ProductName                   string
	PublisherName                 string
	PublisherID                   string
	SubscriptionDescription       string
	SubscriptionID                string
	ChargeStartDate               string
	ChargeEndDate                 string
	UsageDate                     string
	MeterType                     string
	MeterCategory                 string
	MeterID                       string
	MeterSubCategory              string
	MeterName                     string
	MeterRegion                   string
	Unit                          string
	ResourceLocation              string
	CostomerService               string
	ResourceGroup                 string
	ResourceURI                   string
	ChargeType                    string
	UnitPrice                     float64
	Quantity                      float64
	UnitType                      string
	BillingPreTaxTotal            float64
	BillingCurrency               string
	PricingPreTaxTotal            float64
	PricingCurrency               string
	ServiceInfo1                  string
	ServiceInfo2                  string
	Tags                          string
	AdditionalInfo                string
	EffectiveUnitPrice            float64
	PCToBCExchangeRate            int
	PCToBCExchangeRateDate        string
	EntitlementId                 string
	EntitlementDescription        string
	PartnerEarnedCreditPercentage int
	CreditPercentage              int
	CreditType                    string
	BenefitOrderId                string
	BenefitId                     string
	BenefitType                   string
}

// ReadExcelFromReader usa stream para leitura eficiente de arquivos grandes
func ReadExcelFromReader(reader io.Reader) ([]ExcelRecord, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("error to open excel file: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("empty sheet name, please check the file")
	}

	stream, err := f.Rows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("error to open rows: %w", err)
	}
	defer stream.Close()

	var records []ExcelRecord
	rowIndex := 0

	for stream.Next() {
		row, err := stream.Columns()
		if err != nil {
			return nil, fmt.Errorf("error to read rows: %w", err)
		}

		if rowIndex == 0 {
			rowIndex++
			continue
		}

		if isEmptyRow(row) {
			fmt.Printf("Line %d with just %d columns, skipping\n", rowIndex+1, len(row))
			continue
		}

		record := ExcelRecord{
			ParterID:                      safeGet(row, 0),
			PartnerName:                   safeGet(row, 1),
			CustomerID:                    safeGet(row, 2),
			CustomerName:                  safeGet(row, 3),
			CustomerDomainName:            safeGet(row, 4),
			CustomerCountry:               safeGet(row, 5),
			MpnID:                         parseInt(safeGet(row, 6)),
			Tier2MpnID:                    parseInt(safeGet(row, 7)),
			InvoiceNumber:                 safeGet(row, 8),
			ProductID:                     safeGet(row, 9),
			SKUID:                         safeGet(row, 10),
			AvailabilityID:                safeGet(row, 11),
			SKUName:                       safeGet(row, 12),
			ProductName:                   safeGet(row, 13),
			PublisherName:                 safeGet(row, 14),
			PublisherID:                   safeGet(row, 15),
			SubscriptionDescription:       safeGet(row, 16),
			SubscriptionID:                safeGet(row, 17),
			ChargeStartDate:               safeGet(row, 18),
			ChargeEndDate:                 safeGet(row, 19),
			UsageDate:                     safeGet(row, 20),
			MeterType:                     safeGet(row, 21),
			MeterCategory:                 safeGet(row, 22),
			MeterID:                       safeGet(row, 23),
			MeterSubCategory:              safeGet(row, 24),
			MeterName:                     safeGet(row, 25),
			MeterRegion:                   safeGet(row, 26),
			Unit:                          safeGet(row, 27),
			ResourceLocation:              safeGet(row, 28),
			CostomerService:               safeGet(row, 29),
			ResourceGroup:                 safeGet(row, 30),
			ResourceURI:                   safeGet(row, 31),
			ChargeType:                    safeGet(row, 32),
			UnitPrice:                     parseFloat(safeGet(row, 33)),
			Quantity:                      parseFloat(safeGet(row, 34)),
			UnitType:                      safeGet(row, 35),
			BillingPreTaxTotal:            parseFloat(safeGet(row, 36)),
			BillingCurrency:               safeGet(row, 37),
			PricingPreTaxTotal:            parseFloat(safeGet(row, 38)),
			PricingCurrency:               safeGet(row, 39),
			ServiceInfo1:                  safeGet(row, 40),
			ServiceInfo2:                  safeGet(row, 41),
			Tags:                          safeGet(row, 42),
			AdditionalInfo:                safeGet(row, 43),
			EffectiveUnitPrice:            parseFloat(safeGet(row, 44)),
			PCToBCExchangeRate:            parseInt(safeGet(row, 45)),
			PCToBCExchangeRateDate:        safeGet(row, 46),
			EntitlementId:                 safeGet(row, 47),
			EntitlementDescription:        safeGet(row, 48),
			PartnerEarnedCreditPercentage: parseInt(safeGet(row, 49)),
			CreditPercentage:              parseInt(safeGet(row, 50)),
			CreditType:                    safeGet(row, 51),
			BenefitOrderId:                safeGet(row, 52),
			BenefitId:                     safeGet(row, 53),
			BenefitType:                   safeGet(row, 54),
		}

		records = append(records, record)
		rowIndex++
	}

	return records, nil
}

func isEmptyRow(rows []string) bool {
	for _, cell := range rows {
		if cell != "" {
			return false
		}
	}
	return true
}

func safeGet(row []string, i int) string {
	if i < len(row) {
		return row[i]
	}
	return ""
}

func parseFloat(s string) float64 {
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func parseInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}
