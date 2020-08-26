package repository

import (
	"github.com/dablu007/go-lang-libraries/db"
	"github.com/dablu007/go-lang-libraries/models"
)

type EmployeeRepository struct {
	DBService db.DBService
}

func (f EmployeeRepository) FetchEmployeeDetails(employeeId string) models.Employee {
	var employee models.Employee
	dbConnection := f.DBService.GetDB()
	if dbConnection == nil {
		return employee
	}
	dbConnection.Where("employee.employee_id = ?", employeeId).Find(&employee)
	return employee
}

