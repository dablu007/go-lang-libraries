package service

import (
	"fmt"
	"github.com/dablu007/go-lang-libraries/apiclient"
	"github.com/dablu007/go-lang-libraries/db/repository"
	"github.com/dablu007/go-lang-libraries/logger"
	"github.com/dablu007/go-lang-libraries/models"
)

type EmployeeService struct {
	EmployeeRepo repository.EmployeeRepository

}

func (e EmployeeService) FetchDetails(employeeId string) models.Employee {
	logger.SugarLogger.Infow("Inside employee service fetching details")
	return e.EmployeeRepo.FetchEmployeeDetails(employeeId)
}

func (e EmployeeService) FetchdDetailsFromOpenAPI() string {
	fmt.Println("Fetching details from an open api")
	url := "http://dummy.restapiexample.com/api/v1/employees"
	request, _ := apiclient.CreateJSONRequest("GET", url, "", nil)

	_, response, _ := apiclient.RestExecute(request)
	fmt.Println("this is the response ", response)
	return response
}
