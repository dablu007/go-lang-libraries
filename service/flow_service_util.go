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
	mapUtil utility.MapUtil
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
	for _, flow := range flows {
		var versionNumbers []int
		json.Unmarshal([]byte(flow.ModuleVersions), &versionNumbers)
		for _, num := range versionNumbers {
			if completeModuleVersionNumberList[num] == false {
				completeModuleVersionNumberList[num] = true
			}
		}
	}

	dbConnection.Joins("JOIN modules ON modules.id = module_versions.module_id and modules.status = ? and modules.deleted_on is NULL", enum.Active).Where("module_versions.id in (?) and module_versions.deleted_on is NULL", f.mapUtil.GetKeyListFromKeyValueMap(completeModuleVersionNumberList)).Find(&moduleVersions)

	for _, mv := range moduleVersions {
		moduleVersionsMap[mv.Id] = mv
		var sectionNumbers []int
		json.Unmarshal([]byte(mv.SectionVersions), &sectionNumbers)
		for _, num := range sectionNumbers {
			if completeSectionVersionNumberList[num] == false {
				completeSectionVersionNumberList[num] = true
			}
		}
	}

	dbConnection.Joins("JOIN sections ON sections.id = section_versions.section_id and sections.status = ? and sections.deleted_on is NULL", enum.Active).Where("section_versions.id in (?) and section_versions.deleted_on is NULL", f.mapUtil.GetKeyListFromKeyValueMap(completeSectionVersionNumberList)).Find(&sectionVersions)

	for _, sv := range sectionVersions {
		sectionVersionsMap[sv.Id] = sv
		var fieldNumbers []int
		json.Unmarshal([]byte(sv.FieldVersions), &fieldNumbers)
		for _, num := range fieldNumbers {
			if completeFieldVersionNumberList[num] == false {
				completeFieldVersionNumberList[num] = true
			}
		}
	}

	dbConnection.Joins("JOIN fields ON fields.id = field_versions.field_id and fields.status = ? and fields.deleted_on is NULL", enum.Active).Where("field_versions.id in (?) and field_versions.deleted_on is NULL", f.mapUtil.GetKeyListFromKeyValueMap(completeFieldVersionNumberList)).Find(&fieldVersions)

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
