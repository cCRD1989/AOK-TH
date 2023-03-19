package model

import "gorm.io/gorm"

type LogWeb struct {
	gorm.Model
	DataType  string `gorm:"not null"`
	IPAddress string `gorm:"not null"`
}
