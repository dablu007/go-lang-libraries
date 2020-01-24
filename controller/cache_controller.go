package controller

import (
	"flow/cache"
	"flow/logger"
	"flow/model"
	"flow/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CacheController struct {
	requestValidator utility.RequestValidator
}

/*
This method will evict the cache entry for given flow context.
*/

func (u CacheController) DeleteCacheEntry() gin.HandlerFunc {
	methodName := "DeleteCacheEntry:"
	fn := func(c *gin.Context) {
		merchantId := c.Query("merchantId")
		tenantId := c.Query("tenantId")
		channelId := c.Query("channelId")
		if u.requestValidator.IsValidRequest(merchantId, tenantId, channelId) {
			logger.SugarLogger.Info(methodName, "Deleting cache entry for merchant:", merchantId, " tenant:", tenantId, " channel: ", channelId)
			redisClient := cache.GetRedisClient()
			if redisClient == nil {
				logger.SugarLogger.Info(methodName, "Failed to connect with redis client. ")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to connect to redis."})
				return
			}
			redisKey := model.RedisKey{MerchantId: merchantId,
				TenantId:  tenantId,
				ChannelId: channelId}

			_, err := redisClient.Del(redisKey.ToString()).Result()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input parameters."})
	}
	return fn

}
