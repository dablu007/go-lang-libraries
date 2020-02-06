package service

import (
	"flow/config"
	"flow/logger"

	// "flow/model/response_dto"
	// mock_service "flow/test/mocks/service"
	// "flow/utility"

	"testing"
)

func init() {
	service := "flow"
	environment := "dev"
	config.Init(service, environment)
	logger.InitLogger()
}

func TestJourneyService_GetJourneys(t *testing.T) {
	// 	type fields struct {
	// 		JourneyServiceUtil mock_service.NewMockJourneyServiceUtility
	// 		RequestValidator   *utility.RequestValidator
	// 	}
	// 	type args struct {
	// 		merchantId string
	// 		tenantId   string
	// 		channelId  string
	// 	}
	// 	tests := []struct {
	// 		name   string
	// 		fields fields
	// 		args   args
	// 		want   response_dto.JourneyResponsesDto
	// 	}{
	// 		// TODO: Add test cases.
	// 	}
	// 	for _, tt := range tests {
	// 		t.Run(tt.name, func(t *testing.T) {
	// 			u := JourneyService{
	// 				JourneyServiceUtil: tt.fields.JourneyServiceUtil,
	// 				RequestValidator:   tt.fields.RequestValidator,
	// 			}
	// 			if got := u.GetJourneys(tt.args.merchantId, tt.args.tenantId, tt.args.channelId); !reflect.DeepEqual(got, tt.want) {
	// 				t.Errorf("JourneyService.GetJourneys() = %v, want %v", got, tt.want)
	// 			}
	// 		})
	// 	}
	// }

	// func TestJourneyService_GetJourneyById(t *testing.T) {
	// 	type fields struct {
	// 		JourneyServiceUtil mock_service.NewMockJourneyServiceUtility
	// 		RequestValidator   *utility.RequestValidator
	// 	}
	// 	type args struct {
	// 		journeyExternalId string
	// 	}
	// 	tests := []struct {
	// 		name   string
	// 		fields fields
	// 		args   args
	// 		want   response_dto.JourneyResponseDto
	// 	}{
	// 		// TODO: Add test cases.
	// 	}
	// 	for _, tt := range tests {
	// 		t.Run(tt.name, func(t *testing.T) {
	// 			f := JourneyService{
	// 				JourneyServiceUtil: tt.fields.JourneyServiceUtil,
	// 				RequestValidator:   tt.fields.RequestValidator,
	// 			}
	// 			if got := f.GetJourneyById(tt.args.journeyExternalId); !reflect.DeepEqual(got, tt.want) {
	// 				t.Errorf("JourneyService.GetJourneyById() = %v, want %v", got, tt.want)
	// 			}
	// 		})
	// 	}
}
