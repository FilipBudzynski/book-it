package models

import "gorm.io/gorm"

type Genre struct {
	gorm.Model
	Name          string `gorm:"unique;not null"`
}

func (g Genre) String() string {
    return g.Name
}
