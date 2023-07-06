package model

import (
	"time"

	"gorm.io/gorm"
)

type LogTopup struct {
	gorm.Model
	DataType string `gorm:"not null"`
	UserId   string
	Txid     string
	Orderid  string
	Status   string
	Detail   string
	Channel  string
	Price    string
	Bonus    int `gorm:"not null;default:0"`
	Sig      string

	IPAddress string `gorm:"not null"`
}

type LogMailTopup struct {
	Id              int64     `gorm:"uniqueIndex;autoIncrement:true;type:bigint(20)"`
	Eventid         string    `gorm:"type:varchar(50);not null"`
	Senderid        string    `gorm:"type:varchar(50);not null"`
	Sendername      string    `gorm:"type:varchar(32);not null"`
	Receiverid      string    `gorm:"type:varchar(50);not null"`
	Title           string    `gorm:"type:varchar(160);not null"`
	Content         string    `gorm:"not null"`
	Gold            int       `gorm:"type:int(11);default:0;not null"`
	Cash            int       `gorm:"type:int(11);default:0;not null"`
	Currencies      string    `gorm:"not null"`
	Items           string    `gorm:"not null"`
	Isread          byte      `gorm:"default:0;not null"`
	Readtimestamp   time.Time `gorm:"type:timestamp;null;default:null"`
	Isclaim         byte      `gorm:"default:0;not null"`
	Claimtimestamp  time.Time `gorm:"type:timestamp;null;default:null"`
	Isdelete        byte      `gorm:"default:0;not null"`
	Deletetimestamp time.Time `gorm:"type:timestamp;null;default:null"`
	Senttimestamp   time.Time `gorm:"type:timestamp;not null;default:current_timestamp()"`
}

func (m *LogMailTopup) TableName() string {
	return "mail"
}

type Bankingbonus struct {
	Banking string `gorm:"type:varchar(100);not null"`
	Channel string `gorm:"type:varchar(100);not null"`
	Bonus   int    `gorm:"type:int(11);not null;default:0"`
	gorm.Model
}

type Topuprecheck struct {
	UserId    string
	Txid      string
	Orderid   string
	Status    string
	Detail    string
	Channel   string
	Price     string
	Bonus     string
	Sig       string
	IPAddress string

	gorm.Model
}
