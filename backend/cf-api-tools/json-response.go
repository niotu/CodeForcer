package cf_api_tools

import (
	"fmt"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/googlesheets"
)

type FinalJSONData struct {
	Problems     []entities.Problem `json:"problems"`
	Users        []entities.User    `json:"users"`
	CSV          []byte             `json:"csv"`
	GoogleSheets string             `json:"googleSheets"`
}

func combineStatusAndStandings(params *CFContestMethodParams, tableExtraParams ParsingParameters) (*FinalJSONData, error) {
	dataStandings, err := formattedStandings(params)
	if err != nil {
		return nil, err
	}

	dataStatus := &DataFromStatus{
		ProblemMaxPoints: make(map[string]float64),
		Users:            nil,
	}

	for _, problem := range dataStandings.Problems {
		dataStatus.ProblemMaxPoints[problem.Index] = 0.0
	}

	dataStatus, err = formattedStatus(params, dataStatus, dataStandings, tableExtraParams.SubmissionParsingMode)
	if err != nil {
		return nil, err
	}

	for _, problem := range dataStandings.Problems {
		problem.MaxPoints = dataStatus.ProblemMaxPoints[problem.Index]
	}

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

	sheet, err := fillResultsToTable(dataStandings.Name, finalJsonData, tableExtraParams)
	if err != nil {
		return nil, err
	}

	finalJsonData.CSV, err = sheet.GetSpreadsheetCSV()
	if err != nil {
		return nil, err
	}

	finalJsonData.GoogleSheets = sheet.GetSpreadsheetURL()

	return &finalJsonData, nil
}

func fillResultsToTable(name string, resultsData FinalJSONData, extraParams ParsingParameters) (*googlesheets.Spreadsheet, error) {
	var problemNames []string

	for _, p := range resultsData.Problems {
		problemNames = append(problemNames, fmt.Sprintf("Task %s - CodeForces", p.Index))
		problemNames = append(problemNames, fmt.Sprintf("Task %s - Moodle", p.Index))
	}

	csvHeaders := append([]string{"email"}, problemNames...)
	csvHeaders = append(csvHeaders, "Total points - Codeforces")
	csvHeaders = append(csvHeaders, "Total points - Moodle")

	if len(extraParams.ExtraHeaders) > 0 {
		csvHeaders = append(csvHeaders, extraParams.ExtraHeaders...)
	}
	csvHeaders = append(csvHeaders, "Feedback")

	csvData := MakeTableData(resultsData, extraParams)

	sheet, err := MakeGoogleSheet(name, csvHeaders, csvData)
	if err != nil {
		return nil, err
	}

	return sheet, nil
}
