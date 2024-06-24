package entities

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
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

func (c *Client) GetContestStatus(con Contest) interface{} {
	params := &ContestStatusRequestParams{
		GroupCode: con.GroupCode,
		ContestId: con.Id,
		AsManager: true,
		ApiKey:    c.apiKey,
		ApiSecret: c.apiSecret,
		Time:      time.Now().Unix(),
	}
	api := NewApiRequest(CONTEST_STATUS, params)
	fmt.Println(api.ApiSig.Rand, params.Time)

	resp, err := api.MakeApiRequest()
	if err != nil {
		fmt.Println(err)
	}

	var data interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Println(err)
	}

	return data
}

func (c *Client) GetContestStandings(con Contest) interface{} {
	params := &ContestStandingsRequestParams{
		GroupCode: con.GroupCode,
		ContestId: con.Id,
		AsManager: true,
		ApiKey:    c.apiKey,
		ApiSecret: c.apiSecret,
		Time:      time.Now().Unix(),
		Count:     1,
	}
	api := NewApiRequest(CONTEST_STANDINGS, params)

	resp, err := api.MakeApiRequest()
	if err != nil {
		log.Println(err)
	}

	var data interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Println(err)
	}

	return data
}

func ParseContestStatus(data interface{}, dataStatus *DataFromStatus, dataStandings *DataFromStandings) (*DataFromStatus, error) {
	db := make(map[string]*User)

	if resp := data.(map[string]interface{}); resp["status"].(string) == "FAILED" {
		comment := resp["comment"].(string)

		errorMsg := ""

		if strings.Contains(comment, "asManager") {
			errorMsg = "You are not the manager of contest or group. Be sure that on the page of contest you selected the following:\n" +
				"Administration (block on the right) -> Enable manager mode."
		}

		return nil, errors.New(errorMsg)
	}

	result := data.(map[string]interface{})["result"].([]interface{})

	if dataStatus == nil {
		dataStatus = &DataFromStatus{}
	}

	for _, resultElem := range result {
		submissionJson := resultElem.(map[string]interface{})

		username := submissionJson["author"].(map[string]interface{})["members"].([]interface{})[0].(map[string]interface{})["handle"].(string)
		problemIdx := submissionJson["problem"].(map[string]interface{})["index"].(string)

		//go fetchSubmission(nil)

		if _, ok := db[username]; !ok {
			submissions := make(map[string]*Submission)
			for _, p := range dataStandings.Problems {
				submissions[p.Index] = &Submission{
					Index:          p.Index,
					Solution:       "",
					Points:         0,
					SubmissionId:   -1,
					ProgramLang:    "",
					SubmissionTime: 0,
				}
			}

			db[username] = &User{
				Handle:    username,
				Solutions: submissions,
			}
		}
		if s, _ := db[username].Solutions[problemIdx]; s.SubmissionId == -1 {
			_, exists := submissionJson["points"]
			submissionPoints := 0.0
			if exists {
				submissionPoints = submissionJson["points"].(float64)
			}

			problemVerdict := submissionJson["verdict"].(string)
			if value := dataStatus.ProblemMaxPoints[problemIdx]; problemVerdict == "OK" && value == 0.0 {
				dataStatus.ProblemMaxPoints[problemIdx] = submissionPoints
			}

			submission := &Submission{
				Index:          submissionJson["problem"].(map[string]interface{})["index"].(string),
				Solution:       "",
				Points:         submissionPoints,
				SubmissionId:   int64(submissionJson["id"].(float64)),
				ProgramLang:    submissionJson["programmingLanguage"].(string),
				SubmissionTime: int64(submissionJson["creationTimeSeconds"].(float64)),
			}

			db[username].Solutions[problemIdx] = submission
		}

	}

	// test
	//for k, v := range db {
	//	fmt.Println(k)
	//	for _, vv := range v.Solutions {
	//		fmt.Println(*vv)
	//	}
	//}

	dataStatus.Users = db

	return dataStatus, nil
}

type DataFromStandings struct {
	Problems         []*Problem
	DurationSeconds  int64
	StartTimeSeconds int64
}

type DataFromStatus struct {
	ProblemMaxPoints map[string]float64
	Users            map[string]*User
}

type FinalJSONData struct {
	Problems []Problem `json:"problems"`
	Users    []User    `json:"users"`
	CSV      []byte    `json:"csv"`
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

func ParseContestStandings(data interface{}, dataStandings *DataFromStandings) (*DataFromStandings, error) {
	var problems []*Problem

	if resp := data.(map[string]interface{}); resp["status"].(string) == "FAILED" {
		comment := resp["comment"].(string)

		errorMsg := ""

		if strings.Contains(comment, "asManager") {
			errorMsg = "You are not the manager of contest or group. Be sure that on the page of contest you selected the following:\n" +
				"Administration (block on the right) -> Enable manager mode."
		}

		return nil, errors.New(errorMsg)
	}

	result := data.(map[string]interface{})["result"].(map[string]interface{})
	problemsInfo := result["problems"].([]interface{})

	for _, problemJson := range problemsInfo {
		problems = append(problems, &Problem{
			Name:  problemJson.(map[string]interface{})["name"].(string),
			Index: problemJson.(map[string]interface{})["index"].(string),
		})
	}

	contest := result["contest"].(map[string]interface{})
	durationSeconds := contest["durationSeconds"].(float64)
	startTimeSeconds := contest["startTimeSeconds"].(float64)

	// test
	fmt.Println(problems)

	return &DataFromStandings{
		Problems:         problems,
		DurationSeconds:  int64(durationSeconds),
		StartTimeSeconds: int64(startTimeSeconds),
	}, nil
}

func (c *Client) ParseAndFormEntities(con Contest) []byte {
	standings := c.GetContestStandings(con)
	dataStandings, err := ParseContestStandings(standings, nil)
	if err != nil {
		log.Println(err)
	}

	dataStatus := &DataFromStatus{
		ProblemMaxPoints: make(map[string]float64),
		Users:            nil,
	}

	for _, problem := range dataStandings.Problems {
		dataStatus.ProblemMaxPoints[problem.Index] = 0.0
	}

	status := c.GetContestStatus(con)
	dataStatus, err = ParseContestStatus(status, dataStatus, dataStandings)

	for _, problem := range dataStandings.Problems {
		problem.MaxPoints = dataStatus.ProblemMaxPoints[problem.Index]
	}

	//ProblemListToJSON(dataStandings.Problems)

	var u []User
	for _, v := range dataStatus.Users {
		u = append(u, *v)
	}

	var p []Problem
	for _, v := range dataStandings.Problems {
		p = append(p, *v)
	}

	finalJsonData := FinalJSONData{
		Problems: p,
		Users:    u,
	}

	finalJsonData.CSV = MakeCSVFile(finalJsonData).Bytes()

	//UserListToJSON(u)

	return EntitiesToJSON(finalJsonData)
}

//func EntitiesToJSON[T any](jsonData T, filename string) []byte {
//	data, err := json.Marshal(jsonData)
//	if err != nil {
//		log.Println(err)
//	}
//
//	file, _ := os.OpenFile(filename+".json", os.O_CREATE|os.O_TRUNC, 0606)
//	file.Write(data)
//
//	return data
//}

func EntitiesToJSON[T any](jsonData T) []byte {
	data, err := json.Marshal(jsonData)
	if err != nil {
		log.Println(err)
	}

	return data
}

func (c *Client) GetGroupsList(client *http.Client) ([]Group, *http.Client) {
	var err error

	if client == nil {
		client, err = Login(c.handle, c.password)
		if err != nil {
			log.Printf("Login failed: %v", err)
		}
	}

	groups, err := fetchGroups(client)
	if err != nil {
		log.Printf("Failed to fetch groups: %v", err)
	}

	//for _, group := range groups {
	//	fmt.Printf("Group: %s, Link: %s\n", group.GroupName, group.GroupCode)
	//}

	return groups, client
}

func (c *Client) GetContestsList(client *http.Client, groupCode string) ([]Contest, *http.Client) {
	var err error

	if client == nil {
		client, err = Login(c.handle, c.password)
		if err != nil {
			log.Printf("Login failed: %v", err)
		}
	}

	contests, err := fetchContests(client, groupCode)
	if err != nil {
		log.Printf("Failed to fetch groups: %v", err)
	}

	//for _, contest := range contests {
	//	fmt.Println(contest)
	//}

	return contests, client
}

func (c *Client) GetSubmissionCode(client *http.Client, groupCode string, contestId, submissionId int) (string, *http.Client) {
	var err error

	if client == nil {
		client, err = Login(c.handle, c.password)
		if err != nil {
			log.Printf("Login failed: %v", err)
		}
	}

	submission, err := fetchSubmission(client, groupCode, contestId, submissionId)
	if err != nil {
		log.Printf("Failed to fetch submission: %v", err)
	}

	fmt.Println(submission)

	//for _, contest := range contests {
	//	fmt.Println(contest)
	//}

	return submission, client
}

//func main() {
//	api := NewClient()
//	////api.SetApiKey("002d5e9812d982b6e8b353daf3a866cdc3cb012b")
//	////api.SetApiSecret("67b10e4a10dc0379df8bc1775afc9495d91b8055")
//	api.SetApiKey("72bcdcdcf956dc632a5aa98fa94697a1bb06406c")
//	api.SetApiSecret("4c0942196312bbf66eb019fd4f2dfec6534d8c1b")
//	//
//	////con := Contest{
//	////	Id:               530794,
//	////	GroupCode:        "bfRCcT6pgf",
//	////	DurationSeconds:  0,
//	////	StartTimeSeconds: 0,
//	////	Problems:         nil,
//	////}
//	//con := Contest{
//	//	Id:               504401,
//	//	GroupCode:        "CsTlwuSxCL",
//	//	DurationSeconds:  0,
//	//	StartTimeSeconds: 0,
//	//	Problems:         nil,
//	//}
//	//api.GetContestStatus(con)
//
//	g, client := GetGroupsList(nil)
//	//GroupListToJSON(g)
//	fmt.Println(254)
//
//	c, client := GetContestsList(client, g[1])
//	ContestListToJSON(c)
//
//	api.ParseAndFormEntities(c[1])
//}
