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

func GetAllFlowsByMerchantId(merchantId string, tenantId string, channel string) response_dto.FlowResponseDto {
	fmt.Println("Trying to fetch from redis cache.")
	redisClient := cache.GetRedisClient()
	var merchantFlows []int
	var response response_dto.FlowResponseDto
	if redisClient != nil {
		flows, err := redisClient.Get(merchantId).Result()
		fmt.Println("flows", flows)
		if err != nil {
			fmt.Println("Failed to fetch from redis with error: ", err)
		}
		if flows == "" {
			fmt.Println("Fetching all from db")
			flow := FetchAllMerchantFlowsFromDB(merchantId)
			json.Unmarshal([]byte(flow.ModuleVersions), &merchantFlows)
			// merchantFlows = flow.ModuleVersions
			if merchantFlows != nil {
				flowResponseDto, err := ParseMerchantFlowsJsonValues(flow)
				if err != nil {

					fmt.Println(err)
				}
				redisClient.Set(merchantId, flowResponseDto, 0)
				response = flowResponseDto
			}
		}
	} else {
		fmt.Println("Failed to connect with Redis.", redisClient)
	}
	return response
}

func FetchAllMerchantFlowsFromDB(merchantId string) model.Flow {
	fmt.Println("fetching from db")
	dbConnection := db.GetDB()
	var flow *model.Flow
	if dbConnection == nil {
		fmt.Print("failed to connect")
		return *flow
	}
	flow = &model.Flow{MerchantId: uuid.FromStringOrNil(merchantId)}
	dbConnection.First(&flow)
	fmt.Print(flow.Id, flow.Name, flow.CreatedOn)
	return *flow
}

func ParseMerchantFlowsJsonValues(flow model.Flow) (response_dto.FlowResponseDto, error) {
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

	return flowResponseDto, nil
}
