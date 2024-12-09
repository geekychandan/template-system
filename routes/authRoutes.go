package routes

import (
	"template-system/controllers"

	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", controllers.Register)
		authGroup.POST("/login", controllers.Login)
		authGroup.POST("/reset-password", controllers.RequestPasswordReset)
		authGroup.POST("/reset-password/confirm", controllers.ResetPassword)
	}
}
