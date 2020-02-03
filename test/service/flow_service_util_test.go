package service

import (
	"encoding/json"
	"flow/config"
	"flow/logger"
	"flow/model"
	"flow/service"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"testing"
	"time"
)

func init() {
	service := "finbox-integration"
	environment := "dev"
	config.Init(service, environment)
	logger.InitLogger()
}

func TestFetchModuleData(t *testing.T){
	var mockController = gomock.NewController(t);
	defer mockController.Finish()

	var db *gorm.DB
	_, mock, err := sqlmock.NewWithDSN("sqlmock_db_0")
	if err != nil {
		panic("Got an unexpected error.")
	}

	db, err = gorm.Open("postgresmock", "sqlmock_db_0")
	if err != nil {
		panic("Got an unexpected error.")
	}
	defer db.Close()
	var moduleVersions []model.ModuleVersion
	moduleVersions = append(moduleVersions, model.ModuleVersion{
		Id:              0,
		Name:            "",
		ModuleId:        0,
		ExternalId:      uuid.UUID{},
		Version:         "",
		CreatedOn:       time.Time{},
		DeletedOn:       time.Time{},
		Properties:      "",
		SectionVersions: "",
	})
	b, err := json.Marshal(moduleVersions)
	if err != nil {
		fmt.Println(err)
		return
	}
	var a []string
	a = append(a, string(b))
	rs := sqlmock.NewRows(a)
	mock.ExpectQuery(`SELECT "module_version".* FROM "module_version" JOIN module ON module.id = module_version.module_id  WHERE (module_version.id in ($) `).
		WithArgs(1).
		WillReturnRows(rs)
	var flowServiceUtil = &service.FlowServiceUtil{

	}
}