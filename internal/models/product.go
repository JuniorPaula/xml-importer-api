package models

type Product struct {
	ID              string `gorm:"primaryKey" json:"id"`
	SkuID           string `json:"sku_id"`
	AvailabilityID  string `json:"availability_id"`
	SkuName         string `json:"sku_name"`
	Name            string `json:"name"`
	PublisherName   string `json:"publisher_name"`
	PublisherID     string `json:"publisher_id"`
	EntitlementID   string `json:"entitlement_id"`
	EntitlementDesc string `json:"entitlement_desc"`
}
