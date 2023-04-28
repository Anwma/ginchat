package models

import "gorm.io/gorm"

//群信息

type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string
	Desc    string
	Type    int
}

func (table *GroupBasic) TableName() string {
	return "group_basic"

}
