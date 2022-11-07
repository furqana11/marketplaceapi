package main

import (
	"log"
	// "os"
	"github.com/gin-gonic/gin"

	config "github.com/FarisTF/marketplace_api/config"
	routes "github.com/FarisTF/marketplace_api/routes"
)

func main() {
	// Connect DB
	config.Connect()

	// Init Router
	router := gin.Default()
	
	// Route Handlers / Endpoints
	routes.Routes(router)
	log.Fatal(router.Run(":4747"))
}
