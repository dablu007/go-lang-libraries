package utility

import (
	"flow/config"
	"flow/logger"
	. "flow/utility"
	"testing"

	"github.com/google/uuid"
)

func init() {
	config.Init("flow", "dev")
}
func TestRequestValidator_IsValidRequest(t *testing.T) {
	type args struct {
		merchantId string
		tenantId   string
		channelId  string
	}
	tests := []struct {
		name string
		u    RequestValidator
		args args
		want bool
	}{
		{name: "TestWithInvalidMerchantId", args: args{merchantId: "abc123"}, want: false},
		{name: "TastWithValidMerchantIdAndEmptyTenantAndChannelId", args: args{merchantId: uuid.New().String()}, want: true},
		{name: "TestWithValidMerchantIdInvalidTenantId", args: args{merchantId: uuid.New().String(), tenantId: "abc123"}, want: false},
		{name: "TestWithValidMerchantIdAndInvalidChannelId", args: args{merchantId: uuid.New().String(), channelId: "abc123"}, want: false},
		{name: "TestWithValidMerchantIdAndInvalidTenantIdAndChannelId", args: args{merchantId: uuid.New().String(), tenantId: "xyz123", channelId: "abc123"}, want: false},
		{name: "TestWithValidMerchantIdAndValidNonEmptyTenantAndChannelId", args: args{merchantId: uuid.New().String(), tenantId: uuid.New().String(), channelId: uuid.New().String()}, want: true},
	}
	logger.InitLogger()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := RequestValidator{}
			if got := u.IsValidRequest(tt.args.merchantId, tt.args.tenantId, tt.args.channelId); got != tt.want {
				t.Errorf("RequestValidator.IsValidRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidUUID(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "TestNonUUIDString", args: args{uuid: "abc123"}, want: false},
		{name: "TestEmptyString", args: args{}, want: false},
		{name: "TestUUIDString", args: args{uuid: uuid.New().String()}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidUUID(tt.args.uuid); got != tt.want {
				t.Errorf("IsValidUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestValidator_GenerateRedisKey(t *testing.T) {
	type args struct {
		merchantId string
		tenantId   string
		channelId  string
	}
	merchantId := uuid.New().String()
	tenantId := uuid.New().String()
	channelId := uuid.New().String()
	tests := []struct {
		name string
		u    RequestValidator
		args args
		want string
	}{
		{name: "TastWithNonEmptyMerchantIdAndTenantAndEmptyChannelId", args: args{merchantId: merchantId, tenantId: tenantId}, want: merchantId + ":" + tenantId + ":"},
		{name: "TastWithNonEmptyMerchantIdAndEmptyTenantAndChannelId", args: args{merchantId: merchantId, tenantId: tenantId, channelId: channelId}, want: merchantId + ":" + tenantId + ":" + channelId},
		{name: "TastWithNonEmptyMerchantIdAndEmptyTenantAndChannelId", args: args{merchantId: merchantId}, want: merchantId + ":"},
		{name: "TastWithNonEmptyMerchantIdAndChannelIdAndEmptyTenantId", args: args{merchantId: merchantId, channelId: channelId}, want: merchantId + ":" + channelId},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := RequestValidator{}
			if got := u.GenerateRedisKey(tt.args.merchantId, tt.args.tenantId, tt.args.channelId); got != tt.want {
				t.Errorf("RequestValidator.GenerateRedisKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
