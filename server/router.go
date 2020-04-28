package server

import (
	"fmt"
	"github.com/dablu007/go-lang-libraries/auth"
	"github.com/dablu007/go-lang-libraries/controller"
	"os"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
	"github.com/spf13/viper"
)

func NewRouter() *gin.Engine {
	//router := gin.New()
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	health := new(controller.HealthController)
	appName := "Flow CJM"
	cfg := newrelic.NewConfig(appName, viper.GetString("newrelic.licensekey"))
	cfg.Logger = newrelic.NewLogger(os.Stdout)
	cfg.DistributedTracer.Enabled = true
	app, err := newrelic.NewApplication(cfg)
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}
	router.Use(nrgin.Middleware(app))
	router.Use(auth.AuthMiddleware())
	v1 := router.Group("go-lang-libraries/v1")
	{
		group := v1.Group("/")
		{
			group.GET("health", health.Status)

		}
	}
	return router
}
