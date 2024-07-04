package main

import (
	"fmt"

	"github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/Gin-postgres-docker/config"
	"github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/Gin-postgres-docker/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.New()
	config.Connect()
	routes.UserRoute(router)

	router.Run(":8080")

	result := Calculate(3)
	fmt.Println(result)
}

func Calculate(x int) (result int) {
	result = x + 2
	return result
}
