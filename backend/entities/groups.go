package entities

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type Group struct {
	GroupCode string
	GroupName string
	GroupLink string
}

// FetchGroups fetches the groups the logged-in user is part of
func FetchGroups(client *http.Client) ([]Group, error) {
	groupsURL := "https://codeforces.com/groups/my"
	req, err := http.NewRequest("GET", groupsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch groups, status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var groups []Group

	doc.Find(".groupName").Each(func(i int, s *goquery.Selection) {
		groupName := s.Text()
		groupLink, exists := s.Attr("href")
		if exists {
			groups = append(groups, Group{
				GroupCode: strings.Split(groupLink, "/")[2],
				GroupName: strings.TrimSpace(groupName),
				GroupLink: "https://codeforces.com" + groupLink,
			})
		}
	})

	return groups, nil
}
