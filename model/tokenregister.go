package model

import "gorm.io/gorm"

// Status
// 0 ยังไม่ได้รับของ
// 1 รับของไปแล้ว
// 2 หมดอายุไปแล้ว
type LogTokenregister struct {
	Username string `gorm:"type:varchar(32);uniqueIndex;not null"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Tokenid  string `gorm:"type:varchar(32);uniqueIndex;not null"`
	Status   int    `gorm:"not null"`
	gorm.Model
}
