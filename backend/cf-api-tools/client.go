package cf_api_tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/solutions"
	"net/http"
	"time"
)

var (
	ContestStatus    = "contest.status"
	ContestStandings = "contest.standings"
	key              = []byte("X0UcXcvv3hMj9XQpf584VlHxpWnoK7Vx")
)

type Client struct {
	apiKey      []byte
	apiSecret   []byte
	Handle      string `json:"-"`
	password    string
	authClient  *http.Client
	currContest *entities.Contest
}

func encryptClientData(plaintext string) []byte {
	aes, _ := aes.NewCipher(key)

	gcm, _ := cipher.NewGCM(aes)

	nonce := make([]byte, gcm.NonceSize())
	_, _ = rand.Read(nonce)

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return ciphertext
}

func decryptClientData(ciphertext []byte) string {
	aes, _ := aes.NewCipher(key)

	gcm, _ := cipher.NewGCM(aes)

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)

	return string(plaintext)
}

func NewClient(apiKey, apiSecret string) (*Client, error) {
	return &Client{
		apiKey:    encryptClientData(apiKey),
		apiSecret: encryptClientData(apiSecret),
	}, nil
}

func NewClientWithAuth(apiKey, apiSecret, handle, password string) (*Client, error) {
	authClient, err := entities.Login(handle, password)
	if err != nil {
		return nil, err
	}

	return &Client{
		apiKey:     encryptClientData(apiKey),
		apiSecret:  encryptClientData(apiSecret),
		Handle:     handle,
		password:   password,
		authClient: authClient,
	}, nil
}

func (c *Client) DecodeApiKey() string {
	return decryptClientData(c.apiKey)
}

func (c *Client) DecodeApiSecret() string {
	return decryptClientData(c.apiSecret)
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

func (c *Client) GetContestData(groupCode string, contestId int64) (*DataFromStandings, error) {
	params := &CFContestMethodParams{
		GroupCode: groupCode,
		ContestId: contestId,
		AsManager: true,
		ApiKey:    c.DecodeApiKey(),
		ApiSecret: c.DecodeApiSecret(),
		Time:      time.Now().Unix(),
		Count:     1,
	}

	data, err := formattedStandings(params)
	if err != nil {
		return nil, err
	}

	c.currContest = &entities.Contest{
		Id:               contestId,
		Name:             data.Name,
		GroupCode:        groupCode,
		DurationSeconds:  data.DurationSeconds,
		StartTimeSeconds: data.StartTimeSeconds,
	}

	return data, nil
}

func (c *Client) GetStatistics(groupCode string, contestId int64, count int, tableParams ParsingParameters) (FinalJSONData, error) {
	params := &CFContestMethodParams{
		GroupCode: groupCode,
		ContestId: contestId,
		AsManager: true,
		ApiKey:    c.DecodeApiKey(),
		ApiSecret: c.DecodeApiSecret(),
		Time:      time.Now().Unix(),
		Count:     count,
	}

	finalData, err := combineStatusAndStandings(params, tableParams)
	if err != nil {
		return FinalJSONData{}, err
	}

	return *finalData, nil

}

func GetSolutions(srcArchive, userId string, finalData FinalJSONData) error {
	submissonCodes := make(map[int64]entities.User)

	for _, u := range finalData.Users {
		for _, s := range u.Solutions {
			if s.SubmissionId != -1 {
				submissonCodes[s.SubmissionId] = u
			}
		}
	}

	err := solutions.MakeSolutionsArchive(srcArchive, userId, submissonCodes)
	if err != nil {
		return err
	}
	return nil
}
