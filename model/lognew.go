package model

import (
	"gorm.io/gorm"
)

// Datatype
// Event
// Other
type LogNews struct {
	Datatype     string
	Author       string
	Subject      string
	Data         string
	Image        string
	Externallink string
	gorm.Model
}

type LogNewsRequest struct {
	Datatype     string `josn:"datatype" binding:"required"`
	Author       string `josn:"author" binding:"required"`
	Subject      string `josn:"subject" binding:"required"`
	Data         string `josn:"data" binding:"required"`
	Image        string `josn:"image" binding:"required"`
	Externallink string `josn:"externallink" binding:"required"`
}
