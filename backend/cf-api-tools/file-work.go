package cf_api_tools

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
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

func MakeCSVFile(data FinalJSONData) *bytes.Buffer {
	//buff, _ := os.OpenFile("report.csv", os.O_CREATE|os.O_TRUNC, 0606)

	buff := new(bytes.Buffer)

	writer := csv.NewWriter(buff)

	headers := []string{"handle", "points", "comment"}

	writer.Write(headers)

	for _, user := range data.Users {
		var comment []string

		points := 0.0
		for idx, submission := range user.Solutions {
			points += submission.Points
			id := submission.SubmissionId
			if id == -1 {
				comment = append(comment, fmt.Sprintf("%s: %d (no submission)", idx, 0))
			} else {
				comment = append(comment, fmt.Sprintf("%s: %d;", idx, int(submission.Points)))
			}
		}

		sort.Strings(comment)

		row := []string{user.Handle, strconv.Itoa(int(points)), strings.Join(comment, "; ")}

		writer.Write(row)
	}

	writer.Flush()

	return buff

}
