package entity

import "gorm.io/gorm"

// Example {"id":1,"name":"test"}
type Example struct {
	gorm.Model

	Name string `json:"name"`
}
