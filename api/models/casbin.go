package models

type CasbinRule struct {
	PType string `gorm:"type:varchar(100);DEFAULT:NULL"`
	V0    string `gorm:"type:varchar(100);DEFAULT:NULL"`
	V1    string `gorm:"type:varchar(100);DEFAULT:NULL"`
	V2    string `gorm:"type:varchar(100);DEFAULT:NULL"`
	V3    string `gorm:"type:varchar(100);DEFAULT:NULL"`
	V4    string `gorm:"type:varchar(100);DEFAULT:NULL"`
	V5    string `gorm:"type:varchar(100);DEFAULT:NULL"`
}

func (CasbinRule) TableName() string {
	return "casbin_rule"
}
