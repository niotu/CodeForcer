package entities

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"net/http"
	"strings"
)

type Group struct {
	GroupCode   string
	GroupName   string
	AccessLevel string
	GroupLink   string
}

var FetchGroupsFailed = errors.New("failed to fetch groups, please, try later")

// FetchGroups fetches the groups the logged-in user is part of
func FetchGroups(client *http.Client) ([]Group, error) {
	groupsURL := "https://codeforces.com/groups/my"
	req, err := http.NewRequest("GET", groupsURL, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check if response status is OK
	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Errorf("fetch groups status code: %d", resp.StatusCode))
		return nil, FetchGroupsFailed
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var groups []Group

	doc.Find(".groupName").Each(func(i int, s *goquery.Selection) {
		// Get the group name text
		groupName := s.Text()

		accessLevel := s.Parent().Next().Text()

		// Extract the href attribute from the <a> tag
		groupLink, exists := s.Attr("href")
		if exists {
			groups = append(groups, Group{
				GroupCode:   strings.Split(groupLink, "/")[2],
				GroupName:   strings.TrimSpace(groupName),
				AccessLevel: strings.TrimSpace(accessLevel),
				GroupLink:   "https://codeforces.com" + groupLink,
			})
		}
	})

	return groups, nil
}
