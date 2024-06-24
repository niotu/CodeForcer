package entities

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
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

func fetchSubmission(client *http.Client, groupCode string, contestId, submissionId int) (string, error) {
	submissionURL := fmt.Sprintf("https://codeforces.com/group/%s/contest/%d/submission/%d",
		groupCode, contestId, submissionId)

	req, err := http.NewRequest("GET", submissionURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if response status is OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch submission, status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	var sub string

	sub = doc.Find("#program-source-text").Text()

	return sub, nil
}
