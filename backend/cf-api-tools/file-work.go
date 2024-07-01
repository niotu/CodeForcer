package cf_api_tools

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/googlesheets"
	"log"
	"sort"
	"strconv"
	"strings"
)

func EntitiesToJSON[T any](jsonData T) []byte {
	data, err := json.Marshal(jsonData)
	if err != nil {
		log.Println(err)
	}

	return data
}

func MakeCSVFile(data FinalJSONData, headers []string) (*bytes.Buffer, [][]string) {
	//buff, _ := os.OpenFile("report.csv", os.O_CREATE|os.O_TRUNC, 0606)

	buff := new(bytes.Buffer)

	writer := csv.NewWriter(buff)
	writer.Comma = ';'

	writer.Write(headers)

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

		writer.Write(row)
	}

	writer.Flush()

	return buff, rows
}

func MakeGoogleSheet(name string, headers []string, data [][]string) (string, error) {
	_, err := googlesheets.CreateSpreadsheet(name)
	if err != nil {
		return "", err
	}

	googlesheets.WriteHeaders(googlesheets.ToInterfaceSlice(headers))

	var rows [][]interface{}
	for _, row := range data {
		rows = append(rows, googlesheets.ToInterfaceSlice(row))
	}

	googlesheets.WriteData(rows)

	return googlesheets.GetSpreadsheetURL(), nil
}
