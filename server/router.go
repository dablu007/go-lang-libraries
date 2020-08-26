package server

import (
	"github.com/dablu007/go-lang-libraries/controller"
	"github.com/dablu007/go-lang-libraries/db"
	"github.com/dablu007/go-lang-libraries/db/repository"
	"github.com/dablu007/go-lang-libraries/service"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	//router := gin.New()
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	health := new(controller.HealthController)
	//router.Use(auth.AuthMiddleware())
	dbService := new(db.DBService)
	employeeRepo := &repository.EmployeeRepository{DBService: *dbService}
	employeeService := &service.EmployeeService{EmployeeRepo: *employeeRepo}
	employeeController := &controller.EmployeeController{
		EmployeeService: *employeeService,
	}
	v1 := router.Group("go-lang-libraries/v1")
	{
		group := v1.Group("/")
		{
			group.GET("health", health.Status)
			group.GET("employee-details/:id", employeeController.FetchEmployeeDetails)
		}
	}
	return router
}
