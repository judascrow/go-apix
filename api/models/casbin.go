package models

type CabinRule struct {
	PType string `gorm:"type:varchar(100)"`
	V0    string `gorm:"type:varchar(100)"`
	V1    string `gorm:"type:varchar(100)"`
	V2    string `gorm:"type:varchar(100)"`
	V3    string `gorm:"type:varchar(100)"`
	V4    string `gorm:"type:varchar(100)"`
	V5    string `gorm:"type:varchar(100)"`
}

func (CabinRule) TableName() string {
	return "cabin_rule"
}
