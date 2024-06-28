package cf_api_tools

import (
	"math/rand"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

type ApiSignature struct {
	Rand   string      `url:"rand"`
	Method string      `url:"-"`
	Params *url.Values `url:"-"`
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
