package model

import (
	"gorm.io/gorm"
)

type Setting struct {
	Type  byte
	Keys  string
	Value string
	gorm.Model
}
