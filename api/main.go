package api

import "github.com/gin-gonic/gin"

func InitializeRoutes(router *gin.Engine) {
	rincon := router.Group("/rincon", func(c *gin.Context) {})
	rincon.GET("/ping", Ping)
	rincon.GET("/services", GetAllServices)
	rincon.GET("/services/:name", GetService)
	rincon.GET("/services/:name/routes", GetRoutesForService)
	rincon.POST("/services", CreateService)
	rincon.GET("/routes", GetAllRoutes)
	rincon.GET("/routes/:id", GetRoute)
	rincon.POST("/routes", CreateRoute)
}
