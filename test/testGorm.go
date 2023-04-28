package main

import (
	"fmt"
	"ginchat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// 建表
func main() {
	db, err := gorm.Open(mysql.Open("root:anwma@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 迁移 schema
	db.AutoMigrate(&models.UserBasic{})

	// Create
	user := &models.UserBasic{}
	user.Name = "Anwma"
	db.Create(user)

	// Read
	fmt.Println("------------------------")
	fmt.Println(db.First(user, 1))
	//Output: &{0xc000114480 <nil> 1 0xc00028a1c0 0}

	// Update
	db.Model(user).Update("Password", "anwmajeff")
}
