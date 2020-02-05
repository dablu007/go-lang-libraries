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

type FlowServiceUtil struct {
	MapUtil           utility.MapUtil
	DBService         db.DBService
	FlowRepository    repository.FlowRepository
	FieldRepository   repository.FieldRepository
	ModuleRepository  repository.ModuleRepository
	SectionRepository repository.SectionRepository
}

func (f FlowServiceUtil) FetchAllFlowsFromDB(flowContext model.FlowContext) []model.Flow {
	f.FlowRepository = new(repository.FlowRepositoryImpl)
	return f.FlowRepository.FindActiveFlowsByFlowContext(flowContext.MerchantId, flowContext.TenantId, flowContext.ChannelId)
}

func (f FlowServiceUtil) GetParsedFlowsResponse(flows []model.Flow) response_dto.FlowResponsesDto {
	methodName := "GetParsedFlowsResponse"
	logger.SugarLogger.Info(methodName, "fetching the response for flow")

	completeModuleVersionNumberList := make(map[int]bool)
	var moduleVersions []model.ModuleVersion
	moduleVersionsMap := make(map[int]model.ModuleVersion)

	completeSectionVersionNumberList := make(map[int]bool)
	var sectionVersions []model.SectionVersion
	sectionVersionsMap := make(map[int]model.SectionVersion)

	completeFieldVersionNumberList := make(map[int]bool)
	var fieldVersions []model.FieldVersion
	fieldVersionsMap := make(map[int]model.FieldVersion)

	var response response_dto.FlowResponsesDto
	var versionNumbersList []int
	for _, flow := range flows {
		var versionNumbers []int
		json.Unmarshal([]byte(flow.ModuleVersions), &versionNumbers)
		versionNumbersList = append(versionNumbersList, versionNumbers...)
	}
	logger.SugarLogger.Info(methodName, "list of modules ", versionNumbersList)
	for _, num := range versionNumbersList {
		if completeModuleVersionNumberList[num] == false {
			completeModuleVersionNumberList[num] = true
		}
	}
	f.ModuleRepository = new(repository.ModuleRepositoryImpl)
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

	f.SectionRepository = new(repository.SectionRepositoryImpl)
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

	f.FieldRepository = new(repository.FieldRepositoryImpl)
	fieldVersions = f.FieldRepository.FetchFieldVersions(enum.Active, f.MapUtil.GetKeyListFromKeyValueMap(completeFieldVersionNumberList))

	for _, fv := range fieldVersions {
		fieldVersionsMap[fv.Id] = fv
	}

	response = f.ConstructResponse(flows, moduleVersionsMap, sectionVersionsMap, fieldVersionsMap,
		completeModuleVersionNumberList, completeSectionVersionNumberList, completeFieldVersionNumberList)
	logger.SugarLogger.Info(methodName, "Returning the response ", response)
	return response
}

func (f FlowServiceUtil) FetchFlowByIdFromDB(flowExternalId string) model.Flow {
	methodName := "FetchFlowByIdFromDB:"
	logger.SugarLogger.Info(methodName, " Fetching flows from db for flow id ", flowExternalId)
	var flow model.Flow
	f.FlowRepository = new(repository.FlowRepositoryImpl)
	flow = f.FlowRepository.FindByExternalId(flowExternalId)
	return flow
}

func (f FlowServiceUtil) ConstructFlowResponseWithModuleFieldSection(flow model.Flow,
	completeModuleVersionNumberList map[int]bool, moduleVersionsMap map[int]model.ModuleVersion,
	completeSectionVersionNumberList map[int]bool, sectionVersionsMap map[int]model.SectionVersion,
	completeFieldVersionNumberList map[int]bool, fieldVersionsMap map[int]model.FieldVersion) response_dto.FlowResponseDto {
	methodName := "ConstructFlowResponseWithModuleFieldSection"
	logger.SugarLogger.Info(methodName, "fetching the response for flow data")

	flowResponseDto := response_dto.FlowResponseDto{
		Name:       flow.Name,
		ExternalId: flow.ExternalId,
		Version:    flow.Version,
		Type:       flow.Type.String()}
	var moduleVersionNumberList []int
	json.Unmarshal([]byte(flow.ModuleVersions), &moduleVersionNumberList)
	for _, mvn := range moduleVersionNumberList {
		if completeModuleVersionNumberList[mvn] == true {
			moduleVersion := moduleVersionsMap[mvn]
			if (model.ModuleVersion{}) == moduleVersion {
				continue
			}
			flowResponseDto.Modules = append(flowResponseDto.Modules, f.getModuleVersionResponseDto(moduleVersion,
				completeSectionVersionNumberList, sectionVersionsMap, completeFieldVersionNumberList, fieldVersionsMap))
		}
	}

	logger.SugarLogger.Info(methodName, "Returning the response for flow data => ", flowResponseDto)
	return flowResponseDto
}

func (f FlowServiceUtil) FetchModuleData(flow model.Flow) ([]model.ModuleVersion, map[int]bool) {
	methodName := "FetchModuleData"
	var moduleVersionList []int
	completeModuleVersionNumberList := make(map[int]bool)
	var moduleVersions []model.ModuleVersion
	json.Unmarshal([]byte(flow.ModuleVersions), &moduleVersionList)

	logger.SugarLogger.Info(methodName, "list of modules ", moduleVersionList)
	for _, num := range moduleVersionList {
		if completeModuleVersionNumberList[num] == false {
			completeModuleVersionNumberList[num] = true
		}
	}
	moduleVersions = f.ModuleRepository.FetchModuleFromModuleVersion(completeModuleVersionNumberList)
	return moduleVersions, completeModuleVersionNumberList
}

func (f FlowServiceUtil) FetchSectionsData(moduleVersions []model.ModuleVersion) ([]model.SectionVersion, map[int]model.ModuleVersion, map[int]bool) {
	methodName := "FetchSectionsData"
	var sectionNumberList []int
	moduleVersionsMap := make(map[int]model.ModuleVersion)
	completeSectionVersionNumberList := make(map[int]bool)
	var sectionVersions []model.SectionVersion
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

	sectionVersions = f.SectionRepository.FetchSectionFromSectionVersions(completeSectionVersionNumberList)
	return sectionVersions, moduleVersionsMap, completeSectionVersionNumberList
}

func (f FlowServiceUtil) FetchFieldData(sectionVersions []model.SectionVersion) ([]model.FieldVersion, map[int]model.SectionVersion, map[int]bool, map[int]model.FieldVersion) {
	methodName := "FetchFieldData"
	sectionVersionsMap := make(map[int]model.SectionVersion)
	var fieldNumbersList []int
	completeFieldVersionNumberList := make(map[int]bool)
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
	fieldVersionsMap := make(map[int]model.FieldVersion)

	var fieldVersions []model.FieldVersion
	f.FieldRepository = new(repository.FieldRepositoryImpl)
	fieldVersions = f.FieldRepository.FetchFieldFromFieldVersion(completeFieldVersionNumberList)

	for _, fv := range fieldVersions {
		fieldVersionsMap[fv.Id] = fv
	}
	return fieldVersions, sectionVersionsMap, completeFieldVersionNumberList, fieldVersionsMap
}

func (f FlowServiceUtil) ConstructResponse(flows []model.Flow,
	moduleVersionsMap map[int]model.ModuleVersion,
	sectionVersionsMap map[int]model.SectionVersion,
	fieldVersionsMap map[int]model.FieldVersion,
	completeModuleVersionNumberList map[int]bool,
	completeSectionVersionNumberList map[int]bool,
	completeFieldVersionNumberList map[int]bool) response_dto.FlowResponsesDto {
	var response response_dto.FlowResponsesDto
	for _, flow := range flows {
		response.FlowResponses = append(response.FlowResponses, f.ConstructFlowResponseWithModuleFieldSection(flow,
			completeModuleVersionNumberList, moduleVersionsMap,
			completeSectionVersionNumberList, sectionVersionsMap,
			completeFieldVersionNumberList, fieldVersionsMap))
	}
	return response
}

func (f FlowServiceUtil) GetFieldVersionResponseDto(fieldVersion model.FieldVersion) response_dto.FieldVersionsResponseDto {
	fieldVersionResponseDto := response_dto.FieldVersionsResponseDto{
		Name:       fieldVersion.Name,
		ExternalId: fieldVersion.ExternalId,
		IsVisible:  fieldVersion.IsVisible,
		Version:    fieldVersion.Version}
	json.Unmarshal([]byte(fieldVersion.Properties), &fieldVersionResponseDto.Properties)
	return fieldVersionResponseDto
}

func (f FlowServiceUtil) GetSectionVersionResponseDto(sectionVersion model.SectionVersion, completeFieldVersionNumberList map[int]bool, fieldVersionsMap map[int]model.FieldVersion) response_dto.SectionVersionsResponseDto {
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
			sectionVersionResponseDto.Fields = append(sectionVersionResponseDto.Fields, f.GetFieldVersionResponseDto(fieldVersion))
		}
	}
	return sectionVersionResponseDto
}

func (f FlowServiceUtil) getModuleVersionResponseDto(moduleVersion model.ModuleVersion, completeSectionVersionNumberList map[int]bool,
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
				f.GetSectionVersionResponseDto(sectionVersion, completeFieldVersionNumberList, fieldVersionsMap))
		}
	}
	return moduleVersionResponseDto
}
