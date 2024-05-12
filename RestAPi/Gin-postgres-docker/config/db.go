package config

import (
	"github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/Gin-postgres-docker/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	//dsn := "host=localhost user=postgres password=1234567 dbname=gorm1 port=5432"
	db, err := gorm.Open(postgres.Open("postgres://postgres:1234567@localhost:5432/gorm1"), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	DB = db
}
