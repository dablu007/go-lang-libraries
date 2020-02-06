package service

import (
	"encoding/json"
	"flow/db"
	"flow/db/repository"
	"flow/enum"
	"flow/logger"
	"flow/model"
	"flow/model/response_dto"
	"flow/utility"
)

type JourneyServiceUtility interface {
	FetchAllJourneysFromDB(flowContext model.FlowContext) []model.Journey
	FetchJourneyByIdFromDB(flowExternalId string) model.Journey
	GetModuleSectionAndFieldVersionsAndActiveVersionNumberList(journeys ...model.Journey) (
		moduleVersionsMap map[int]model.ModuleVersion, sectionVersionsMap map[int]model.SectionVersion, fieldVersionsMap map[int]model.FieldVersion,
		completeModuleVersionNumberList map[int]bool, completeSectionVersionNumberList map[int]bool, completeFieldVersionNumberList map[int]bool)
	ConstructJourneysResponse(journeys []model.Journey, moduleVersionsMap map[int]model.ModuleVersion,
		sectionVersionsMap map[int]model.SectionVersion, fieldVersionsMap map[int]model.FieldVersion,
		completeModuleVersionNumberList map[int]bool, completeSectionVersionNumberList map[int]bool,
		completeFieldVersionNumberList map[int]bool) response_dto.JourneyResponsesDto
	ConstructFlowResponseWithModuleFieldSection(journey model.Journey,
		completeModuleVersionNumberList map[int]bool, moduleVersionsMap map[int]model.ModuleVersion,
		completeSectionVersionNumberList map[int]bool, sectionVersionsMap map[int]model.SectionVersion,
		completeFieldVersionNumberList map[int]bool, fieldVersionsMap map[int]model.FieldVersion) response_dto.JourneyResponseDto
}

type JourneyServiceUtil struct {
	MapUtil           utility.MapUtility
	DBService         *db.DBService
	JourneyRepository repository.JourneyRepository
	FieldRepository   repository.FieldRepository
	ModuleRepository  repository.ModuleRepository
	SectionRepository repository.SectionRepository
}

func NewJourneyServiceUtil(util utility.MapUtility, dBService *db.DBService, journeyRepository repository.JourneyRepository,
	fieldRepository repository.FieldRepository, moduleRepository repository.ModuleRepository,
	sectionRepository repository.SectionRepository) *JourneyServiceUtil {
	service := &JourneyServiceUtil{
		MapUtil:           util,
		DBService:         dBService,
		JourneyRepository: journeyRepository,
		FieldRepository:   fieldRepository,
		ModuleRepository:  moduleRepository,
		SectionRepository: sectionRepository,
	}
	return service
}

func (f JourneyServiceUtil) FetchAllJourneysFromDB(flowContext model.FlowContext) []model.Journey {
	return f.JourneyRepository.FindActiveJourneysByJourneyContext(flowContext.MerchantId, flowContext.TenantId, flowContext.ChannelId)
}

func (f JourneyServiceUtil) FetchJourneyByIdFromDB(flowExternalId string) model.Journey {
	methodName := "FetchJourneyByIdFromDB:"
	logger.SugarLogger.Info(methodName, " Fetching flows from db for journey id ", flowExternalId)
	var journey model.Journey
	journey = f.JourneyRepository.FindByExternalId(flowExternalId)
	return journey
}

func (f JourneyServiceUtil) GetModuleSectionAndFieldVersionsAndActiveVersionNumberList(journeys ...model.Journey) (
	moduleVersionsMap map[int]model.ModuleVersion, sectionVersionsMap map[int]model.SectionVersion, fieldVersionsMap map[int]model.FieldVersion,
	completeModuleVersionNumberList map[int]bool, completeSectionVersionNumberList map[int]bool, completeFieldVersionNumberList map[int]bool) {
	methodName := "GetModuleSectionAndFieldVersionsAndActiveVersionNumberList"
	// logger.SugarLogger.Info(methodName, "fetching the response for flow")

	completeModuleVersionNumberList = make(map[int]bool)
	var moduleVersions []model.ModuleVersion
	moduleVersionsMap = make(map[int]model.ModuleVersion)

	completeSectionVersionNumberList = make(map[int]bool)
	var sectionVersions []model.SectionVersion
	sectionVersionsMap = make(map[int]model.SectionVersion)

	completeFieldVersionNumberList = make(map[int]bool)
	var fieldVersions []model.FieldVersion
	fieldVersionsMap = make(map[int]model.FieldVersion)

	var moduleVersionNumbersListForFlow []int
	for _, journey := range journeys {
		var versionNumbers []int
		json.Unmarshal([]byte(journey.ModuleVersions), &versionNumbers)
		moduleVersionNumbersListForFlow = append(moduleVersionNumbersListForFlow, versionNumbers...)
	}
	logger.SugarLogger.Info(methodName, "list of modules ", moduleVersionNumbersListForFlow)
	for _, num := range moduleVersionNumbersListForFlow {
		if completeModuleVersionNumberList[num] == false {
			completeModuleVersionNumberList[num] = true
		}
	}

	moduleVersions = f.ModuleRepository.FetchModuleVersions(enum.Active, f.MapUtil.GetKeyListFromKeyValueMap(completeModuleVersionNumberList))

	var sectionNumberList []int
	for _, mv := range moduleVersions {
		moduleVersionsMap[mv.Id] = mv
		var sectionNumbers []int
		json.Unmarshal([]byte(mv.SectionVersions), &sectionNumbers)
		sectionNumberList = append(sectionNumberList, sectionNumbers...)
	}
	logger.SugarLogger.Info(methodName, "list of sections ", sectionNumberList)
	for _, num := range sectionNumberList {
		if completeSectionVersionNumberList[num] == false {
			completeSectionVersionNumberList[num] = true
		}
	}

	sectionVersions = f.SectionRepository.FetchSectionVersions(enum.Active, f.MapUtil.GetKeyListFromKeyValueMap(completeSectionVersionNumberList))

	var fieldNumbersList []int
	for _, sv := range sectionVersions {
		sectionVersionsMap[sv.Id] = sv
		var fieldNumbers []int
		json.Unmarshal([]byte(sv.FieldVersions), &fieldNumbers)
		fieldNumbersList = append(fieldNumbersList, fieldNumbers...)
	}

	logger.SugarLogger.Info(methodName, "list of fields ", fieldNumbersList)
	for _, num := range fieldNumbersList {
		if completeFieldVersionNumberList[num] == false {
			completeFieldVersionNumberList[num] = true
		}
	}

	fieldVersions = f.FieldRepository.FetchFieldVersions(enum.Active, f.MapUtil.GetKeyListFromKeyValueMap(completeFieldVersionNumberList))

	for _, fv := range fieldVersions {
		fieldVersionsMap[fv.Id] = fv
	}

	return moduleVersionsMap, sectionVersionsMap, fieldVersionsMap,
		completeModuleVersionNumberList, completeSectionVersionNumberList, completeFieldVersionNumberList
}

func (f JourneyServiceUtil) ConstructJourneysResponse(journeys []model.Journey, moduleVersionsMap map[int]model.ModuleVersion,
	sectionVersionsMap map[int]model.SectionVersion, fieldVersionsMap map[int]model.FieldVersion,
	completeModuleVersionNumberList map[int]bool, completeSectionVersionNumberList map[int]bool,
	completeFieldVersionNumberList map[int]bool) response_dto.JourneyResponsesDto {
	var response response_dto.JourneyResponsesDto
	for _, journey := range journeys {
		response.JourneyResponses = append(response.JourneyResponses, f.ConstructFlowResponseWithModuleFieldSection(journey,
			completeModuleVersionNumberList, moduleVersionsMap,
			completeSectionVersionNumberList, sectionVersionsMap,
			completeFieldVersionNumberList, fieldVersionsMap))
	}
	return response
}

func (f JourneyServiceUtil) ConstructFlowResponseWithModuleFieldSection(journey model.Journey,
	completeModuleVersionNumberList map[int]bool, moduleVersionsMap map[int]model.ModuleVersion,
	completeSectionVersionNumberList map[int]bool, sectionVersionsMap map[int]model.SectionVersion,
	completeFieldVersionNumberList map[int]bool, fieldVersionsMap map[int]model.FieldVersion) response_dto.JourneyResponseDto {
	methodName := "ConstructFlowResponseWithModuleFieldSection"
	logger.SugarLogger.Info(methodName, "fetching the response for journey data")

	journeyResponseDto := response_dto.JourneyResponseDto{
		Name:       journey.Name,
		ExternalId: journey.ExternalId,
		Version:    journey.Version,
		Type:       journey.Type.String()}
	var moduleVersionNumberList []int
	json.Unmarshal([]byte(journey.ModuleVersions), &moduleVersionNumberList)
	for _, mvn := range moduleVersionNumberList {
		if completeModuleVersionNumberList[mvn] == true {
			moduleVersion := moduleVersionsMap[mvn]
			if (model.ModuleVersion{}) == moduleVersion {
				continue
			}
			journeyResponseDto.Modules = append(journeyResponseDto.Modules, f.getModuleVersionResponseDto(moduleVersion,
				completeSectionVersionNumberList, sectionVersionsMap, completeFieldVersionNumberList, fieldVersionsMap))
		}
	}

	logger.SugarLogger.Info(methodName, "Returning the response for journey data => ", journeyResponseDto)
	return journeyResponseDto
}

func (f JourneyServiceUtil) getModuleVersionResponseDto(moduleVersion model.ModuleVersion, completeSectionVersionNumberList map[int]bool,
	sectionVersionsMap map[int]model.SectionVersion, completeFieldVersionNumberList map[int]bool,
	fieldVersionsMap map[int]model.FieldVersion) response_dto.ModuleVersionResponseDto {
	moduleVersionResponseDto := response_dto.ModuleVersionResponseDto{
		Name:       moduleVersion.Name,
		ExternalId: moduleVersion.ExternalId,
		Version:    moduleVersion.Version}
	json.Unmarshal([]byte(moduleVersion.Properties), &moduleVersionResponseDto.Properties)
	var sectionVersionNumberList []int
	json.Unmarshal([]byte(moduleVersion.SectionVersions), &sectionVersionNumberList)
	for _, svn := range sectionVersionNumberList {
		if completeSectionVersionNumberList[svn] == true {
			sectionVersion := sectionVersionsMap[svn]
			if (model.SectionVersion{}) == sectionVersion {
				continue
			}
			moduleVersionResponseDto.Sections = append(moduleVersionResponseDto.Sections,
				f.getSectionVersionResponseDto(sectionVersion, completeFieldVersionNumberList, fieldVersionsMap))
		}
	}
	return moduleVersionResponseDto
}

func (f JourneyServiceUtil) getSectionVersionResponseDto(sectionVersion model.SectionVersion, completeFieldVersionNumberList map[int]bool, fieldVersionsMap map[int]model.FieldVersion) response_dto.SectionVersionsResponseDto {
	sectionVersionResponseDto := response_dto.SectionVersionsResponseDto{
		Name:       sectionVersion.Name,
		ExternalId: sectionVersion.ExternalId,
		Version:    sectionVersion.Version,
		IsVisible:  sectionVersion.IsVisible}
	json.Unmarshal([]byte(sectionVersion.Properties), &sectionVersionResponseDto.Properties)
	var fieldVersionNumberList []int
	json.Unmarshal([]byte(sectionVersion.FieldVersions), &fieldVersionNumberList)
	for _, fvn := range fieldVersionNumberList {
		if completeFieldVersionNumberList[fvn] == true {
			fieldVersion := fieldVersionsMap[fvn]
			if (model.FieldVersion{}) == fieldVersion {
				continue
			}
			sectionVersionResponseDto.Fields = append(sectionVersionResponseDto.Fields, f.getFieldVersionResponseDto(fieldVersion))
		}
	}
	return sectionVersionResponseDto
}

func (f JourneyServiceUtil) getFieldVersionResponseDto(fieldVersion model.FieldVersion) response_dto.FieldVersionsResponseDto {
	fieldVersionResponseDto := response_dto.FieldVersionsResponseDto{
		Name:       fieldVersion.Name,
		ExternalId: fieldVersion.ExternalId,
		IsVisible:  fieldVersion.IsVisible,
		Version:    fieldVersion.Version}
	json.Unmarshal([]byte(fieldVersion.Properties), &fieldVersionResponseDto.Properties)
	return fieldVersionResponseDto
}

func (f JourneyServiceUtil) ConstructFlowResponseAsList(journey model.Journey, moduleVersions map[int]model.ModuleVersion,
	sectionVersions map[int]model.SectionVersion, fieldVersions map[int]model.FieldVersion) response_dto.JourneyResponseDtoList {
	methodName := "ConstructFlowResponseAsList"
	logger.SugarLogger.Info(methodName, "fetching the response for journey data")

	journeyResponseDtoList := response_dto.JourneyResponseDtoList{
		Name:       journey.Name,
		ExternalId: journey.ExternalId,
		Version:    journey.Version,
		Type:       journey.Type.String()}
	var moduleVersionList []response_dto.ResponseDTO
	var sectionVersionList []response_dto.ResponseDTO
	var fieldVersionsList []response_dto.ResponseDTO
	for _, value := range moduleVersions {
		dto := response_dto.ResponseDTO{
			Name:       value.Name,
			Version:    value.Version,
			ExternalId: value.ExternalId,
		}
		moduleVersionList = append(moduleVersionList,dto)
	}

	for _, value := range sectionVersions {
		dto := response_dto.ResponseDTO{
			Name:       value.Name,
			Version:    value.Version,
			ExternalId: value.ExternalId,
		}
		sectionVersionList = append(sectionVersionList,dto)
	}
	for _, value := range fieldVersions {
		dto := response_dto.ResponseDTO{
			Name:       value.Name,
			Version:    value.Version,
			ExternalId: value.ExternalId,
		}
		fieldVersionsList = append(fieldVersionsList,dto)
	}
	journeyResponseDtoList.Modules = moduleVersionList
	journeyResponseDtoList.Sections = sectionVersionList
	journeyResponseDtoList.Fields = fieldVersionsList

	return journeyResponseDtoList
}