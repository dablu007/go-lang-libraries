package repository

import (
	"flow/db"
	"flow/enum"
	"flow/logger"
	"flow/model"
	"flow/utility"
)

type ModuleRepository interface {
	FetchModuleFromModuleVersion(completeModuleVersionNumberList map[int]bool) []model.ModuleVersion
	FetchModuleVersions(moduleStatus enum.Status, moduleVersionNumbers []int) []model.ModuleVersion
}

type ModuleRepositoryImpl struct {
	MapUtil   utility.MapUtil
	DBService db.DBService
}

func NewModuleRepository() *ModuleRepositoryImpl {
	repo := &ModuleRepositoryImpl{
		MapUtil:utility.MapUtil{},
		DBService: db.DBService{},
	}
	return repo
}

func (f ModuleRepositoryImpl) FetchModuleFromModuleVersion(completeModuleVersionNumberList map[int]bool) []model.ModuleVersion {
	methodName := "FetchModuleFromModuleVersion"
	logger.SugarLogger.Info(methodName, " Fetching the module data with join on module versions")
	var moduleVersions []model.ModuleVersion
	dbConnection := f.DBService.GetDB()
	if dbConnection == nil {
		return moduleVersions
	}
	dbConnection.Joins("JOIN module ON module.id = module_version.module_id ").Where("module_version.id in (?) ", f.MapUtil.GetKeyListFromKeyValueMap(completeModuleVersionNumberList)).Find(&moduleVersions)
	return moduleVersions
}

func (f ModuleRepositoryImpl) FetchModuleVersions(moduleStatus enum.Status, moduleVersionNumbers []int) []model.ModuleVersion {
	methodName := "FetchModuleVersions"
	logger.SugarLogger.Info(methodName, " Fetching the module data with join on module versions with module status %s", moduleStatus.String())
	var moduleVersions []model.ModuleVersion
	dbConnection := f.DBService.GetDB()
	if dbConnection == nil {
		return moduleVersions
	}
	dbConnection.Joins("JOIN module ON module.id = module_version.module_id and module.status = ? and module.deleted_on is NULL", moduleStatus).Where("module_version.id in (?) and module_version.deleted_on is NULL", moduleVersionNumbers).Find(&moduleVersions)
	return moduleVersions
}
