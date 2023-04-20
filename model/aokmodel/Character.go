package aokmodel

type Character struct {
	CharacterName string `gorm:"type:varchar(32)"`
	Level         int    `gorm:"type:int(11)"`
}

func (n *Character) TableName() string {
	return "characters"
}
