package main

import (
	"fmt"
	"test-pp-back/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	routes.RegisterRoutes(server)

	err := server.Run(":8080")
	if err != nil {
		fmt.Println("Failed to start server")
		panic(err)
	}
}
