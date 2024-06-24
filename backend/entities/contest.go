package entities

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Contest struct {
	Id               int
	Name             string
	GroupCode        string
	ContestLink      string
	DurationString   string
	StartTimeString  string
	DurationSeconds  int64
	StartTimeSeconds int64
	Problems         []Problem
}

// fetchContests fetches the contests of the current group
func fetchContests(client *http.Client, groupCode string) ([]Contest, error) {
	groupURL := "https://codeforces.com/group/" + groupCode

	contestsURL := groupURL + "/contests"
	req, err := http.NewRequest("GET", contestsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch groups, status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var contests []Contest

	doc.Find(".highlighted-row").Each(func(i int, s *goquery.Selection) {
		contestIdString, exists := s.Attr("data-contestid")
		contestId, _ := strconv.Atoi(contestIdString)

		contestName := strings.Split(
			strings.TrimSpace(
				s.Find("td:nth-child(1)").Text()), "\n",
		)[0]

		if exists {
			contestStartTimeString := strings.TrimSpace(s.Find("a[target='_blank']").Text())
			contestDurationString := strings.TrimSpace(s.Find("td:nth-child(3)").Text())

			contestStartTime, _ := time.Parse("Jan/02/2006 15:04", contestStartTimeString)
			parts := strings.Split(contestDurationString, ":")
			days, _ := strconv.Atoi(parts[0])
			hours, _ := strconv.Atoi(parts[1])
			minutes, _ := strconv.Atoi(parts[2])

			contestDuration := (days * 86400) + (hours * 3600) + (minutes * 60)

			contests = append(contests, Contest{
				Id:               contestId,
				Name:             contestName,
				GroupCode:        groupCode,
				DurationSeconds:  int64(contestDuration),
				StartTimeSeconds: contestStartTime.Unix(),
				DurationString:   contestDurationString,
				StartTimeString:  contestStartTimeString,
				ContestLink:      groupURL + "/contest/" + strconv.Itoa(contestId),
				Problems:         nil,
			})
		}
	})

	return contests, nil
}

func ContestListToJSON(contests []Contest) []byte {
	data, err := json.Marshal(contests)
	if err != nil {
		log.Fatal(err)
	}

	file, _ := os.OpenFile("contests.json", os.O_CREATE|os.O_TRUNC, 0606)
	file.Write(data)

	return data
}
