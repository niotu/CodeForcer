package cf_api_tools

import (
	"fmt"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/db"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/googlesheets"
	"strconv"
	"strings"
)

type ParsingParameters struct {
	TasksWeights          []int
	ExtraHeaders          []string
	LatePenalty           int
	LateEndSeconds        int64
	LateDurationSeconds   int64
	SubmissionParsingMode string
}

func feedbackFormulaPattern(startCol int) string {
	startLetter := string(rune('A' + startCol))
	return fmt.Sprintf("& CHAR(10)&CHAR(10) & IF(COLUMN() <= %d, \"\", JOIN(CHAR(10)&CHAR(10), ARRAYFORMULA(INDIRECT(\"$%s$1\"):INDIRECT(ADDRESS(1, COLUMN()-1)) & \": \" & INDIRECT(\"$%s\"&ROW()):INDIRECT(ADDRESS(ROW(),COLUMN()-1)))))",
		startCol+1, startLetter, startLetter)
}

func totalFormulaPattern(start int) string {
	startLetter := string(rune('A' + start - 1))
	return fmt.Sprintf("=SUM(INDIRECT(\"%s\" & ROW() & \":\" & ADDRESS(ROW(), COLUMN() - 1)))",
		startLetter)
}

func MakeTableData(resultsData FinalJSONData, extraParams ParsingParameters, mandatoryCols int) [][]string {
	TotalField := totalFormulaPattern(mandatoryCols)

	mapHandleToEmail := db.GetUsers()

	var rows [][]string

	for userIdx, user := range resultsData.Users {
		totalCF := 0.0
		totalMoodle := 0.0
		var points []string
		i := 0

		var feedbackPart []string

		for _, p := range resultsData.Problems {
			submission := user.Solutions[p.Index]
			// codeforces task points
			points = append(points, strconv.Itoa(int(submission.Points)))
			// converted to moodle according to weight
			moodlePoints := submission.Points / resultsData.Problems[i].MaxPoints * float64(extraParams.TasksWeights[i])

			id := submission.SubmissionId
			taskStatus := ""
			if id == -1 {
				taskStatus = "(no submission)"
			} else if submission.Late && submission.SubmissionTime <= extraParams.LateEndSeconds {
				taskStatus = fmt.Sprintf("(late submission: -%d%%)", extraParams.LatePenalty)
				moodlePoints *= 1.0 - float64(extraParams.LatePenalty)/100
			} else if submission.Late && submission.SubmissionTime > extraParams.LateEndSeconds {
				taskStatus = "(submission after extended deadline => 0 points)"
				moodlePoints = 0.0
			}

			points = append(points, strconv.Itoa(int(moodlePoints)))

			feedbackPart = append(feedbackPart, fmt.Sprintf("Task %s: %d/%d %s",
				submission.Index, int(submission.Points), int(resultsData.Problems[i].MaxPoints), taskStatus))

			totalCF += submission.Points
			totalMoodle += moodlePoints

			i++
		}

		userEmail := mapHandleToEmail[user.Handle]

		row := append([]string{fmt.Sprintf("User%d", userIdx+1), userEmail}, points...)
		//row := append([]string{user.Handle, userEmail}, points...)
		row = append(row, strconv.Itoa(int(totalCF)), strconv.Itoa(int(totalMoodle)))
		if len(extraParams.ExtraHeaders) > 0 {
			row = append(row, make([]string, len(extraParams.ExtraHeaders))...)
		}
		row = append(row, TotalField)
		row = append(row, "=\"Passing test:\n"+strings.Join(feedbackPart, "; ")+"\""+feedbackFormulaPattern(mandatoryCols))

		rows = append(rows, row)

	}

	return rows
}

func MakeGoogleSheet(name string, headers []string, data [][]string) (*googlesheets.Spreadsheet, error) {
	ss, err := googlesheets.CreateSpreadsheet(name, "ramazanatzuf10@gmail.com")
	if err != nil {
		return nil, err
	}

	err = ss.WriteHeaders(googlesheets.ToInterfaceSlice(headers))
	if err != nil {
		return nil, err
	}

	var rows [][]interface{}
	for _, row := range data {
		rows = append(rows, googlesheets.ToInterfaceSlice(row))
	}

	err = ss.WriteData(rows)
	if err != nil {
		return nil, err
	}

	return ss, nil
}
