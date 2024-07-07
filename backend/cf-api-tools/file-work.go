package cf_api_tools

import (
	"fmt"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/googlesheets"
	"strconv"
	"strings"
)

type ParsingParameters struct {
	TasksWeights          []int
	ExtraHeaders          []string
	LatePenalty           int
	LateTime              int64
	SubmissionParsingMode string
}

func feedbackFormulaPattern(start int) string {
	startLetter := string(rune('A' + start - 1))
	return fmt.Sprintf("& CHAR(10)&CHAR(10) & IF(COLUMN() <= %d, \"\", JOIN(CHAR(10)&CHAR(10), ARRAYFORMULA(\"$%s$1\":INDIRECT(ADDRESS(1, COLUMN()-1)) & \": \" & INDIRECT(\"$%s\"&ROW()):INDIRECT(ADDRESS(ROW(), COLUMN()-1)))))",
		start, startLetter, startLetter)
}

var FeedbackFormula string

func MakeTableData(resultsData FinalJSONData, extraParams ParsingParameters) [][]string {
	FeedbackFormula = feedbackFormulaPattern(4 + len(resultsData.Problems)*2)

	var rows [][]string

	count := 1

	for _, user := range resultsData.Users {
		totalCF := 0.0
		totalMoodle := 0.0
		var points []string
		i := 0

		var feedbackPart []string

		for _, submission := range user.Solutions {

			// codeforces task points
			points = append(points, strconv.Itoa(int(submission.Points)))
			// converted to moodle according to weight
			moodlePoints := submission.Points / resultsData.Problems[i].MaxPoints * float64(extraParams.TasksWeights[i])

			id := submission.SubmissionId
			taskStatus := ""
			if id == -1 {
				taskStatus = "(no submission)"
			} else if submission.Late && submission.SubmissionTime <= extraParams.LateTime {
				taskStatus = fmt.Sprintf("(late submission: -%d%%)", extraParams.LatePenalty)
				moodlePoints *= 1.0 - float64(extraParams.LatePenalty)/100
			} else if submission.Late && submission.SubmissionTime > extraParams.LateTime {
				taskStatus = "(no submission)"
				submission.Solution = ""
			}

			points = append(points, strconv.Itoa(int(moodlePoints)))

			feedbackPart = append(feedbackPart, fmt.Sprintf("Task %s: %d/%d %s",
				submission.Index, int(submission.Points), int(resultsData.Problems[i].MaxPoints), taskStatus))

			totalCF += submission.Points
			totalMoodle += moodlePoints

			i++
		}

		row := append([]string{fmt.Sprintf("User%d", count)}, points...)
		row = append(row, strconv.Itoa(int(totalCF)), strconv.Itoa(int(totalMoodle)))
		row = append(row, "=\"Passing test:\n"+strings.Join(feedbackPart, "; ")+"\""+FeedbackFormula)

		rows = append(rows, row)

		count++

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
