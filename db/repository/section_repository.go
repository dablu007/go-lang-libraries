package repository

import (
	"flow/db"
	"flow/enum"
	"flow/logger"
	"flow/model"
	"flow/utility"
)

type SectionRepository interface {
	FetchSectionFromSectionVersions(completeSectionVersionNumberList map[int]bool) []model.SectionVersion
	FetchSectionVersions(sectionStatus enum.Status, sectionVersionNumbers []int) []model.SectionVersion
}

type SectionRepositoryImpl struct {
	MapUtil   utility.MapUtil
	DBService db.DBService
}

func NewSectionRepository() *SectionRepositoryImpl {
	repo := &SectionRepositoryImpl{
		MapUtil: utility.MapUtil{},
		DBService: db.DBService{},
	}
	return repo
}

func (f SectionRepositoryImpl) FetchSectionFromSectionVersions(completeSectionVersionNumberList map[int]bool) []model.SectionVersion {
	methodName := "FetchFieldFromFieldVersion"
	logger.SugarLogger.Info(methodName, " Fetching the field data with join on field versions")
	var sectionVersions []model.SectionVersion
	dbConnection := f.DBService.GetDB()
	if dbConnection == nil {
		return sectionVersions
	}
	dbConnection.Joins("JOIN section ON section.id = section_version.section_id").Where("section_version.id in (?) ", f.MapUtil.GetKeyListFromKeyValueMap(completeSectionVersionNumberList)).Find(&sectionVersions)

	return sectionVersions
}

func (f SectionRepositoryImpl) FetchSectionVersions(sectionStatus enum.Status, sectionVersionNumbers []int) []model.SectionVersion {
	methodName := "FetchSectionVersions"
	logger.SugarLogger.Info(methodName, " Fetching the section data with join on section versions with section status %s", sectionStatus.String())
	var sectionVersions []model.SectionVersion
	dbConnection := f.DBService.GetDB()
	if dbConnection == nil {
		return sectionVersions
	}
	dbConnection.Joins("JOIN section ON section.id = section_version.section_id and section.status = ? and section.deleted_on is NULL", sectionStatus).Where("section_version.id in (?) and section_version.deleted_on is NULL", sectionVersionNumbers).Find(&sectionVersions)
	return sectionVersions
}
