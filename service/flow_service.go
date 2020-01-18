package service

import (
	"encoding/json"
	"flow/cache"
	"flow/db"
	"flow/enum"
	"flow/model"
	"flow/model/response_dto"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type FlowService struct {
}

func (u FlowService) GetMerchantFlows(merchantId string) response_dto.FlowResponsesDto {
	//Add log
	redisClient := cache.GetRedisClient()
	var flowsResponse response_dto.FlowResponsesDto
	if redisClient != nil {
		cachedFlow, err := redisClient.Get(merchantId).Result()
		//Add log
		if err != nil {
			//Add log
		}
		if cachedFlow == "" {
			//Add log
			flows := FetchAllMerchantFlowsFromDB(merchantId)
			flowsResponse, err := GetParsedMerchantFlows(flows)
			if err != nil {
				//Add log
			} else {
				//Add expiry time for cache entry.
				response, err := json.Marshal(flowsResponse)
				if err == nil {
					setStatus := redisClient.Set(merchantId, response, 0)
					fmt.Println(setStatus.Result())
				}
			}
			return flowsResponse
		} else {
			json.Unmarshal([]byte(cachedFlow), &flowsResponse)
		}
	} else {
		//Log redis connection failure.
	}
	return flowsResponse
}

func FetchAllMerchantFlowsFromDB(merchantId string) []model.Flow {
	fmt.Println("fetching from db")
	dbConnection := db.GetDB()
	var flows []model.Flow
	if dbConnection == nil {
		fmt.Print("failed to connect")
	} else {
		dbConnection.Where("merchantid = ?", uuid.FromStringOrNil(merchantId)).Find(&flows)
	}
	//Add log
	return flows
}

func GetParsedMerchantFlows(flows []model.Flow) (response_dto.FlowResponsesDto, error) {
	var response response_dto.FlowResponsesDto
	for _, flow := range flows {
		//Todo: write logic to parse entire json
		var flowResponseDto response_dto.FlowResponseDto
		flowResponseDto.Name = flow.Name
		flowResponseDto.Version = flow.Version
		flowResponseDto.Type = flow.Type
		flowResponseDto.Status = flow.Status
		var moduleVersionResponses []response_dto.ModuleVersionResponseDto
		var moduleList []int
		dbConnection := db.GetDB()
		err := json.Unmarshal([]byte(flow.ModuleVersions), &moduleList)
		if err != nil {
			fmt.Println(err)
		}
		for i := 0; i < len(moduleList); i++ {
			var moduleVersionResponseDto response_dto.ModuleVersionResponseDto
			module := &model.Module{Id: int(moduleList[i]), Status: enum.Active}
			dbConnection.Debug().Find(&module).Related(&module.ModuleVersions)

			moduleVersionResponseDto.Name = module.Name
			moduleVersionResponseDto.Status = module.Status

			for _, moduleVersion := range module.ModuleVersions {
				moduleVersionResponseDto.Version = moduleVersion.Version
				moduleVersionResponseDto.ExternalId = moduleVersion.ExternalId
				json.Unmarshal([]byte(moduleVersion.Properties), &moduleVersionResponseDto.Properties)

				var sectionList []int
				err := json.Unmarshal([]byte(moduleVersion.Sections), &sectionList)
				if err != nil {
					fmt.Println(err)
				}

				//Find sections against each module versions
				//Todo: this should be section

				var sectionVersionResponseDto response_dto.SectionVersionsResponseDto
				for _, sectionId := range sectionList {
					var section model.Section
					section.Id = sectionId
					dbConnection.Find(&section).Related(&section.SectionVersions).Related(&section.Fields)
					sectionVersionResponseDto.Name = section.Name
					sectionVersionResponseDto.Status = section.Status
					sectionVersionResponseDto.IsVisible = section.IsVisible

					//Find all the Sectionsversions and map them to fields
					for _, sectionVersion := range section.SectionVersions {
						fmt.Println(sectionVersion)
						for _, field := range section.Fields {
							var fieldModel model.Field
							fieldModel.Id = field.Id
							dbConnection.Find(&fieldModel).Related(&fieldModel.FieldVersions)
							var fieldVersionsResponseDto response_dto.FieldVersionsResponseDto
							fieldVersionsResponseDto.Name = fieldModel.Name
							fieldVersionsResponseDto.Status = fieldModel.Status
							fieldVersionsResponseDto.IsVisible = fieldModel.IsVisible
							for _, fieldVersion := range fieldModel.FieldVersions {
								fieldVersionsResponseDto.Version = fieldVersion.Version
								fieldVersionsResponseDto.ExternalId = fieldVersion.ExternalId
								fieldVersionsResponseDto.IsVisible = fieldVersion.IsVisible
								sectionVersionResponseDto.Fields = append(sectionVersionResponseDto.Fields, fieldVersionsResponseDto)
							}
						}

					}
					moduleVersionResponseDto.Sections = append(moduleVersionResponseDto.Sections, sectionVersionResponseDto)
				}
			}

			moduleVersionResponses = append(moduleVersionResponses, moduleVersionResponseDto)
		}

		flowResponseDto.Modules = moduleVersionResponses
		response.FlowResponses = append(response.FlowResponses, flowResponseDto)
	}
	return response, nil
}
