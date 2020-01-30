package server

import (
	"flow/auth"
	"flow/controller"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
	"github.com/spf13/viper"
	"os"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	health := new(controller.HealthController)
	cacheController := new(controller.CacheController)
	appName := "Flow CJM"
	cfg := newrelic.NewConfig(appName, viper.GetString("newrelic.licensekey"))
	cfg.Logger = newrelic.NewDebugLogger(os.Stdout)
	app, err := newrelic.NewApplication(cfg)
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}
	router.Use(nrgin.Middleware(app))
	router.GET("rest/flows/health", health.Status)
	router.Use(auth.AuthMiddleware())
	router.GET("rest/flows/refresh", cacheController.DeleteCacheEntry())

	v1 := router.Group("rest/v1")
	{
		group := v1.Group("/")
		{
			flowController := new(controller.FlowController)
			group.GET("flows", flowController.GetFlows())
			group.GET("flows/:flowId", flowController.GetFlowById())

		}
	}
	return router
}
