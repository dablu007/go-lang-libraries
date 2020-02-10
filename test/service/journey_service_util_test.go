package service

import (
	"flow/config"
	"flow/enum"
	"flow/enum/flow_status"
	"flow/logger"
	"flow/model"
	"flow/model/response_dto"
	"flow/service"
	mock_repository "flow/test/mocks/repository"
	mock_utility "flow/test/mocks/utility"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/magiconair/properties/assert"
)

func init() {
	service := "flow"
	environment := "dev"
	config.Init(service, environment)
	logger.InitLogger()
}

func TestFetchAllJourneysFromDB(t *testing.T) {
	var controller = gomock.NewController(t)
	defer controller.Finish()

	journeys := []model.Journey{}

	var merchantId = "m123"

	//test for empty jorney list.
	var journeyRepository = mock_repository.NewMockJourneyRepository(controller)
	journeyRepository.EXPECT().FindActiveJourneysByJourneyContext(merchantId, "", "").Return(journeys)
	var journeyService = &service.JourneyServiceUtil{
		JourneyRepository: journeyRepository,
	}
	var journeyActual = journeyService.FetchAllJourneysFromDB(model.FlowContext{MerchantId: merchantId})
	assert.Equal(t, journeyActual, journeys)

	//test for non empty journey list.
	var externalId = uuid.New()
	journeys = []model.Journey{{Id: 1, CreatedOn: time.Now(), ExternalId: externalId, Name: "CLflow", Type: enum.CL, Status: enum.Active, ModuleVersions: "[1,2,3]", Version: "1.1", FlowContext: ""}}
	journeyRepository.EXPECT().FindActiveJourneysByJourneyContext(merchantId, "", "").Return(journeys)

	journeyActual = journeyService.FetchAllJourneysFromDB(model.FlowContext{MerchantId: merchantId})
	assert.Equal(t, journeyActual, journeys)

}

func TestFetchjourneyByIdFromDB(t *testing.T) {
	var controller = gomock.NewController(t)
	defer controller.Finish()

	var journey model.Journey
	journey.Id = 1
	journey.CreatedOn = time.Now()

	var journeyRepository = mock_repository.NewMockJourneyRepository(controller)
	var externalId = "123"
	journeyRepository.EXPECT().FindByExternalId(externalId).Return(journey)

	var journeyService = &service.JourneyServiceUtil{
		JourneyRepository: journeyRepository,
	}
	var journeyActual = journeyService.FetchJourneyByIdFromDB(externalId)
	assert.Equal(t, journeyActual.Id, journey.Id)

	journeyRepository.EXPECT().FindByExternalId("").Return(model.Journey{})

	journeyActual = journeyService.FetchJourneyByIdFromDB("")
	assert.Equal(t, len(journeyActual.Name), 0)
}

func TestGetModuleSectionAndFieldVersionsAndActiveVersionNumberList(t *testing.T) {
	//Todo: fix this one
	var controller = gomock.NewController(t)
	defer controller.Finish()

	var journey model.Journey = model.Journey{}
	var moduleRepo = mock_repository.NewMockModuleRepository(controller)
	var sectionRepo = mock_repository.NewMockSectionRepository(controller)
	var fieldRepo = mock_repository.NewMockFieldRepository(controller)
	var mapUtil = mock_utility.NewMockMapUtility(controller)

	var journeyServiceUtil = &service.JourneyServiceUtil{
		FieldRepository:   fieldRepo,
		ModuleRepository:  moduleRepo,
		SectionRepository: sectionRepo,
		MapUtil:           mapUtil,
	}

	var statusActive = enum.Status.String(enum.Active)

	//test with no journey
	mapUtil.EXPECT().GetKeyListFromKeyValueMap(map[int]bool{}).Return([]int{})
	moduleRepo.EXPECT().FetchModuleVersions(statusActive, []int{}).Return([]model.ModuleVersion{})
	sectionRepo.EXPECT().FetchSectionVersions(statusActive, []int{}).Return([]model.SectionVersion{})
	fieldRepo.EXPECT().FetchFieldVersions(statusActive, []int{}).Return([]model.FieldVersion{})

	var actualModuleVersionsMap map[int]model.ModuleVersion
	var actualSectionVersionsMap map[int]model.SectionVersion
	var actualFieldVersionsMap map[int]model.FieldVersion
	actualModuleVersionsMap, actualSectionVersionsMap, actualFieldVersionsMap =
		journeyServiceUtil.GetModuleSectionAndFieldVersionsAndActiveVersionNumberList(journey)
	fmt.Println("length", len(actualModuleVersionsMap))
	assert.Equal(t, len(actualModuleVersionsMap), 0)
	assert.Equal(t, len(actualSectionVersionsMap), 0)
	assert.Equal(t, len(actualFieldVersionsMap), 0)
}

func TestConstructJourneysResponse(t *testing.T) {
	var controller = gomock.NewController(t)
	defer controller.Finish()

	var journeyServiceUtil = &service.JourneyServiceUtil{}

	var journeys []model.Journey
	var moduleVersionsMap map[int]model.ModuleVersion
	var sectionVersionsMap map[int]model.SectionVersion
	var fieldVersionsMap map[int]model.FieldVersion

	//test for empty response.
	var actualResponse = journeyServiceUtil.ConstructJourneysResponse(journeys, moduleVersionsMap,
		sectionVersionsMap, fieldVersionsMap)

	assert.Equal(t, actualResponse, response_dto.JourneyResponsesDto{})

	//test for non empty response.

	moduleVersionsMap = make(map[int]model.ModuleVersion)
	sectionVersionsMap = make(map[int]model.SectionVersion)
	fieldVersionsMap = make(map[int]model.FieldVersion)

	journeys = append(journeys, model.Journey{Id: 1, ModuleVersions: "[1,2]", Status: flow_status.Active, Type: enum.CL})
	moduleVersionsMap[1] = model.ModuleVersion{ModuleId: 1, Name: "Mod1", SectionVersions: "[1,2,3]"}
	moduleVersionsMap[2] = model.ModuleVersion{ModuleId: 2, Name: "Mod2", SectionVersions: "[1]"}

	sectionVersionsMap[1] = model.SectionVersion{SectionId: 1, Name: "Sec1", FieldVersions: "[1]"}
	sectionVersionsMap[2] = model.SectionVersion{SectionId: 2, Name: "Sec2", FieldVersions: "[1,2]"}

	fieldVersionsMap[1] = model.FieldVersion{FieldId: 1, Name: "Field1"}
	fieldVersionsMap[2] = model.FieldVersion{FieldId: 2, Name: "Field2"}

	var expectedResponse = response_dto.JourneyResponsesDto{
		JourneyResponses: []response_dto.JourneyResponseDto{
			{Type: "CL", Modules: []response_dto.ModuleVersionResponseDto{
				{Name: "Mod1", Sections: []response_dto.SectionVersionsResponseDto{
					{Name: "Sec1", Fields: []response_dto.FieldVersionsResponseDto{
						{Name: "Field1"},
					}},
					{Name: "Sec2", Fields: []response_dto.FieldVersionsResponseDto{
						{Name: "Field1"},
						{Name: "Field2"},
					}},
				}},
				{Name: "Mod2", Sections: []response_dto.SectionVersionsResponseDto{
					{Name: "Sec1", Fields: []response_dto.FieldVersionsResponseDto{
						{Name: "Field1"},
					}},
				}},
			}},
		},
	}

	actualResponse = journeyServiceUtil.ConstructJourneysResponse(journeys, moduleVersionsMap,
		sectionVersionsMap, fieldVersionsMap)

	assert.Equal(t, actualResponse, expectedResponse)
}

func TestConstructFlowResponseWithModuleFieldSection(t *testing.T) {
	var controller = gomock.NewController(t)
	defer controller.Finish()

	var journeyServiceUtil = &service.JourneyServiceUtil{}
	var journey = model.Journey{}
	var moduleVersionsMap map[int]model.ModuleVersion
	var sectionVersionsMap map[int]model.SectionVersion
	var fieldVersionsMap map[int]model.FieldVersion

	//test for empty response.
	var actualResponse = journeyServiceUtil.ConstructFlowResponseWithModuleFieldSection(journey,
		moduleVersionsMap, sectionVersionsMap, fieldVersionsMap)

	assert.Equal(t, len(actualResponse.Modules), 0)
	assert.Equal(t, len(actualResponse.Name), 0)

	//test for non empty response.

	moduleVersionsMap = make(map[int]model.ModuleVersion)
	sectionVersionsMap = make(map[int]model.SectionVersion)
	fieldVersionsMap = make(map[int]model.FieldVersion)

	journey = model.Journey{Id: 1, ModuleVersions: "[1,2]", Status: flow_status.Active, Type: enum.CL}

	moduleVersionsMap[1] = model.ModuleVersion{ModuleId: 1, Name: "Mod1", SectionVersions: "[1,2,3]"}
	moduleVersionsMap[2] = model.ModuleVersion{ModuleId: 2, Name: "Mod2", SectionVersions: "[1]"}

	sectionVersionsMap[1] = model.SectionVersion{SectionId: 1, Name: "Sec1", FieldVersions: "[1]"}
	sectionVersionsMap[2] = model.SectionVersion{SectionId: 2, Name: "Sec2", FieldVersions: "[1,2]"}

	fieldVersionsMap[1] = model.FieldVersion{FieldId: 1, Name: "Field1"}
	fieldVersionsMap[2] = model.FieldVersion{FieldId: 2, Name: "Field2"}

	var expectedResponse = response_dto.JourneyResponseDto{
		Type: "CL", Modules: []response_dto.ModuleVersionResponseDto{
			{Name: "Mod1", Sections: []response_dto.SectionVersionsResponseDto{
				{Name: "Sec1", Fields: []response_dto.FieldVersionsResponseDto{
					{Name: "Field1"},
				}},
				{Name: "Sec2", Fields: []response_dto.FieldVersionsResponseDto{
					{Name: "Field1"},
					{Name: "Field2"},
				}},
			}},
			{Name: "Mod2", Sections: []response_dto.SectionVersionsResponseDto{
				{Name: "Sec1", Fields: []response_dto.FieldVersionsResponseDto{
					{Name: "Field1"},
				}},
			}},
		},
	}

	actualResponse = journeyServiceUtil.ConstructFlowResponseWithModuleFieldSection(journey,
		moduleVersionsMap, sectionVersionsMap, fieldVersionsMap)

	assert.Equal(t, actualResponse, expectedResponse)
}

func TestConstructFlowResponseAsList(t *testing.T) {
	var controller = gomock.NewController(t)
	defer controller.Finish()

	var journeyServiceUtil = &service.JourneyServiceUtil{}
	var journey = model.Journey{}
	var moduleVersions = map[int]model.ModuleVersion{}
	var sectionVersions = map[int]model.SectionVersion{}
	var fieldVersions = map[int]model.FieldVersion{}

	//Test for empty response.
	var actualResponse = journeyServiceUtil.ConstructFlowResponseAsList(journey, moduleVersions,
		sectionVersions, fieldVersions)

	assert.Equal(t, actualResponse, response_dto.JourneyResponseDtoList{})

	//Test for non-empty respopnse.
	journey = model.Journey{Id: 1, Name: "abc", Type: enum.CL, ModuleVersions: "[1,2]"}
	moduleVersions = make(map[int]model.ModuleVersion)
	moduleVersions[1] = model.ModuleVersion{Id: 1, Name: "Mod1"}
	moduleVersions[2] = model.ModuleVersion{Id: 2, Name: "Mod2"}
	sectionVersions = make(map[int]model.SectionVersion)
	sectionVersions[1] = model.SectionVersion{Id: 1, Name: "Sec1"}
	sectionVersions[2] = model.SectionVersion{Id: 2, Name: "Sec2"}
	fieldVersions = make(map[int]model.FieldVersion)
	fieldVersions[1] = model.FieldVersion{Id: 1, Name: "Field1"}
	fieldVersions[2] = model.FieldVersion{Id: 2, Name: "Field2"}

	actualResponse = journeyServiceUtil.ConstructFlowResponseAsList(journey, moduleVersions,
		sectionVersions, fieldVersions)

	var expectedResponse = response_dto.JourneyResponseDtoList{Name: "abc", Type: enum.FlowType.String(enum.CL), Modules: []response_dto.ResponseDTO{
		{Name: "Mod1"},
		{Name: "Mod2"},
	}, Sections: []response_dto.ResponseDTO{{Name: "Sec1"}, {Name: "Sec2"}},
		Fields: []response_dto.ResponseDTO{{Name: "Field1"}, {Name: "Field2"}}}
	assert.Equal(t, actualResponse, expectedResponse)
}
