package backend_test

import (
	cfapitools "gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/cf-api-tools"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"testing"
)

func TestAuthorization(t *testing.T) {
	logger.Init()
	_, err := cfapitools.NewClientWithAuth(
		"33bcdcdcf956dc632a5aa98fa94697a1bb06406c",
		"4c0942196312bbf66eb019fd4f2dfec6534d8c1w",
		"rammao",
		"rammaio",
	)

	if err == nil {
		t.Errorf("Should produce an error since authorization data is incorrect")
	}

}
