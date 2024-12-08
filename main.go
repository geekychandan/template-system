package main

import (
	"template-system/config"
	"template-system/routes"
	"template-system/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	utils.InitDB()
	utils.InitCache()

	r := gin.Default()

	routes.InitAuthRoutes(r)
	routes.InitTemplateRoutes(r)
	routes.InitDocumentRoutes(r)

	r.Run()
}
