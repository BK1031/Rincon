package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rincon/model"
	"rincon/service"
)

func GetAllRoutes(c *gin.Context) {
	result := service.GetAllRoutes()
	c.JSON(http.StatusOK, result)
}

func GetRoute(c *gin.Context) {
	result := service.GetRouteByID(c.Param("id"))
	if result.Route == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "No route " + c.Param("id") + " found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func GetRoutesForService(c *gin.Context) {
	result := service.GetRoutesByServiceName(c.Param("name"))
	c.JSON(http.StatusOK, result)
}

func CreateRoute(c *gin.Context) {
	var input model.Route
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := service.CreateRoute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, service.GetRouteByID(input.Route))
}
