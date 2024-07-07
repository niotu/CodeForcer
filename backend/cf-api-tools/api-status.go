package cf_api_tools

import (
	"encoding/json"
	"errors"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"strings"
)

type DataFromStatus struct {
	ProblemMaxPoints map[string]float64
	Users            map[string]*entities.User
}

const (
	BestSolutionMode = "best"
	LastSolutionMode = "last"
)

func checkResponseError(comment string) error {
	errorMsg := ""

	if strings.Contains(comment, "Incorrect API key") {
		errorMsg = "Incorrect API key"
	} else if strings.Contains(comment, "asManager") || strings.Contains(comment, "Incorrect signature") {
		errorMsg = "Incorrect API secret, please check it on the settings page\n"
		errorMsg += "\tOR\n"
		errorMsg += "You are not the manager of the contest. Be sure that on the page of contest you selected the following:\n" +
			"Administration (block on the right) -> Enable manager mode."
		errorMsg += "\tOR\n"
		errorMsg += "Check the correctness of entered URL"
	} else {
		return nil
	}

	return errors.New(errorMsg)
}

func getContestStatus(params *CFContestMethodParams) (interface{}, error) {
	api := NewApiRequest(ContestStatus, params)

	resp, err := api.MakeApiRequest()
	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseContestStatus(data interface{}, dataStatus *DataFromStatus, dataStandings *DataFromStandings, mode string) (*DataFromStatus, error) {
	db := make(map[string]*entities.User)

	if mode != LastSolutionMode && mode != BestSolutionMode {
		return nil, errors.New("incorrect mode")
	}

	if resp := data.(map[string]interface{}); resp["status"].(string) == "FAILED" {
		comment := resp["comment"].(string)

		if err := checkResponseError(comment); err != nil {
			return nil, err
		}
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
					Late:           false,
					ProgramLang:    submissionJson["programmingLanguage"].(string),
					SubmissionTime: int64(submissionJson["creationTimeSeconds"].(float64)),
				}

			}

			db[username] = &entities.User{
				Handle:    username,
				Solutions: submissions,
			}
		}

		_, exists := submissionJson["points"]
		submissionPoints := 0.0
		if exists {
			submissionPoints = submissionJson["points"].(float64)
		}

		currUser := db[username]

		if s := currUser.Solutions[problemIdx]; s.SubmissionId == -1 || (mode == BestSolutionMode &&
			submissionPoints > currUser.Solutions[problemIdx].Points) {

			problemVerdict := submissionJson["verdict"].(string)
			if value := dataStatus.ProblemMaxPoints[problemIdx]; problemVerdict == "OK" && value == 0.0 {
				dataStatus.ProblemMaxPoints[problemIdx] = submissionPoints
			}

			id := int64(submissionJson["id"].(float64))

			if s.SubmissionTime > dataStandings.StartTimeSeconds+dataStandings.DurationSeconds {
				currUser.Solutions[problemIdx].Late = true
			}
			currUser.Solutions[problemIdx].Points = submissionPoints
			currUser.Solutions[problemIdx].SubmissionId = id
		}
	}

	dataStatus.Users = db

	return dataStatus, nil
}

func formattedStatus(params *CFContestMethodParams, dataStatus *DataFromStatus, dataStandings *DataFromStandings, mode string) (*DataFromStatus, error) {
	status, err := getContestStatus(params)
	if err != nil {
		return nil, err
	}
	dataStatus, err = parseContestStatus(status, dataStatus, dataStandings, mode)
	if err != nil {
		return nil, err
	}

	return dataStatus, nil
}
