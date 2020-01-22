package service

import (
	"encoding/json"
	"flow/cache"
	"flow/logger"
	"flow/model"
	"flow/model/response_dto"
)

type FlowService struct {
	FlowServiceUtil FlowServiceUtil
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
	redisKey := model.RedisKey{MerchantId: merchantId,
		TenantId:  tenantId,
		ChannelId: channelId}
	redisClient.Expire(redisKey.ToString(), 0)
	cachedFlow, err := redisClient.Get(redisKey.ToString()).Result()
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
			logger.SugarLogger.Info(methodName, "Failed to fetch parsed flows associated with merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId, " with error: ", err)
		} else {
			//Add expiry time for cache entry.
			response, err := json.Marshal(flowsResponse)
			if err == nil {
				logger.SugarLogger.Info(methodName, " Updating the redis client with the response")
				setStatus := redisClient.Set(redisKey.ToString(), response, 0)
				logger.SugarLogger.Info(methodName, " set redis status: ", setStatus.Val())
			}
		}
		return flowsResponse
	}
	logger.SugarLogger.Info(methodName, " UnMarshlling the cached flow response")
	json.Unmarshal([]byte(cachedFlow), &flowsResponse)
	return flowsResponse
}
