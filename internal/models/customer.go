package models

type Customer struct {
	ID      string `gorm:"primaryKey" json:"id"`
	Name    string `json:"name"`
	Domain  string `json:"domain"`
	Country string `json:"country"`
}
