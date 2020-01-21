package service

import (
	"encoding/json"
	"flow/cache"
	"flow/db"
	"flow/enum"
	"flow/logger"
	"flow/model"
	"flow/model/response_dto"
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
	var response response_dto.FlowResponsesDto
	for _, flow := range flows {
		flowResponseDto := response_dto.FlowResponseDto{
			Name:       flow.Name,
			ExternalId: flow.ExternalId,
			Version:    flow.Version,
			Type:       flow.Type}
		var moduleVersionNumberList []int
		json.Unmarshal([]byte(flow.ModuleVersions), &moduleVersionNumberList)
		for i := 0; i < len(moduleVersionNumberList); i++ {
			var moduleVersion model.ModuleVersion
			dbConnection.Debug().Joins("JOIN modules ON modules.id = module_versions.module_id and modules.status = ? and modules.deleted_on is NULL", enum.Active).Where("module_versions.id = ? and module_versions.deleted_on is NULL", moduleVersionNumberList[i]).Find(&moduleVersion)
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
			for i := 0; i < len(sectionVersionNumberList); i++ {
				var sectionVersion model.SectionVersion
				dbConnection.Debug().Joins("JOIN sections ON sections.id = section_versions.section_id and sections.status = ? and sections.deleted_on is NULL", enum.Active).Where("section_versions.id = ? and section_versions.deleted_on is NULL", sectionVersionNumberList[i]).Find(&sectionVersion)
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
				for i := 0; i < len(fieldVersionNumberList); i++ {
					var fieldVersion model.FieldVersion
					dbConnection.Debug().Joins("JOIN fields ON fields.id = field_versions.field_id and fields.status = ? and fields.deleted_on is NULL", enum.Active).Where("field_versions.id = ? and field_versions.deleted_on is NULL", fieldVersionNumberList[i]).Find(&fieldVersion)
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
				moduleVersionResponseDto.Sections = append(moduleVersionResponseDto.Sections, sectionVersionResponseDto)
			}
			flowResponseDto.Modules = append(flowResponseDto.Modules, moduleVersionResponseDto)
		}
		response.FlowResponses = append(response.FlowResponses, flowResponseDto)
	}
	return response, nil
}
