package utility

import (
	"flow/logger"
	"regexp"
)

type RequestValidatorUtil interface {
	IsValidRequest(merchantId string, tenantId string, channelId string) bool
	IsValidUUID(uuid string) bool
	GenerateRedisKey(merchantId string, tenantId string, channelId string) string
}

type RequestValidator struct {
}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

func (u RequestValidator) IsValidRequest(merchantId string, tenantId string, channelId string) bool {
	methodName := "IsValidRequest:"
	logger.SugarLogger.Info(methodName, " Validating the request for merchant ", merchantId, " tenantId ", tenantId, " and channel ", channelId)
	if !IsValidUUID(merchantId) {
		return false
	}
	if !(tenantId == "" || IsValidUUID(tenantId)) {
		return false
	}
	if !(channelId == "" || IsValidUUID(channelId)) {
		return false
	}
	return true
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func (u RequestValidator) GenerateRedisKey(merchantId string, tenantId string, channelId string) string {
	var key = ""
	if len(merchantId) > 0 {
		key = key + merchantId + ":"
	}
	if len(tenantId) > 0 {
		key = key + tenantId + ":"
	}
	if len(channelId) > 0 {
		key = key + channelId
	}
	return key
}
