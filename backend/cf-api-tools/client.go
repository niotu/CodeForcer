package cf_api_tools

import (
	"fmt"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"log"
	"net/http"
	"sync"
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
	wg *sync.WaitGroup, groupCode string, contestId, submissionId int64) {
	fmt.Println(submissionId)
	var err error

	if client == nil {
		client, err = entities.Login(c.handle, c.password)
		if err != nil {
			log.Fatalf("Login failed: %v", err)
		}
	}

	entities.FetchSubmission(client, ch, wg, groupCode, contestId, submissionId)
	if err != nil {
		log.Printf("Failed to fetch submission: %v", err)
	}

	fmt.Println(55)

}

func (c *Client) GetStatistics(client *http.Client, groupCode string, contestId int64) []byte {
	finalData := parseAndFormEntities(c, groupCode, contestId)

	submissionCodeChan := make(map[int64]chan entities.SubmissionCodeChanObject)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	var err error
	if client == nil {
		client, err = entities.Login(c.handle, c.password)
		if err != nil {
			log.Fatalf("Login failed: %v", err)
		}
	}

	fmt.Println(finalData)

	fmt.Println(24)

	wg.Wait()

	for k, v := range submissionCodeChan {
		fmt.Println(k)
		fmt.Println(<-v)
	}

	return EntitiesToJSON(finalData)

}
