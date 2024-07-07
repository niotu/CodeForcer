package entities

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sync"
	"time"
)

type Submission struct {
	Index          string
	Solution       string
	Points         float64
	SubmissionId   int64
	ProgramLang    string
	SubmissionTime int64
	Late           bool
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

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
}

var count = 0
var userAgentIdx = 0
var currUserAgent = userAgents[0]

func getNextUserAgent() string {

	//return userAgents[rand.Intn(len(userAgents))]
	userAgentIdx = (userAgentIdx + 1) % len(userAgents)
	currUserAgent = userAgents[userAgentIdx]
	return currUserAgent
}

func FetchSubmission(client *http.Client, codeChan chan SubmissionCodeChanObject, mutex *sync.Mutex,
	groupCode string, contestId, submissionId int64) {

	submissionURL := fmt.Sprintf("https://codeforces.com/group/%s/contest/%d/submission/%d",
		groupCode, contestId, submissionId)

	req, err := http.NewRequest("GET", submissionURL, nil)
	if err != nil {
		fmt.Println(err)
		codeChan <- SubmissionCodeChanObject{
			Code:  "",
			Error: err,
		}
		close(codeChan)
		return
	}
	req.Header.Set("User-Agent", currUserAgent)
	if count%15 == 0 {
		req.Header.Set("User-Agent", getNextUserAgent())
		fmt.Println("change agent!")
	}

	fmt.Println(time.Now())
	mutex.Lock()
	resp, err := client.Do(req)
	mutex.Unlock()

	if err != nil {
		fmt.Println(err)
		codeChan <- SubmissionCodeChanObject{
			Code:  "",
			Error: err,
		}
		close(codeChan)
		return
	}
	defer resp.Body.Close()

	// Check if response status is OK
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Request.URL, "\n", resp.StatusCode)
		codeChan <- SubmissionCodeChanObject{
			Code:  "",
			Error: fmt.Errorf("failed to fetch submission, status code: %d", resp.StatusCode),
		}
		close(codeChan)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		codeChan <- SubmissionCodeChanObject{
			Code:  "",
			Error: err,
		}
		close(codeChan)
		return
	}

	var sub string

	sub = doc.Find("#program-source-text").Text()

	codeChan <- SubmissionCodeChanObject{
		Code:  sub,
		Error: nil,
	}
	close(codeChan)

	fmt.Println("fin ", count)
	count++

	return
}
