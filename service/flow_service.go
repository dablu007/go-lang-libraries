package service

import (
	"encoding/json"
	"flow/cache"
	"flow/db"
	"flow/enum"
	"flow/logger"
	"flow/model"
	"flow/model/response_dto"
	"fmt"
)

type FlowService struct {
	redisKey model.RedisKey
}

func (u FlowService) GetFlows(merchantId string, tenantId string, channelId string) response_dto.FlowResponsesDto {
	methodName := "GetFlows"
	logger.SugarLogger.Info(methodName, "Recieved request to get all the flows associated with merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId)
	redisClient := cache.GetRedisClient()
	var flowsResponse response_dto.FlowResponsesDto
	if redisClient != nil {
		logger.SugarLogger.Info(methodName, "Fetching flows from redis cache for merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId)
		redisKey := model.RedisKey{MerchantId: merchantId,
			TenantId:  tenantId,
			ChannelId: channelId}

		redisClient.Expire(redisKey.ToString(), 0)
		cachedFlow, err := redisClient.Get(redisKey.ToString()).Result()
		if err != nil {
			logger.SugarLogger.Info(methodName, "Failed to fetch flows from redis cache for merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId, " with error: ", err)
		}
		if cachedFlow == "" {
			logger.SugarLogger.Info(methodName, "No flows exist in redis cache for merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId)
			flowContext := model.FlowContext{
				MerchantId: merchantId,
				TenantId:   tenantId,
				ChannelId:  channelId}
			flows := FetchAllFlowsFromDB(flowContext)
			flowsResponse, err := GetParsedFlowsResponse(flows)
			if err != nil {
				logger.SugarLogger.Info(methodName, "Failed to fetch parsed flows associated with merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId, " with error: ", err)
			} else {
				//Add expiry time for cache entry.
				response, err := json.Marshal(flowsResponse)
				if err == nil {
					setStatus := redisClient.Set(redisKey.ToString(), response, 0)
					logger.SugarLogger.Info(methodName, " set redis status: ", setStatus, " key: ", redisKey.ToString)
				}
			}
			return flowsResponse

		} else {
			json.Unmarshal([]byte(cachedFlow), &flowsResponse)
		}
	} else {
		logger.SugarLogger.Info(methodName, "Failed to connect with redis client. ")
	}
	return flowsResponse
}

func FetchAllFlowsFromDB(flowContext model.FlowContext) []model.Flow {
	dbConnection := db.GetDB()
	var flows []model.Flow
	if dbConnection == nil {
		return flows
	}
	dbConnection.Debug().Where("flow_context->>'MerchantId' = ? and flow_context->>'TenantId' = ? and flow_context->>'ChannelId' = ? and status = ? and deleted_on is NULL", flowContext.MerchantId, flowContext.TenantId, flowContext.ChannelId, enum.Active).Find(&flows)
	return flows
}

func GetParsedFlowsResponse(flows []model.Flow) (response_dto.FlowResponsesDto, error) {
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

	keys := []int{}
	for k, _ := range completeModuleVersionNumberList {
		keys = append(keys, k)
	}

	fmt.Println(keys)

	dbConnection.Debug().Joins("JOIN modules ON modules.id = module_versions.module_id and modules.status = ? and modules.deleted_on is NULL", enum.Active).Where("module_versions.id in (?) and module_versions.deleted_on is NULL", keys).Find(&moduleVersions)

	for _, mv := range moduleVersions {
		if (moduleVersionsMap[mv.Id] == model.ModuleVersion{}) {
			moduleVersionsMap[mv.Id] = mv
		}
	}

	for _, mv := range moduleVersions {
		var sectionNumbers []int
		json.Unmarshal([]byte(mv.SectionVersions), &sectionNumbers)
		for _, num := range sectionNumbers {
			if completeSectionVersionNumberList[num] == false {
				completeSectionVersionNumberList[num] = true
			}
		}
	}

	keys1 := []int{}
	for k, _ := range completeSectionVersionNumberList {
		keys1 = append(keys1, k)
	}

	fmt.Println(keys1)

	dbConnection.Debug().Joins("JOIN sections ON sections.id = section_versions.section_id and sections.status = ? and sections.deleted_on is NULL", enum.Active).Where("section_versions.id in (?) and section_versions.deleted_on is NULL", keys1).Find(&sectionVersions)

	for _, sv := range sectionVersions {
		if (sectionVersionsMap[sv.Id] == model.SectionVersion{}) {
			sectionVersionsMap[sv.Id] = sv
		}
	}

	for _, sv := range sectionVersions {
		var fieldNumbers []int
		json.Unmarshal([]byte(sv.FieldVersions), &fieldNumbers)
		for _, num := range fieldNumbers {
			if completeFieldVersionNumberList[num] == false {
				completeFieldVersionNumberList[num] = true
			}
		}
	}

	keys2 := []int{}
	for k, _ := range completeFieldVersionNumberList {
		keys2 = append(keys2, k)
	}

	fmt.Println(keys2)

	dbConnection.Debug().Joins("JOIN fields ON fields.id = field_versions.field_id and fields.status = ? and fields.deleted_on is NULL", enum.Active).Where("field_versions.id in (?) and field_versions.deleted_on is NULL", keys2).Find(&fieldVersions)

	for _, fv := range fieldVersions {
		if (fieldVersionsMap[fv.Id] == model.FieldVersion{}) {
			fieldVersionsMap[fv.Id] = fv
		}
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
