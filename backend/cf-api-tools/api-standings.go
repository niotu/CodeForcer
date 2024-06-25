package cf_api_tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"log"
	"strings"
	"time"
)

type DataFromStandings struct {
	Problems         []*entities.Problem
	DurationSeconds  int64
	StartTimeSeconds int64
}

func getContestStandings(c *Client, groupCode string, contestId int64) interface{} {
	params := &ContestStandingsRequestParams{
		GroupCode: groupCode,
		ContestId: contestId,
		AsManager: true,
		ApiKey:    c.apiKey,
		ApiSecret: c.apiSecret,
		Time:      time.Now().Unix(),
		Count:     1,
	}
	api := NewApiRequest(CONTEST_STANDINGS, params)

	resp, err := api.MakeApiRequest()
	if err != nil {
		log.Println(err)
	}

	var data interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Println(err)
	}

	return data
}

func parseContestStandings(data interface{}) (*DataFromStandings, error) {
	var problems []*entities.Problem

	if resp := data.(map[string]interface{}); resp["status"].(string) == "FAILED" {
		comment := resp["comment"].(string)

		errorMsg := ""

		if strings.Contains(comment, "asManager") {
			errorMsg = "You are not the manager of contest or group. Be sure that on the page of contest you selected the following:\n" +
				"Administration (block on the right) -> Enable manager mode."
		}

		return nil, errors.New(errorMsg)
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
	durationSeconds := contest["durationSeconds"].(float64)
	startTimeSeconds := contest["startTimeSeconds"].(float64)

	// test
	fmt.Println(problems)

	return &DataFromStandings{
		Problems:         problems,
		DurationSeconds:  int64(durationSeconds),
		StartTimeSeconds: int64(startTimeSeconds),
	}, nil
}
