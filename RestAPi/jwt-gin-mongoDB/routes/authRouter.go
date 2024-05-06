package routes

import (
	controller "github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/jwt-gin-mongoDB/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes, _ gin.Engine) {
	incomingRoutes.POST("/users/signup", controller.Signup())
	incomingRoutes.POST("/users/login", controller.Login())

}
