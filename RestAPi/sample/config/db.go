package config

import (
	"github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/sample/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connection() {

	db, err := gorm.Open(postgres.Open("postgres://postgres:1234567@localhost:5432/gosample"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})

	DB = db
}
