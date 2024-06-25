package cf_api_tools

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"net/url"
)

type RequestParams interface {
	GetKey() string
	GetSecret() string
}

type ApiRequest struct {
	Method string
	Params RequestParams
	ApiSig *ApiSignature
}

func NewApiRequest(method string, params RequestParams) *ApiRequest {

	return &ApiRequest{
		Method: method,
		Params: params,
		ApiSig: NewApiSignature(),
	}
}

func (a *ApiRequest) GetUrlParams() url.Values {
	vals, _ := query.Values(a.Params)
	return vals
}

func (a *ApiRequest) GetApiSigHash() string {
	u, _ := url.Parse(fmt.Sprint(a.ApiSig.Rand) + "/" + a.Method)

	args := a.GetUrlParams()

	u.RawQuery = args.Encode()

	res := u.String() + "#" + fmt.Sprint(a.Params.GetSecret())

	sha := sha512.New()
	sha.Write([]byte(res))

	fmt.Printf("%x\n", sha.Sum(nil))

	return fmt.Sprintf("%x", sha.Sum(nil))
}

func (a *ApiRequest) GetApiSig() string {
	return fmt.Sprint(a.ApiSig.Rand) + a.GetApiSigHash()
}

func (a *ApiRequest) MakeApiRequest() ([]byte, error) {
	if a.Method != CONTEST_STATUS && a.Method != CONTEST_STANDINGS {
		return nil, errors.New("no such api request method")
	}

	u, _ := url.Parse("https://codeforces.com/api/" + a.Method)

	params := a.GetUrlParams()

	params.Add("apiSig", a.GetApiSig())

	u.RawQuery = params.Encode()

	fmt.Println(u.String())

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
