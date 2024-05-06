package controller

import (
	database "github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/jwt-gin-mongoDB/databases"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword()

func VerifyPassword()

func Signup()

func Login()

func GetUser()

func GetUsers()
