package model

import "gorm.io/gorm"

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
