package models

import "gorm.io/gorm"

type Genre struct {
	gorm.Model
    ID            uint   `gorm:"primary_key"`
	Name          string `gorm:"unique;not null"`
}

func (g Genre) String() string {
    return g.Name
}
