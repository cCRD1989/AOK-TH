package model

import "gorm.io/gorm"

type LogTopup struct {
	gorm.Model
	DataType  string `gorm:"not null"`
	UserId    string
	Txid      string
	Orderid   string
	Status    string
	Detail    string
	Channel   string
	Price     string
	Sig       string
	IPAddress string `gorm:"not null"`
}
