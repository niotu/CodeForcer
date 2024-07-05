package cf_api_tools

import (
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"net/http"
	"sync"
	"time"
)

var (
	ContestStatus    = "contest.status"
	ContestStandings = "contest.standings"
)

type Client struct {
	apiKey     string
	apiSecret  string
	Handle     string
	password   string
	authClient *http.Client
}

func NewClient(apiKey, apiSecret, handle, password string) (*Client, error) {
	authClient, err := entities.Login(handle, password)
	if err != nil {
		return nil, err
	}

	return &Client{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		Handle:     handle,
		password:   password,
		authClient: authClient,
	}, nil
}

func (c *Client) Authenticate() error {
	if c.authClient == nil || entities.IsCookieExpired(c.authClient) {
		client, err := entities.Login(c.Handle, c.password)
		if err != nil {
			return err
		}
		c.authClient = client
	}
	return nil
}

func (c *Client) GetGroupsList() ([]entities.Group, error) {
	var err error

	if err = c.Authenticate(); err != nil {
		return nil, err
	}

	groups, err := entities.FetchGroups(c.authClient)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (c *Client) GetContestsList(groupCode string) ([]entities.Contest, error) {
	var err error

	if err = c.Authenticate(); err != nil {
		return nil, err
	}

	contests, err := entities.FetchContests(c.authClient, groupCode)
	if err != nil {
		return nil, err
	}

	return contests, nil
}

func (c *Client) GetSubmissionCode(ch chan entities.SubmissionCodeChanObject,
	mutex *sync.Mutex, groupCode string, contestId, submissionId int64) error {
	var err error

	if err = c.Authenticate(); err != nil {
		return err
	}

	entities.FetchSubmission(c.authClient, ch, mutex, groupCode, contestId, submissionId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetStatistics(groupCode string, contestId int64, count int, weights []int) (FinalJSONData, error) {
	params := &CFContestMethodParams{
		GroupCode: groupCode,
		ContestId: contestId,
		AsManager: true,
		ApiKey:    c.apiKey,
		ApiSecret: c.apiSecret,
		Time:      time.Now().Unix(),
		Count:     count,
	}

	finalData, err := parseAndFormEntities(params, weights)
	if err != nil {
		return FinalJSONData{}, err
	}

	//submissionCodeChan := make(map[int64]chan entities.SubmissionCodeChanObject)

	//wg := &sync.WaitGroup{}
	//defer wg.Wait()
	//
	//mutex := &sync.Mutex{}

	//var err error
	//if client == nil {
	//	client, err = entities.Login(c.Handle, c.password)
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

	return *finalData, nil

}
