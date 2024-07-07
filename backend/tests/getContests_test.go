package backend

import (
	cfapitools "gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/cf-api-tools"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"testing"
)

func TestGetContests(t *testing.T) {
	logger.Init()

	client, err := cfapitools.NewClientWithAuth(
		"33bcdcdcf956dc632a5aa98fa94697a1bb06406c",
		"4c0942196312bbf66eb019fd4f2dfec6534d8c1w",
		"ramzestwo",
		"ramazan1810",
	)
	if err != nil {
		t.Fatalf("Should not produce an error since authorization data is correct: %v", err)
	}

	contests, err := client.GetContestsList("bfRCcT6pgf")
	if err != nil {
		t.Fatalf("Should not produce an error since client is authorized and groupCode is correct: %v", err)
	}

	expectedContest := entities.Contest{
		Id:               530794,
		Name:             "contest1",
		GroupCode:        "bfRCcT6pgf",
		ContestLink:      "https://codeforces.com/group/bfRCcT6pgf/contest/530794",
		DurationString:   "7:00:00",
		StartTimeString:  "Jun/18/2024 13:25",
		DurationSeconds:  604800,
		StartTimeSeconds: 1718717100,
		Problems:         []entities.Problem{},
	}

	for _, c := range contests {
		if c.Id != expectedContest.Id {
			t.Errorf("Expected Id %d, but got %d", expectedContest.Id, c.Id)
		}
		if c.Name != expectedContest.Name {
			t.Errorf("Expected Name '%s', but got '%s'", expectedContest.Name, c.Name)
		}
		if c.GroupCode != expectedContest.GroupCode {
			t.Errorf("Expected GroupCode '%s', but got '%s'", expectedContest.GroupCode, c.GroupCode)
		}
		if c.ContestLink != expectedContest.ContestLink {
			t.Errorf("Expected ContestLink '%s', but got '%s'", expectedContest.ContestLink, c.ContestLink)
		}
		if c.DurationString != expectedContest.DurationString {
			t.Errorf("Expected DurationString '%s', but got '%s'", expectedContest.DurationString, c.DurationString)
		}
		if c.StartTimeString != expectedContest.StartTimeString {
			t.Errorf("Expected StartTimeString '%s', but got '%s'", expectedContest.StartTimeString, c.StartTimeString)
		}
		if c.DurationSeconds != expectedContest.DurationSeconds {
			t.Errorf("Expected DurationSeconds %d, but got %d", expectedContest.DurationSeconds, c.DurationSeconds)
		}
		if c.StartTimeSeconds != expectedContest.StartTimeSeconds {
			t.Errorf("Expected StartTimeSeconds %d, but got %d", expectedContest.StartTimeSeconds, c.StartTimeSeconds)
		}
		if !compareProblems(c.Problems, expectedContest.Problems) {
			t.Errorf("Expected Problems '%v', but got '%v'", expectedContest.Problems, c.Problems)
		}
	}
}

// Helper function to compare slices of problems
func compareProblems(a, b []entities.Problem) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
