package cf_api_tools

import (
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"log"
)

type FinalJSONData struct {
	Problems     []entities.Problem `json:"problems"`
	Users        []entities.User    `json:"users"`
	CSV          []byte             `json:"csv"`
	GoogleSheets string             `json:"googleSheets"`
}

func parseAndFormEntities(params *CFContestMethodParams) *FinalJSONData {
	standings := getContestStandings(params)
	dataStandings, err := parseContestStandings(standings)
	if err != nil {
		log.Println(err)
	}

	dataStatus := &DataFromStatus{
		ProblemMaxPoints: make(map[string]float64),
		Users:            nil,
	}

	for _, problem := range dataStandings.Problems {
		dataStatus.ProblemMaxPoints[problem.Index] = 0.0
	}

	status := getContestStatus(params)
	dataStatus, err = parseContestStatus(status, dataStatus, dataStandings)

	for _, problem := range dataStandings.Problems {
		problem.MaxPoints = dataStatus.ProblemMaxPoints[problem.Index]
	}

	//ProblemListToJSON(dataStandings.Problems)

	var u []entities.User
	for _, v := range dataStatus.Users {
		u = append(u, *v)
	}

	var p []entities.Problem
	for _, v := range dataStandings.Problems {
		p = append(p, *v)
	}

	finalJsonData := FinalJSONData{
		Problems: p,
		Users:    u,
	}

	csvHeaders := []string{"handle", "points", "comment"}
	csvBuff, csvData := MakeCSVFile(finalJsonData, csvHeaders)

	finalJsonData.CSV = csvBuff.Bytes()

	sheetURL, err := MakeGoogleSheet(dataStandings.Name, csvHeaders, csvData)
	if err != nil {
		log.Println(err)
	}
	finalJsonData.GoogleSheets = sheetURL

	return &finalJsonData
}
