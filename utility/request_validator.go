package utility

import (
	"flow/logger"
	"regexp"
)

type RequestValidator struct {
}

func (u RequestValidator)IsValidRequest(merchantId string, tenantId string, channelId string) bool{
	methodName := "IsValidRequest:"
	logger.SugarLogger.Info(methodName, " Validating the request for merchant ", merchantId, " tenantId ", tenantId , " and channel ", channelId)
	if !IsValidUUID(merchantId){
		return false
	}
	if !(tenantId == "" || IsValidUUID(tenantId)){
		return false
	}
	if !(channelId == "" || IsValidUUID(channelId)){
		return false
	}
	return true
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
