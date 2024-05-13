package routes

import (
	"github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/sample/controller"
	"github.com/gin-gonic/gin"
)

func UserRoutes(routes *gin.Engine) {

	routes.GET("/", controller.GetUsers)
	routes.POST("/", controller.CreateUser)
	routes.PUT("/:id", controller.UpdateUser)
	routes.DELETE("/:id", controller.DeleteUser)
	routes.GET("/:id", controller.GetUserById)
}
