package models

import "gorm.io/gorm"

//人员关系

type Contact struct {
	gorm.Model
	OwnerId  int64  //谁的关系信息
	TargetId uint   //对应的谁
	Type     int    //对应的类型 0 1 2
	Desc     string //描述信息
}

func (table *Contact) TableName() string {
	return "contact"

}
