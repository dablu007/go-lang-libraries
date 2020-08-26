package controller

import (
	"fmt"
	"github.com/dablu007/go-lang-libraries/logger"
	"github.com/dablu007/go-lang-libraries/models"
	"github.com/dablu007/go-lang-libraries/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EmployeeController struct {
	EmployeeService service.EmployeeService
}
func (e EmployeeController) FetchEmployeeDetails(c *gin.Context) {
	logger.SugarLogger.Info("Inside the employee controller")
	fmt.Print("Inside the employee controller")
	var data models.Employee = e.EmployeeService.FetchDetails(c.Param("id"))
	fmt.Println("this is the data recieced ", data)
	c.String(http.StatusOK, data.Name)
}

func (e EmployeeController) FetchDetails(c *gin.Context){
	logger.SugarLogger.Infow("Inside fetch details of employee controller")
	fmt.Println("Inside fetch details of employee controller")
	var data = e.EmployeeService.FetchdDetailsFromOpenAPI()
	c.JSON(http.StatusOK, data)
}
