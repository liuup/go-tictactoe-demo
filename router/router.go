package router

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	// ----- 公共接口 -----
	public := r.Group("/api")

	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)

	// ----- 私有接口 -----

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleWare())

	// protected.GET("/user", controllers.CurrentUser)
	protected.GET("/ws", controllers.WsStart)

	return r
}
