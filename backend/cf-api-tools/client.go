package cf_api_tools

import (
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	CONTEST_STATUS    = "contest.status"
	CONTEST_STANDINGS = "contest.standings"
)

type Client struct {
	apiKey    string
	apiSecret string
	handle    string
	password  string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SetApiKey(apiKey string) {
	c.apiKey = apiKey
}

func (c *Client) SetApiSecret(apiSecret string) {
	c.apiSecret = apiSecret
}

func (c *Client) SetHandle(handle string) {
	c.handle = handle
}

func (c *Client) SetPassword(password string) {
	c.password = password
}

func (c *Client) GetGroupsList(client *http.Client) ([]entities.Group, *http.Client) {
	var err error

	if client == nil {
		client, err = entities.Login(c.handle, c.password)
		if err != nil {
			log.Printf("Login failed: %v", err)
		}
	}

	groups, err := entities.FetchGroups(client)
	if err != nil {
		log.Printf("Failed to fetch groups: %v", err)
	}

	return groups, client
}

func (c *Client) GetContestsList(client *http.Client, groupCode string) ([]entities.Contest, *http.Client) {
	var err error

	if client == nil {
		client, err = entities.Login(c.handle, c.password)
		if err != nil {
			log.Printf("Login failed: %v", err)
		}
	}

	contests, err := entities.FetchContests(client, groupCode)
	if err != nil {
		log.Printf("Failed to fetch groups: %v", err)
	}

	return contests, client
}

func (c *Client) GetSubmissionCode(client *http.Client, ch chan entities.SubmissionCodeChanObject,
	mutex *sync.Mutex, groupCode string, contestId, submissionId int64) {
	var err error

	if client == nil {
		client, err = entities.Login(c.handle, c.password)
		if err != nil {
			log.Fatalf("Login failed: %v", err)
		}
	}

	entities.FetchSubmission(client, ch, mutex, groupCode, contestId, submissionId)
	if err != nil {
		log.Printf("Failed to fetch submission: %v", err)
	}
}

func (c *Client) GetStatistics(client *http.Client, groupCode string, contestId int64, count int, weights []int) []byte {
	params := &CFContestMethodParams{
		GroupCode: groupCode,
		ContestId: contestId,
		AsManager: true,
		ApiKey:    c.apiKey,
		ApiSecret: c.apiSecret,
		Time:      time.Now().Unix(),
		Count:     count,
	}

	finalData := parseAndFormEntities(params)

	//submissionCodeChan := make(map[int64]chan entities.SubmissionCodeChanObject)

	//wg := &sync.WaitGroup{}
	//defer wg.Wait()
	//
	//mutex := &sync.Mutex{}

	//var err error
	//if client == nil {
	//	client, err = entities.Login(c.handle, c.password)
	//	if err != nil {
	//		log.Fatalf("Login failed: %v", err)
	//	}
	//}
	//
	//for _, u := range finalData.Users {
	//	for _, s := range u.Solutions {
	//		if s.SubmissionId != -1 {
	//			submissionCodeChan[s.SubmissionId] = make(chan entities.SubmissionCodeChanObject, 1)
	//			ch := submissionCodeChan[s.SubmissionId]
	//
	//			wg.Add(1)
	//			go func(group *sync.WaitGroup, channel chan entities.SubmissionCodeChanObject, id int64) {
	//				defer group.Done()
	//
	//				c.GetSubmissionCode(client, channel, mutex,
	//					groupCode, contestId, id)
	//			}(wg, ch, s.SubmissionId)
	//			time.Sleep(2100 * time.Millisecond)
	//		}
	//	}
	//}
	//
	//wg.Wait()
	//
	//for _, u := range finalData.Users {
	//	for key, sub := range u.Solutions {
	//		if sub.SubmissionId == -1 {
	//			continue
	//		}
	//		code := <-submissionCodeChan[sub.SubmissionId]
	//
	//		if code.Error != nil {
	//			fmt.Println(code.Error)
	//		} else {
	//			u.Solutions[key].Solution = code.Code
	//		}
	//	}
	//}

	return EntitiesToJSON(finalData)

}
