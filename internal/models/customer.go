package models

type Customer struct {
	ID      string `gorm:"primaryKey"`
	Name    string
	Domain  string
	Country string
}
