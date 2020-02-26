package utility

import (
	"flow/logger"
	"regexp"

	"github.com/google/uuid"
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

//IsValidUUID : returns true if the UUID is valid UUID 
func IsValidUUID(id string) bool {
	parsedId, err := uuid.Parse(id)
	return err!=nil && parsedId!=uuid.Nil
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
