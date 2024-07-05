package cf_api_tools

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/googlesheets"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"sort"
	"strconv"
	"strings"
)

func EntitiesToJSON[T any](jsonData T) []byte {
	data, err := json.Marshal(jsonData)
	if err != nil {
		logger.Error(err)
	}

	return data
}

func MakeCSVFile(data FinalJSONData, headers []string) (*bytes.Buffer, [][]string) {
	//buff, _ := os.OpenFile("report.csv", os.O_CREATE|os.O_TRUNC, 0606)

	buff := new(bytes.Buffer)

	writer := csv.NewWriter(buff)
	writer.Comma = ';'

	_ = writer.Write(headers)

	var rows [][]string

	count := 1

	for _, user := range data.Users {
		var comment []string

		points := 0.0
		for idx, submission := range user.Solutions {
			// TODO: add logic for calculating weights of tasks

			points += submission.Points
			id := submission.SubmissionId
			if id == -1 {
				comment = append(comment, fmt.Sprintf("%s: %d (no submission)", idx, 0))
			} else {
				comment = append(comment, fmt.Sprintf("%s: %d", idx, int(submission.Points)))
			}
		}

		sort.Strings(comment)

		row := []string{fmt.Sprintf("User%d", count), strconv.Itoa(int(points)), strings.Join(comment, "; ")}
		rows = append(rows, row)
		count++

		_ = writer.Write(row)
	}

	writer.Flush()

	return buff, rows
}

func MakeGoogleSheet(name string, headers []string, data [][]string) (string, error) {
	ss, err := googlesheets.CreateSpreadsheet(name, "ramazanatzuf10@gmail.com")
	if err != nil {
		return "", err
	}

	err = ss.WriteHeaders(googlesheets.ToInterfaceSlice(headers))
	if err != nil {
		return "", err
	}

	var rows [][]interface{}
	for _, row := range data {
		rows = append(rows, googlesheets.ToInterfaceSlice(row))
	}

	err = ss.WriteData(rows)
	if err != nil {
		return "", err
	}

	return ss.GetSpreadsheetURL(), nil
}
