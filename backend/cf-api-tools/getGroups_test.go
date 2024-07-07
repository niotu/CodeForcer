package cf_api_tools

import (
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"testing"
)

func TestGetGroups(t *testing.T) {
	logger.Init()

	client, err := NewClientWithAuth(
		"33bcdcdcf956dc632a5aa98fa94697a1bb06406c",
		"4c0942196312bbf66eb019fd4f2dfec6534d8c1w",
		"ramzestwo",
		"ramazan1810",
	)
	if err != nil {
		t.Errorf("Should not produce an error since authorization data is correct")
	}

	groups, err := client.GetGroupsList()
	if err != nil {
		t.Errorf("Should not produce an error since client is authorized")
	}

	for i, g := range groups {
		if g.AccessLevel == "Manager" {
			switch i {
			case 0:
				if g.GroupCode != "bfRCcT6pgf" {
					t.Errorf("Expected GroupCode 'bfRCcT6pgf', but got '%s'", g.GroupCode)
				}
				if g.GroupName != "test_swp" {
					t.Errorf("Expected GroupName 'test_swp', but got '%s'", g.GroupName)
				}
				if g.GroupLink != "https://codeforces.com/group/bfRCcT6pgf" {
					t.Errorf("Expected GroupLink 'https://codeforces.com/group/bfRCcT6pgf', but got '%s'", g.GroupLink)
				}
			case 1:
				if g.GroupCode != "CsTlwuSxCL" {
					t.Errorf("Expected GroupCode 'CsTlwuSxCL', but got '%s'", g.GroupCode)
				}
				if g.GroupName != "IU TCS Spring 2024" {
					t.Errorf("Expected GroupName 'IU TCS Spring 2024', but got '%s'", g.GroupName)
				}
				if g.GroupLink != "https://codeforces.com/group/CsTlwuSxCL" {
					t.Errorf("Expected GroupLink 'https://codeforces.com/group/CsTlwuSxCL', but got '%s'", g.GroupLink)
				}
			}
		}
	}
}
