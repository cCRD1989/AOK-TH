package aokmodel

import "time"

// `id` varchar(50) NOT NULL,
// `username` varchar(32) NOT NULL,
// `password` varchar(72) NOT NULL,
// `gold` int(11) NOT NULL DEFAULT 0,
// `cash` int(11) NOT NULL DEFAULT 0,
// `email` varchar(50) NOT NULL DEFAULT '',
// `isEmailVerified` tinyint(1) NOT NULL DEFAULT 0,
// `authType` tinyint(3) unsigned NOT NULL DEFAULT 1,
// `accessToken` varchar(36) NOT NULL DEFAULT '',
// `userLevel` tinyint(3) unsigned NOT NULL DEFAULT 0,
// `unbanTime` bigint(20) NOT NULL DEFAULT 0,
// `createAt` timestamp NOT NULL DEFAULT current_timestamp(),
// `updateAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
// PRIMARY KEY (`id`),
// UNIQUE KEY `username` (`username`)

type Userlogin struct {
	Id              string `gorm:"type:varchar(50);not null;primaryKey"`
	Username        string `gorm:"type:varchar(32);not null;uniqueIndex"`
	Password        string `gorm:"type:varchar(72);not null"`
	Gold            int
	Cash            int
	Email           string `gorm:"type:varchar(50);not null;default:"`
	IsEmailVerified int
	AuthType        int
	AccessToken     string `gorm:"type:varchar(36);not null;default:"`
	UserLevel       int
	UnbanTime       int
	CreateAt        time.Time
	UpdateAt        time.Time
}

func (n *Userlogin) TableName() string {
	return "userlogin"
}
