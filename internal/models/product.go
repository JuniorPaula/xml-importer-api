package models

type Product struct {
	ID              string `gorm:"primaryKey"`
	SkuID           string
	AvailabilityID  string
	SkuName         string
	Name            string // ProductName
	PublisherName   string
	PublisherID     string
	EntitlementID   string
	EntitlementDesc string
}
