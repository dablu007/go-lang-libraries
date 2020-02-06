package service

import (
	"encoding/json"
	"flow/cache"
	"flow/logger"
	"flow/model"
	"flow/model/response_dto"
	"flow/utility"
)

type JourneyService struct {
	JourneyServiceUtil *JourneyServiceUtil
	RequestValidator   *utility.RequestValidator
}

func NewJourneyService(journeyService *JourneyServiceUtil, validator *utility.RequestValidator) *JourneyService {
	service := &JourneyService{
		JourneyServiceUtil: journeyService,
		RequestValidator:   validator,
	}
	return service
}

func (u JourneyService) GetJourneys(merchantId string, tenantId string, channelId string) response_dto.JourneyResponsesDto {
	methodName := "GetJourneys"
	logger.SugarLogger.Info(methodName, "Recieved request to get all the flows associated with merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId)
	redisClient := cache.GetRedisClient()
	var journeyResponsesDto response_dto.JourneyResponsesDto
	if redisClient == nil {
		logger.SugarLogger.Info(methodName, "Failed to connect with redis client. ")
		return journeyResponsesDto
	}
	logger.SugarLogger.Info(methodName, "Fetching flows from redis cache for merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId)
	redisKey := u.RequestValidator.GenerateRedisKey(merchantId, tenantId, channelId)
	cachedFlow, err := redisClient.Get(redisKey).Result()
	if err != nil {
		logger.SugarLogger.Info(methodName, "Failed to fetch flow from redis cache for merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId, " with error: ", err)
	}
	if cachedFlow == "" {
		logger.SugarLogger.Info(methodName, "No flows exist in redis cache for merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId)
		flowContext := model.FlowContext{
			MerchantId: merchantId,
			TenantId:   tenantId,
			ChannelId:  channelId}
		flows := u.JourneyServiceUtil.FetchAllJourneysFromDB(flowContext)
		moduleVersionsMap, sectionVersionsMap, fieldVersionsMap, completeModuleVersionNumberList, completeSectionVersionNumberList,
			completeFieldVersionNumberList := u.JourneyServiceUtil.GetModuleSectionAndFieldVersionsAndActiveVersionNumberList(flows...)
		journeyResponse := u.JourneyServiceUtil.ConstructJourneysResponse(flows, moduleVersionsMap, sectionVersionsMap, fieldVersionsMap,
			completeModuleVersionNumberList, completeSectionVersionNumberList, completeFieldVersionNumberList)

		//Do not set redis key when there is no entry for given flowContext.
		if len(journeyResponse.JourneyResponses) == 0 {
			return journeyResponse
		}

		response, err := json.Marshal(journeyResponse)
		if err != nil {
			logger.SugarLogger.Error(methodName, " couldn't update redis as failed to marshal response with err: ", err)
			return journeyResponse
		}

		logger.SugarLogger.Info(methodName, " Adding redis key: ", redisKey)
		setStatus := redisClient.Set(redisKey, response, 0)
		logger.SugarLogger.Info(methodName, " Set redis key status: ", setStatus.Val(), " for key: ", redisKey)
		return journeyResponse
	}
	logger.SugarLogger.Info(methodName, " UnMarshlling the cached flow response")
	json.Unmarshal([]byte(cachedFlow), &journeyResponsesDto)
	return journeyResponsesDto
}

func (f JourneyService) GetJourneyById(journeyExternalId string) response_dto.JourneyResponseDto {
	methodName := "GetJourneyById"
	logger.SugarLogger.Info(methodName, "Recieved request to get flow id ", journeyExternalId)
	redisClient := cache.GetRedisClient()
	var journeyResponseDto response_dto.JourneyResponseDto
	if redisClient == nil {
		logger.SugarLogger.Info(methodName, "Failed to connect with redis client. ")
		return journeyResponseDto
	}
	logger.SugarLogger.Info("Fetching the flow data from redis for journeyExternalId ", journeyExternalId)
	//redisClient.FlushAll()
	key := "journeyId:" + journeyExternalId + ":nested"
	cachedFlow, err := redisClient.Get(key).Result()
	if err != nil {
		logger.SugarLogger.Info(methodName, "Failed to fetch flow from redis cache for journeyExternalId: ", journeyExternalId, " with error: ", err)
	}
	if len(cachedFlow) == 0 {
		flow := f.JourneyServiceUtil.FetchJourneyByIdFromDB(journeyExternalId)
		if len(flow.Name) <= 0 {
			logger.SugarLogger.Error(methodName, " Invalid flow id passed : ", journeyExternalId)
			return journeyResponseDto
		}
		moduleVersionsMap, sectionVersionsMap, fieldVersionsMap, completeModuleVersionNumberList, completeSectionVersionNumberList,
			completeFieldVersionNumberList := f.JourneyServiceUtil.GetModuleSectionAndFieldVersionsAndActiveVersionNumberList(flow)
		flowsResponse := f.JourneyServiceUtil.ConstructFlowResponseWithModuleFieldSection(flow, completeModuleVersionNumberList,
			moduleVersionsMap, completeSectionVersionNumberList, sectionVersionsMap, completeFieldVersionNumberList, fieldVersionsMap)

		response, err := json.Marshal(flowsResponse)
		if err != nil {
			logger.SugarLogger.Error(methodName, " failed to marshal response with err: will not be able to update redis", err)
			return flowsResponse
		}
		logger.SugarLogger.Info(methodName, " Adding redis key: ", journeyExternalId)
		setStatus := redisClient.Set(key, response, 0)
		logger.SugarLogger.Info(methodName, " Set redis key status: ", setStatus.Val(), " for key: ", journeyExternalId)
		return flowsResponse
	}
	logger.SugarLogger.Info(methodName, " UnMarshlling the cached flow response")
	json.Unmarshal([]byte(cachedFlow), &journeyResponseDto)
	return journeyResponseDto
}

func (f JourneyService) GetJourneyByIdNotNested(journeyExternalId string) response_dto.JourneyResponseDtoList {
	methodName := "GetJourneyById"
	logger.SugarLogger.Info(methodName, "Recieved request to get flow id ", journeyExternalId)
	redisClient := cache.GetRedisClient()
	var journeyResponseDto response_dto.JourneyResponseDtoList
	if redisClient == nil {
		logger.SugarLogger.Info(methodName, "Failed to connect with redis client. ")
		return journeyResponseDto
	}
	logger.SugarLogger.Info("Fetching the flow data from redis for journeyExternalId ", journeyExternalId)
	//redisClient.FlushAll()
	key := "journeyId:" + journeyExternalId
	cachedFlow, err := redisClient.Get(key).Result()
	if err != nil {
		logger.SugarLogger.Info(methodName, "Failed to fetch flow from redis cache for journeyExternalId: ", journeyExternalId, " with error: ", err)
	}
	if len(cachedFlow) == 0 {
		flow := f.JourneyServiceUtil.FetchJourneyByIdFromDB(journeyExternalId)
		if len(flow.Name) <= 0 {
			logger.SugarLogger.Error(methodName, " Invalid flow id passed : ", journeyExternalId)
			return journeyResponseDto
		}
		moduleVersionsMap, sectionVersionsMap, fieldVersionsMap, _, _, _ := f.JourneyServiceUtil.GetModuleSectionAndFieldVersionsAndActiveVersionNumberList(flow)
		flowsResponse := f.JourneyServiceUtil.ConstructFlowResponseNotNested(flow,moduleVersionsMap,sectionVersionsMap,fieldVersionsMap)

		response, err := json.Marshal(flowsResponse)
		if err != nil {
			logger.SugarLogger.Error(methodName, " failed to marshal response with err: will not be able to update redis", err)
			return flowsResponse
		}
		logger.SugarLogger.Info(methodName, " Adding redis key: ", journeyExternalId)
		setStatus := redisClient.Set(key, response, 0)
		logger.SugarLogger.Info(methodName, " Set redis key status: ", setStatus.Val(), " for key: ", journeyExternalId)
		return flowsResponse
	}
	logger.SugarLogger.Info(methodName, " UnMarshlling the cached flow response")
	json.Unmarshal([]byte(cachedFlow), &journeyResponseDto)
	return journeyResponseDto
}
