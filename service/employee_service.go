package service

import (
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
