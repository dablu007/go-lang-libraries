package service

import (
	"flow/config"
	"flow/enum"
	"flow/logger"
	"flow/model"
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
	var controller = gomock.NewController(t)
	defer controller.Finish()

	var journey model.Journey = model.Journey{}
	var moduleRepo = mock_repository.NewMockModuleRepository(controller)
	var sectionRepo = mock_repository.NewMockSectionRepository(controller)
	var fieldRepo = mock_repository.NewMockFieldRepository(controller)
	var mapUtil = mock_utility.NewMockMapUtility(controller)

	var journeyService = &service.JourneyServiceUtil{
		FieldRepository:   fieldRepo,
		ModuleRepository:  moduleRepo,
		SectionRepository: sectionRepo,
		MapUtil:           mapUtil,
	}

	//test with no journey
	mapUtil.EXPECT().GetKeyListFromKeyValueMap(map[int]bool{}).Return([]int{})
	moduleRepo.EXPECT().FetchModuleVersions(enum.Active, []int{}).Return([]model.ModuleVersion{})
	sectionRepo.EXPECT().FetchSectionVersions(enum.Active, []int{}).Return([]model.SectionVersion{})
	fieldRepo.EXPECT().FetchFieldVersions(enum.Active, []int{}).Return([]model.FieldVersion{})

	var actualModuleVersionsMap map[int]model.ModuleVersion
	var actualSectionVersionsMap map[int]model.SectionVersion
	var actualFieldVersionsMap map[int]model.FieldVersion
	var actualCompleteModuleVersionNumberList map[int]bool
	var actualCompleteSectionVersionNumberList map[int]bool
	var actualCompleteFieldVersionNumberList map[int]bool
	actualModuleVersionsMap, actualSectionVersionsMap, actualFieldVersionsMap,
		actualCompleteModuleVersionNumberList, actualCompleteSectionVersionNumberList, actualCompleteFieldVersionNumberList =
		journeyService.GetModuleSectionAndFieldVersionsAndActiveVersionNumberList(journey)
	fmt.Println("length", len(actualModuleVersionsMap))
	assert.Equal(t, len(actualModuleVersionsMap), 0)
	assert.Equal(t, len(actualSectionVersionsMap), 0)
	assert.Equal(t, len(actualFieldVersionsMap), 0)
	assert.Equal(t, len(actualCompleteModuleVersionNumberList), 0)
	assert.Equal(t, len(actualCompleteSectionVersionNumberList), 0)
	assert.Equal(t, len(actualCompleteFieldVersionNumberList), 0)

}

func TestConstructJourneysResponse(t *testing.T) {

}

func TestConstructFlowResponseWithModuleFieldSection(t *testing.T) {

}
