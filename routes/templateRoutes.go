package routes

import (
	"template-system/controllers"
	"template-system/middleware"

	"github.com/gin-gonic/gin"
)

func InitTemplateRoutes(r *gin.Engine) {
	templateGroup := r.Group("api/templates", middleware.AuthMiddleware())
	{
		templateGroup.POST("/upload", controllers.UploadTemplate)
		templateGroup.GET("/", controllers.GetTemplates)
		templateGroup.GET("/:id/placeholders", controllers.GetPlaceholders)
	}
}
