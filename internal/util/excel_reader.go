package util

import (
	"fmt"
	"io"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type ExcelRecord struct {
	PartnerID     string
	PartnerName   string
	CustomerID    string
	CustomerName  string
	InvoiceNumber string
	ProductID     string
	CreditType    string
	CreditPerc    int
	UnitPrice     float64
	Quantity      float64
	Total         float64
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
			PartnerID:     safeGet(row, 0),
			PartnerName:   safeGet(row, 1),
			CustomerID:    safeGet(row, 2),
			CustomerName:  safeGet(row, 3),
			InvoiceNumber: safeGet(row, 4),
			ProductID:     safeGet(row, 5),
			CreditType:    safeGet(row, 6),
			CreditPerc:    parseInt(safeGet(row, 7)),
			UnitPrice:     parseFloat(safeGet(row, 8)),
			Quantity:      parseFloat(safeGet(row, 9)),
			Total:         parseFloat(safeGet(row, 10)),
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
