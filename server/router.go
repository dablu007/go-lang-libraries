package server

import (
	"flow/auth"
	"flow/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// logger.InitLogger()
	health := new(controller.HealthController)

	router.GET("flow/health", health.Status)
	router.Use(auth.AuthMiddleware())

	v1 := router.Group("flow/api/v1")
	{
		group := v1.Group("/")
		{
			flowController := new(controller.FlowController)
			group.GET("flows", flowController.GetFlows())
		}
	}
	return router
}
