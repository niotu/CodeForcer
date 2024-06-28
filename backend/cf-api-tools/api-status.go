package cf_api_tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"log"
	"strings"
)

type DataFromStatus struct {
	ProblemMaxPoints map[string]float64
	Users            map[string]*entities.User
}

func getContestStatus(params *CFContestMethodParams) interface{} {
	api := NewApiRequest(CONTEST_STATUS, params)
	fmt.Println(api.ApiSig.Rand, params.Time)

	resp, err := api.MakeApiRequest()
	if err != nil {
		fmt.Println(err)
	}

	var data interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Println(err)
	}

	return data
}

func parseContestStatus(data interface{}, dataStatus *DataFromStatus, dataStandings *DataFromStandings) (*DataFromStatus, error) {
	db := make(map[string]*entities.User)

	if resp := data.(map[string]interface{}); resp["status"].(string) == "FAILED" {
		comment := resp["comment"].(string)

		errorMsg := ""

		if strings.Contains(comment, "asManager") {
			errorMsg = "You are not the manager of contest or group. Be sure that on the page of contest you selected the following:\n" +
				"Administration (block on the right) -> Enable manager mode."
		}

		return nil, errors.New(errorMsg)
	}

	result := data.(map[string]interface{})["result"].([]interface{})

	if dataStatus == nil {
		dataStatus = &DataFromStatus{}
	}

	for _, resultElem := range result {
		submissionJson := resultElem.(map[string]interface{})

		username := submissionJson["author"].(map[string]interface{})["members"].([]interface{})[0].(map[string]interface{})["handle"].(string)
		problemIdx := submissionJson["problem"].(map[string]interface{})["index"].(string)

		//go fetchSubmission(nil)

		if _, ok := db[username]; !ok {
			submissions := make(map[string]*entities.Submission)
			for _, p := range dataStandings.Problems {
				submissions[p.Index] = &entities.Submission{
					Index:          p.Index,
					SubmissionId:   -1,
					ProgramLang:    submissionJson["programmingLanguage"].(string),
					SubmissionTime: int64(submissionJson["creationTimeSeconds"].(float64)),
				}

			}

			db[username] = &entities.User{
				Handle:    username,
				Solutions: submissions,
			}
		}
		if s, _ := db[username].Solutions[problemIdx]; s.SubmissionId == -1 {
			_, exists := submissionJson["points"]
			submissionPoints := 0.0
			if exists {
				submissionPoints = submissionJson["points"].(float64)
			}

			problemVerdict := submissionJson["verdict"].(string)
			if value := dataStatus.ProblemMaxPoints[problemIdx]; problemVerdict == "OK" && value == 0.0 {
				dataStatus.ProblemMaxPoints[problemIdx] = submissionPoints
			}

			id := int64(submissionJson["id"].(float64))

			db[username].Solutions[problemIdx].Points = submissionPoints
			db[username].Solutions[problemIdx].SubmissionId = id

			//db[username].Solutions[problemIdx].Solution = go entities.FetchSubmission(cf, groupCode, contestId, id)
		}
	}

	dataStatus.Users = db

	return dataStatus, nil
}
