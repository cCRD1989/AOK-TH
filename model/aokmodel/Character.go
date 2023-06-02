package aokmodel

type Character struct {
	Id            string
	Userid        string
	Dataid        int
	Charactername string `gorm:"type:varchar(32)"`
	Level         int    `gorm:"type:int(11)"`
	Factionid     int
	Currenthp     int
	Currentmp     int
	Guildid       int
	Guildids      string
}

func (c *Character) TableName() string {
	return "characters"
}

type Guild struct {
	Id        int    `gorm:"uniqueIndex;type:int(11);not null"`
	Guildname string `gorm:"type:varchar(32)"`
	Leaderid  string `gorm:"type:varchar(50)"`
}

func (c *Guild) TableName() string {
	return "guild"
}
