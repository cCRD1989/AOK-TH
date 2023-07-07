package model

import (
	"gorm.io/gorm"
)

type Setup struct {
	Type  string
	Value string
	gorm.Model
}
