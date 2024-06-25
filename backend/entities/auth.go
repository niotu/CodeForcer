package entities

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

// generateCSRFToken fetches the CSRF token from the login page
func generateCSRFToken(body []byte) (string, error) {
	reg := regexp.MustCompile(`csrf='(.+?)'`)
	matches := reg.FindSubmatch(body)
	if len(matches) < 2 {
		return "", errors.New("Cannot find CSRF token")
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

	//fmt.Println("auth.go/findHandle: ", string(body))

	matches := reg.FindSubmatch(body)
	if len(matches) < 2 {
		return "", errors.New("Not logged in")
	}
	return string(matches[1]), nil
}

// Login logs into Codeforces and returns an authenticated HTTP client
func Login(handle, password string) (*http.Client, error) {
	fmt.Println("login process:\n")
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginURL := "https://codeforces.com/enter"
	resp, err := client.Get(loginURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", loginURL)

	loginResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer loginResp.Body.Close()

	loginBody, err := io.ReadAll(loginResp.Body)
	if err != nil {
		return nil, err
	}

	handle, err = findHandle(loginBody)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Login successful. Welcome, %s!\n", handle)
	return client, nil
}
