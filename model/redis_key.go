package model

import "fmt"

type RedisKey struct {
	MerchantId string
	TenantId   string
	ChannelId  string
}

func (r RedisKey) ToString() string {
	return fmt.Sprintf("[%s, %s, %s]", r.MerchantId, r.TenantId, r.ChannelId)
}
