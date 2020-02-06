package server

import (
	"flow/auth"
	"flow/controller"
	"flow/db"
	"flow/db/repository"
	"flow/service"
	"flow/utility"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
	"github.com/spf13/viper"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	health := new(controller.HealthController)
	cacheController := new(controller.CacheController)
	appName := "Flow CJM"
	cfg := newrelic.NewConfig(appName, viper.GetString("newrelic.licensekey"))
	cfg.Logger = newrelic.NewLogger(os.Stdout)
	app, err := newrelic.NewApplication(cfg)
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}
	router.Use(nrgin.Middleware(app))
	router.GET("journey-definition/health", health.Status)
	router.Use(auth.AuthMiddleware())
	router.GET("journey-definition/refresh", cacheController.DeleteCacheEntry())

	v1 := router.Group("journey-definition/v1")
	{
		group := v1.Group("/")
		{
			//flowController := new(controller.JourneyController)
			validator := utility.NewRequestValidator()
			mapUtil := utility.NewMapUtil()
			dbService := new(db.DBService)
			journeyRepo := repository.NewJourneyRepository()
			fieldRepo := repository.NewFieldRepository()
			moduleRepo := repository.NewModuleRepository()
			sectionRepo := repository.NewSectionRepository()
			journeyServiceUtil := service.NewJourneyServiceUtil(mapUtil, dbService, journeyRepo, fieldRepo, moduleRepo, sectionRepo)
			service := service.NewJourneyService(journeyServiceUtil, validator)
			flowController := controller.NewJourneyController(service, validator)
			group.GET("journeys", flowController.GetJourneys())
			group.GET("journeys/:journeyId", flowController.GetJourneyById())

		}
	}
	return router
}
