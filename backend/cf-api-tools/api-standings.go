package cf_api_tools

import (
	"encoding/json"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
)

type DataFromStandings struct {
	Name             string
	Problems         []*entities.Problem
	DurationSeconds  int64
	StartTimeSeconds int64
}

func getContestStandings(params *CFContestMethodParams) (interface{}, error) {
	api := NewApiRequest(ContestStandings, params)

	resp, err := api.MakeApiRequest()
	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		logger.Error(err)
		return nil, ApiRequestError
	}

	return data, nil
}

func parseContestStandings(data interface{}) (*DataFromStandings, error) {
	var problems []*entities.Problem

	if resp := data.(map[string]interface{}); resp["status"].(string) == "FAILED" {
		comment := resp["comment"].(string)

		if err := checkResponseError(comment); err != nil {
			return nil, err
		}
	}

	result := data.(map[string]interface{})["result"].(map[string]interface{})
	problemsInfo := result["problems"].([]interface{})

	for _, problemJson := range problemsInfo {
		problems = append(problems, &entities.Problem{
			Name:  problemJson.(map[string]interface{})["name"].(string),
			Index: problemJson.(map[string]interface{})["index"].(string),
		})
	}

	contest := result["contest"].(map[string]interface{})

	name := contest["name"].(string)
	durationSeconds := contest["durationSeconds"].(float64)
	startTimeSeconds := contest["startTimeSeconds"].(float64)

	return &DataFromStandings{
		Name:             name,
		Problems:         problems,
		DurationSeconds:  int64(durationSeconds),
		StartTimeSeconds: int64(startTimeSeconds),
	}, nil
}
