package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Name      string
	Email     string
	Tel       string
	Product []OrderItem `gorm:"foreignKey:OrderID"`
}
