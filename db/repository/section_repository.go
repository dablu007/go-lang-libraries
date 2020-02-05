package repository

import (
	"flow/db"
	"flow/logger"
	"flow/model"
	"flow/utility"
)

type SectionRepository struct {
	MapUtil utility.MapUtil
	DBService db.DBService
}

func (f SectionRepository) FetchSectionFromSectionVersions(completeSectionVersionNumberList map[int]bool) []model.SectionVersion {
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

