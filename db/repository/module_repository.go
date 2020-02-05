package repository

import (
	"flow/db"
	"flow/logger"
	"flow/model"
	"flow/utility"
)

type ModuleRepository struct {
	MapUtil utility.MapUtil
	DBService db.DBService
}

func (f ModuleRepository) FetchModuleFromModuleVersion(completeModuleVersionNumberList map[int]bool) []model.ModuleVersion {
	methodName := "FetchFieldFromFieldVersion"
	logger.SugarLogger.Info(methodName, " Fetching the module data with join on field versions")
	var moduleVersions []model.ModuleVersion
	dbConnection := f.DBService.GetDB()
	if dbConnection == nil {
		return moduleVersions
	}
	dbConnection.Debug().Joins("JOIN module ON module.id = module_version.module_id ").Where("module_version.id in (?) ", f.MapUtil.GetKeyListFromKeyValueMap(completeModuleVersionNumberList)).Find(&moduleVersions)
	return moduleVersions
}
