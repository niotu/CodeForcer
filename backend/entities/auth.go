package entities

import (
	"errors"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var LoginFailed = errors.New("login failed: try later")

// generateCSRFToken fetches the CSRF token from the login page
func generateCSRFToken(body []byte) (string, error) {
	reg := regexp.MustCompile(`csrf='(.+?)'`)
	matches := reg.FindSubmatch(body)
	if len(matches) < 2 {
		logger.Error(errors.New("cannot find CSRF token"))
		return "", errors.New("authorization cannot be proceeded, please, try later")
	}
	return string(matches[1]), nil
}

// generateFtaa generates a random ftaa value
func generateFtaa() string {
	return RandString(18)
}

// generateBfaa generates a static bfaa value
func generateBfaa() string {
	return "f1b3f18c715565b589b7823cda7448ce"
}

// RandString generates a random string of specified length
func RandString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// findHandle extracts the handle from the response body
func findHandle(body []byte) (string, error) {
	reg := regexp.MustCompile(`handle = "([\s\S]+?)"`)

	matches := reg.FindSubmatch(body)
	if len(matches) < 2 {
		return "", errors.New("login failed: incorrect handle or password, try later")
	}
	return string(matches[1]), nil
}

// Login logs into Codeforces and returns an authenticated HTTP client
func Login(handle, password string) (*http.Client, error) {
	logger.Logger().Info("Login process...",
		zap.String("Handle", handle))

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginURL := "https://codeforces.com/enter"
	resp, err := client.Get(loginURL)
	if err != nil {
		logger.Error(err)
		return nil, LoginFailed
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, LoginFailed
	}

	csrfToken, err := generateCSRFToken(body)
	if err != nil {
		return nil, err
	}

	ftaa := generateFtaa()
	bfaa := generateBfaa()

	form := url.Values{}
	form.Add("csrf_token", csrfToken)
	form.Add("action", "enter")
	form.Add("ftaa", ftaa)
	form.Add("bfaa", bfaa)
	form.Add("handleOrEmail", handle)
	form.Add("password", password)
	form.Add("_tta", "176")
	form.Add("remember", "on")

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(form.Encode()))
	if err != nil {
		logger.Error(err)
		return nil, LoginFailed
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", loginURL)

	loginResp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, LoginFailed
	}
	defer loginResp.Body.Close()

	loginBody, err := io.ReadAll(loginResp.Body)
	if err != nil {
		logger.Error(err)
		return nil, LoginFailed
	}

	_, err = findHandle(loginBody)
	if err != nil {
		logger.Logger().Error(err.Error(), zap.String("handle", handle))
		return nil, err
	}

	logger.Logger().Info("Login success: ", zap.String("handle", handle))
	return client, nil
}

func IsCookieExpired(client *http.Client) bool {
	jar := client.Jar
	if jar == nil {
		return true
	}

	// Get all cookies for the codeforces domain
	u, _ := url.Parse("https://codeforces.com")
	cookies := jar.Cookies(u)

	// Check if any cookie is expired
	for _, cookie := range cookies {
		// If a cookie has an expiration date and it's before now, consider it expired
		if !cookie.Expires.IsZero() && cookie.Expires.Before(time.Now()) {
			return true
		}
	}
	return false
}
