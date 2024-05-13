package main

import (
	"github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/sample/config"
	"github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/sample/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.New()

	config.Connection()
	routes.UserRoutes(router)

	router.Run(":9000")
}
