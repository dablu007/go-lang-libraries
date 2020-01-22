package service

import (
	"encoding/json"
	"flow/db"
	"flow/enum"
	"flow/model"
	"flow/model/response_dto"
	"flow/utility"
)

type FlowServiceUtil struct {
	MapUtil utility.MapUtil
}

func (f FlowServiceUtil) FetchAllFlowsFromDB(flowContext model.FlowContext) []model.Flow {
	dbConnection := db.GetDB()
	var flows []model.Flow
	if dbConnection == nil {
		return flows
	}
	dbConnection.Where("flow_context->>'MerchantId' = ? and flow_context->>'TenantId' = ? and flow_context->>'ChannelId' = ? and status = ? and deleted_on is NULL", flowContext.MerchantId, flowContext.TenantId, flowContext.ChannelId, enum.Active).Find(&flows)
	return flows
}

func (f FlowServiceUtil) GetParsedFlowsResponse(flows []model.Flow) (response_dto.FlowResponsesDto, error) {
	dbConnection := db.GetDB()
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

	for _, num := range versionNumbersList {
		if completeModuleVersionNumberList[num] == false {
			completeModuleVersionNumberList[num] = true
		}
	}
	logrs.Println("fetching from db")
	dbConnection.Joins("JOIN module ON module.id = module_version.module_id and module.status = ? and module.deleted_on is NULL", enum.Active).Where("module_version.id in (?) and module_version.deleted_on is NULL", f.MapUtil.GetKeyListFromKeyValueMap(completeModuleVersionNumberList)).Find(&moduleVersions)

	var sectionNumberList []int
	for _, mv := range moduleVersions {
		moduleVersionsMap[mv.Id] = mv
		var sectionNumbers []int
		json.Unmarshal([]byte(mv.SectionVersions), &sectionNumbers)
		sectionNumberList = append(sectionNumberList, sectionNumbers...)
	}
	for _, num := range sectionNumberList {
		if completeSectionVersionNumberList[num] == false {
			completeSectionVersionNumberList[num] = true
		}
	}

	dbConnection.Joins("JOIN section ON section.id = section_version.section_id and section.status = ? and section.deleted_on is NULL", enum.Active).Where("section_version.id in (?) and section_version.deleted_on is NULL", f.MapUtil.GetKeyListFromKeyValueMap(completeSectionVersionNumberList)).Find(&sectionVersions)

	var fieldNumbersList []int
	for _, sv := range sectionVersions {
		sectionVersionsMap[sv.Id] = sv
		var fieldNumbers []int
		json.Unmarshal([]byte(sv.FieldVersions), &fieldNumbers)
		fieldNumbersList = append(fieldNumbersList, fieldNumbers...)
	}

	for _, num := range fieldNumbersList {
		if completeFieldVersionNumberList[num] == false {
			completeFieldVersionNumberList[num] = true
		}
	}

	dbConnection.Joins("JOIN field ON field.id = field_version.field_id and field.status = ? and field.deleted_on is NULL", enum.Active).Where("field_version.id in (?) and field_version.deleted_on is NULL", f.MapUtil.GetKeyListFromKeyValueMap(completeFieldVersionNumberList)).Find(&fieldVersions)

	for _, fv := range fieldVersions {
		fieldVersionsMap[fv.Id] = fv
	}

	for _, flow := range flows {
		flowResponseDto := response_dto.FlowResponseDto{
			Name:       flow.Name,
			ExternalId: flow.ExternalId,
			Version:    flow.Version,
			Type:       flow.Type}
		var moduleVersionNumberList []int
		json.Unmarshal([]byte(flow.ModuleVersions), &moduleVersionNumberList)
		for _, mvn := range moduleVersionNumberList {
			if completeModuleVersionNumberList[mvn] == true {
				moduleVersion := moduleVersionsMap[mvn]
				if (model.ModuleVersion{}) == moduleVersion {
					continue
				}
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
								fieldVersionResponseDto := response_dto.FieldVersionsResponseDto{
									Name:       fieldVersion.Name,
									ExternalId: fieldVersion.ExternalId,
									IsVisible:  fieldVersion.IsVisible,
									Version:    fieldVersion.Version}
								json.Unmarshal([]byte(fieldVersion.Properties), &fieldVersionResponseDto.Properties)
								sectionVersionResponseDto.Fields = append(sectionVersionResponseDto.Fields, fieldVersionResponseDto)
							}
						}
						moduleVersionResponseDto.Sections = append(moduleVersionResponseDto.Sections, sectionVersionResponseDto)
					}
				}
				flowResponseDto.Modules = append(flowResponseDto.Modules, moduleVersionResponseDto)
			}
		}
		response.FlowResponses = append(response.FlowResponses, flowResponseDto)
	}
	return response, nil
}
