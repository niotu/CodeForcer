package cf_api_tools

import (
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"slices"
	"testing"
)

func TestMakeTableData(t *testing.T) {
	resData := FinalJSONData{
		Problems: []entities.Problem{
			{
				Name:      "Task 1",
				Index:     "A",
				MaxPoints: 5.0,
			},
			{
				Name:      "Task 2",
				Index:     "B",
				MaxPoints: 5.0,
			},
		},
		Users: []entities.User{
			{
				Handle: "niotu",
				Solutions: map[string]*entities.Submission{
					"A": {
						Index:          "A",
						Points:         4.0,
						SubmissionId:   123456,
						ProgramLang:    "GNU C++20",
						SubmissionTime: 1720371915,
						Late:           false,
					},
					"B": {
						Index:          "B",
						Points:         2.0,
						SubmissionId:   123459,
						ProgramLang:    "GNU C++20",
						SubmissionTime: 1720371917,
						Late:           false,
					},
				},
			},
			{
				Handle: "karabas2004",
				Solutions: map[string]*entities.Submission{
					"A": {
						Index:          "A",
						Points:         5.0,
						SubmissionId:   123450,
						ProgramLang:    "Java 8",
						SubmissionTime: 1720371925,
						Late:           false,
					},
					"B": {
						Index:          "B",
						Points:         5.0,
						SubmissionId:   123459,
						ProgramLang:    "GNU C++20",
						SubmissionTime: 1720372000,
						Late:           true,
					},
				},
			},
		},
	}

	extra := ParsingParameters{
		TasksWeights:          []int{10, 10},
		ExtraHeaders:          nil,
		LatePenalty:           50,
		LateTime:              1720372124,
		SubmissionParsingMode: "best",
	}

	actual := MakeTableData(resData, extra, 8)

	expected := [][]string{
		{"User1", "r@gmatil.com", "4", "8", "2", "4", "6", "12"},
		{"User2", "f@f.ru", "5", "10", "5", "5", "10", "15"},
	}

	for i, s := range actual {
		curr := s[:len(s)-1]

		if slices.Compare(curr, expected[i]) != 0 {
			t.Errorf("MakeTableData()[%d] is expected: %v, but got: %v\n", i, expected[i], curr)
			return
		}
	}
}
