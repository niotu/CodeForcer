package cf_api_tools

type CFContestMethodParams struct {
	GroupCode string `url:"groupCode"`
	ContestId int64  `url:"contestId"`
	AsManager bool   `url:"asManager"`
	ApiKey    string `url:"apiKey"`
	ApiSecret string `url:"-"`
	Time      int64  `url:"time"`
	Count     int    `url:"count,omitempty"`
}

func (c *CFContestMethodParams) GetKey() string {
	return c.ApiKey
}

func (c *CFContestMethodParams) GetSecret() string {
	return c.ApiSecret
}
