package models

type Employee struct {
	Id                int                    `gorm:"primary_key; AUTO_INCREMENT; column:id"`
	EmployeeId        string                 `gorm:"column:employee_id"`
	Name              string                 `gorm:"column:name"`
	Address	          string                  `gorm:"column:address"`
}
