package service

import (
	"encoding/json"
	"flow/cache"
	"flow/logger"
	"flow/model"
	"flow/model/response_dto"
	"flow/utility"
	"fmt"
)

type FlowService struct {
	FlowServiceUtil FlowServiceUtil
	RequestValidator utility.RequestValidator
}

func (u FlowService) GetFlows(merchantId string, tenantId string, channelId string) response_dto.FlowResponsesDto {
	methodName := "GetFlows"
	logger.SugarLogger.Info(methodName, "Recieved request to get all the flows associated with merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId)
	redisClient := cache.GetRedisClient()
	var flowsResponse response_dto.FlowResponsesDto
	if redisClient == nil {
		logger.SugarLogger.Info(methodName, "Failed to connect with redis client. ")
		return flowsResponse
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
		flows := u.FlowServiceUtil.FetchAllFlowsFromDB(flowContext)
		flowsResponse, err := u.FlowServiceUtil.GetParsedFlowsResponse(flows)
		if err != nil {
			logger.SugarLogger.Error(methodName, "Failed to fetch parsed flows associated with merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId, " with error: ", err)
			return flowsResponse
		}
		//Do not set redis key when there is no entry for given flowContext.
		if len(flowsResponse.FlowResponses) == 0 {
			return flowsResponse
		}

		response, err := json.Marshal(flowsResponse)
		if err != nil{
			logger.SugarLogger.Error(methodName, " couldn't update redis as failed to marshal response with err: ", err)
			return flowsResponse
		}

		logger.SugarLogger.Info(methodName, " Adding redis key: ",redisKey)
		setStatus := redisClient.Set(redisKey, response, 0)
		logger.SugarLogger.Info(methodName, " Set redis key status: ", setStatus.Val(), " for key: ",redisKey)
		return flowsResponse
	}
	logger.SugarLogger.Info(methodName, " UnMarshlling the cached flow response")
	json.Unmarshal([]byte(cachedFlow), &flowsResponse)
	return flowsResponse
}

func (f FlowService) GetFlowById(flowExternalId string) response_dto.FlowResponseDto {
	methodName := "GetFlowById"
	logger.SugarLogger.Info(methodName, "Recieved request to get flow id ", flowExternalId)
	redisClient := cache.GetRedisClient()
	var flowsResponse response_dto.FlowResponseDto
	if redisClient == nil {
		logger.SugarLogger.Info(methodName, "Failed to connect with redis client. ")
		return flowsResponse
	}
	logger.SugarLogger.Info("Fetching the flow data from redis for flowExternalId ", flowExternalId)
	cachedFlow, err := redisClient.Get(flowExternalId).Result()
	if err != nil {
		logger.SugarLogger.Info(methodName, "Failed to fetch flow from redis cache for flowExternalId: ", flowExternalId, " with error: ", err)
	}
	if len(cachedFlow) == 0 {
		flow := f.FlowServiceUtil.FetchFlowByIdFromDB(flowExternalId)
		if len(flow.Name) <= 0 {
			logger.SugarLogger.Error(methodName, " Invalid flow id passed : ", flowExternalId)
			return flowsResponse
		}
		moduleVersions, completeModuleVersionNumberList := f.FlowServiceUtil.FetchModuleData(flow)
		sectionVersions, moduleVersionsMap, completeSectionVersionNumberList := f.FlowServiceUtil.FetchSectionsData(moduleVersions)
		fieldVersions,sectionVersionsMap, completeFieldVersionNumberList, fieldVersionsMap := f.FlowServiceUtil.FetchFieldData(sectionVersions)
		fmt.Println("==== ", fieldVersions)
		flowsResponse, err := f.FlowServiceUtil.GetFlowModuleSectionAndFieldData(flow, completeModuleVersionNumberList,
			moduleVersionsMap, completeSectionVersionNumberList,sectionVersionsMap,completeFieldVersionNumberList,fieldVersionsMap)
		if err != nil{
			logger.SugarLogger.Error(methodName, " couldn't update redis as failed to marshal response with err: ", err)
			return flowsResponse
		}
		response, err := json.Marshal(flowsResponse)
		if err != nil{
			logger.SugarLogger.Error(methodName, " failed to marshal response with err: will not be able to update redis", err)
			return flowsResponse
		}
		logger.SugarLogger.Info(methodName, " Adding redis key: ", flowExternalId)
		setStatus := redisClient.Set(flowExternalId, response, 0)
		logger.SugarLogger.Info(methodName, " Set redis key status: ", setStatus.Val(), " for key: ", flowExternalId)
		return flowsResponse
	}
	logger.SugarLogger.Info(methodName, " UnMarshlling the cached flow response")
	json.Unmarshal([]byte(cachedFlow), &flowsResponse)
	return flowsResponse
}
