package repository

import (
	"flow/db"
	"flow/enum"
	"flow/logger"
	"flow/model"
	"flow/utility"
)

type FieldRepository interface {
	FetchFieldFromFieldVersion(completeFieldVersionNumberList map[int]bool) []model.FieldVersion
	FetchFieldVersions(fieldStatus enum.Status, fieldVersionNumbers []int) []model.FieldVersion
}

type FieldRepositoryImpl struct {
	MapUtil   utility.MapUtil
	DBService db.DBService
}

func NewFieldRepository() *FieldRepositoryImpl {
	repo := &FieldRepositoryImpl{
		MapUtil:utility.MapUtil{},
		DBService:db.DBService{},
	}
	return repo
}

func (f *FieldRepositoryImpl) FetchFieldFromFieldVersion(completeFieldVersionNumberList map[int]bool) []model.FieldVersion {
	methodName := "FetchFieldFromFieldVersion"
	logger.SugarLogger.Info(methodName, " Fetching the field data with join on field versions")
	var fieldVersions []model.FieldVersion
	dbConnection := f.DBService.GetDB()
	if dbConnection == nil {
		return nil
	}
	dbConnection.Joins("JOIN field ON field.id = field_version.field_id ").Where("field_version.id in (?) ", f.MapUtil.GetKeyListFromKeyValueMap(completeFieldVersionNumberList)).Find(&fieldVersions)

	return fieldVersions
}

func (f *FieldRepositoryImpl) FetchFieldVersions(fieldStatus enum.Status, fieldVersionNumbers []int) []model.FieldVersion {
	methodName := "FetchFieldVersions"
	logger.SugarLogger.Info(methodName, " Fetching the field data with join on field versions with field status %s", fieldStatus.String())
	var fieldVersions []model.FieldVersion
	dbConnection := f.DBService.GetDB()
	if dbConnection == nil {
		return fieldVersions
	}
	dbConnection.Joins("JOIN field ON field.id = field_version.field_id and field.status = ? and field.deleted_on is NULL", fieldStatus).Where("field_version.id in (?) and field_version.deleted_on is NULL", fieldVersionNumbers).Find(&fieldVersions)
	return fieldVersions
}
