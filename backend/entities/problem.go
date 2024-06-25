package entities

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sync"
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

type SubmissionCodeChanObject struct {
	Code  string
	Error error
}

func FetchSubmission(client *http.Client, codeChan chan SubmissionCodeChanObject, wg *sync.WaitGroup,
	groupCode string, contestId, submissionId int64) {

	submissionURL := fmt.Sprintf("https://codeforces.com/group/%s/contest/%d/submission/%d",
		groupCode, contestId, submissionId)

	req, err := http.NewRequest("GET", submissionURL, nil)
	if err != nil {
		codeChan <- SubmissionCodeChanObject{
			Code:  "",
			Error: err,
		}
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		codeChan <- SubmissionCodeChanObject{
			Code:  "",
			Error: err,
		}
		return
	}
	defer resp.Body.Close()

	// Check if response status is OK
	if resp.StatusCode != http.StatusOK {
		codeChan <- SubmissionCodeChanObject{
			Code:  "",
			Error: fmt.Errorf("failed to fetch submission, status code: %d", resp.StatusCode),
		}
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		codeChan <- SubmissionCodeChanObject{
			Code:  "",
			Error: err,
		}
	}

	var sub string

	sub = doc.Find("#program-source-text").Text()

	codeChan <- SubmissionCodeChanObject{
		Code:  sub,
		Error: nil,
	}
	close(codeChan)
	fmt.Println("finish")

	return
}
