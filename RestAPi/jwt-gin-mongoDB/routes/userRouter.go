package routes

import (
	controller "github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/jwt-gin-mongoDB/controllers"
	"github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/jwt-gin-mongoDB/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:userid", controller.GetUser())
}
