package entities

import (
	"encoding/json"
	"log"
	"os"
)

type Submission struct {
	Index          string
	Solution       string
	Points         float64
	SubmissionId   int64
	ProgramLang    string
	SubmissionTime int64
}

type Problem struct {
	Name      string
	Index     string
	MaxPoints float64
}

func ProblemListToJSON(problems []Problem) []byte {
	data, err := json.Marshal(problems)
	if err != nil {
		log.Fatal(err)
	}

	file, _ := os.OpenFile("problems.json", os.O_CREATE|os.O_TRUNC, 0606)
	file.Write(data)

	return data
}
