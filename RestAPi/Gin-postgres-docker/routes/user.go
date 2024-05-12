package routes

import (
	"github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/Gin-postgres-docker/controller"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.GET("/", controller.GetUsers)
	router.POST("/", controller.CreateUser)
	router.DELETE("/:id", controller.DeleteUser)
	router.PUT("/:id", controller.UpdateUser)
	router.GET("/:id", controller.GetUserById)
}
