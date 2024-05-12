package models

type User struct {
	Id       int    `json:"id" gorm:"primary_key"`
	Name     string `json:"name" validation:"required"`
	Email    string `json:"email" validation:"required"`
	Password string `json:"password" validation:"required"`
}
