package routes

import (
	"template-system/controllers"
	"template-system/middleware"

	"github.com/gin-gonic/gin"
)

func InitDocumentRoutes(r *gin.Engine) {
	documentGroup := r.Group("api/documents", middleware.AuthMiddleware())
	{
		documentGroup.POST("/:id/generate", controllers.GenerateDocument)
		documentGroup.GET("/", controllers.GetGeneratedDocuments)
		documentGroup.GET("/:id/download", controllers.DownloadDocument)
	}
}
