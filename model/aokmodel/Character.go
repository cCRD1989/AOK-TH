package aokmodel

type Character struct {
	Charactername string `gorm:"type:varchar(32)"`
	Level         int    `gorm:"type:int(11)"`

	//
}

func (c *Character) TableName() string {
	return "characters"
}
