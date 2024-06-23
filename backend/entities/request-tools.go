package entities

import (
	"math/rand"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
)

type ApiSignature struct {
	Rand   string      `url:"rand"`
	Method string      `url:"-"`
	Params *url.Values `url:"-"`
}

type ContestStatusRequestParams struct {
	GroupCode string `url:"groupCode"`
	ContestId int    `url:"contestId"`
	AsManager bool   `url:"asManager"`
	ApiKey    string `url:"apiKey"`
	ApiSecret string `url:"-"`
	Time      int64  `url:"time"`
	Count     int    `url:"count,omitempty"`
}

func (c *ContestStatusRequestParams) GetKey() string {
	return c.ApiKey
}

func (c *ContestStatusRequestParams) GetSecret() string {
	return c.ApiSecret
}

func (c *ContestStatusRequestParams) UpdateTime() {
	c.Time = time.Now().Unix()
}

type ContestStandingsRequestParams struct {
	GroupCode string `url:"groupCode"`
	ContestId int    `url:"contestId"`
	AsManager bool   `url:"asManager"`
	ApiKey    string `url:"apiKey"`
	ApiSecret string `url:"-"`
	Time      int64  `url:"time"`
	Count     int    `url:"count,omitempty"`
}

func (c *ContestStandingsRequestParams) GetKey() string {
	return c.ApiKey
}

func (c *ContestStandingsRequestParams) GetSecret() string {
	return c.ApiSecret
}

func NewApiSignature() *ApiSignature {
	randNum := strconv.Itoa(rand.Intn(999999))

	return &ApiSignature{
		Rand: strings.Repeat("0", 6-len(randNum)) + randNum,
	}
}

func (as *ApiSignature) SortApiSigParams() string {
	paramsSlice := strings.Split(as.Params.Encode(), "&")
	slices.Sort(paramsSlice)

	return strings.Join(paramsSlice, "&")
}
