package models

type Partner struct {
	ID         string `gorm:"primaryKey"`
	Name       string
	MpnID      int
	Tier2MpnID int
}
