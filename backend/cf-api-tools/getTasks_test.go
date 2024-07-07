package cf_api_tools

import (
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"testing"
)

func TestGetTasksFailure(t *testing.T) {
	logger.Init()

	client, err := NewClient(
		"33bcdcdcf956dc632a5aa98fa94697a1bb06406c",
		"4c0942196312bbf66eb019fd4f2dfec6534d8c1w",
	)
	if err != nil {
		t.Errorf("Should not produce an error since authorization data is correct")
	}

	_, err = client.GetContestData("CsTlwuSxCL", 504401)
	if err == nil {
		t.Errorf("Should produce an error since apiKey is not correct")
		return
	}

}

func TestGetTasksSuccess(t *testing.T) {
	logger.Init()

	client, err := NewClient(
		"72bcdcdcf956dc632a5aa98fa94697a1bb06406c",
		"4c0942196312bbf66eb019fd4f2dfec6534d8c1b",
	)
	if err != nil {
		t.Errorf("Should not produce an error since authorization data is correct")
	}

	data, err := client.GetContestData("CsTlwuSxCL", 504401)
	if err != nil {
		t.Errorf("Should not produce an error since authorization, groupCode and contestId are correct")
		return
	}

	tasks := data.Problems

	expected := []entities.Problem{
		{
			Name:      "Task 1",
			Index:     "A",
			MaxPoints: 0,
		},
		{
			Name:      "Task 2",
			Index:     "B",
			MaxPoints: 0,
		},
	}

	for i, tsk := range tasks {
		if tsk.Name != expected[i].Name {
			t.Errorf("Expected task name '%s', but got '%s'", expected[i].Name, tsk.Name)
		}
		if tsk.Index != expected[i].Index {
			t.Errorf("Expected task index '%s', but got '%s'", expected[i].Index, tsk.Index)
		}
		if tsk.MaxPoints != expected[i].MaxPoints {
			t.Errorf("Expected task max points '%f', but got '%f'", expected[i].MaxPoints, tsk.MaxPoints)
		}
	}
}
