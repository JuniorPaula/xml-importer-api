package models

import "time"

type ImportStatus struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	ImportID  string    `json:"import_id" gorm:"uniqueIndex;not null"`
	FileName  string    `json:"file_name" gorm:"not null"`
	Status    string    `json:"status" gorm:"not null"` // e.g., "pending", "in_progress", "completed", "failed"
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
