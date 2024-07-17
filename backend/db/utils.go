package db

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	//cfapitools "gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/cf-api-tools"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var usersFilePath = "./db/users.json"
var clientsFilePath = "./db/clients.json"

func UploadUsersToFile(csvTable []byte) error {
	reader := bytes.NewReader(csvTable)

	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	records, err := csvReader.ReadAll()
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("unable to upload csv file")
	}

	data := GetUsers()

	headers := records[0]
	emailIdx, handleIdx := 0, 1
	if headers[1] == "email" {
		emailIdx, handleIdx = 1, 0
	}

	records = records[1:]

	for _, row := range records {
		first, second := row[0], row[1]

		if !strings.Contains(first, "@") && !strings.Contains(second, "@") {
			logger.Error(err)
			return fmt.Errorf("csv format do not match the needed one")
		}

		data[row[handleIdx]] = row[emailIdx]
	}

	buff, err := json.Marshal(data)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("unable to upload csv file")
	}

	path, _ := filepath.Abs(usersFilePath)

	err = os.WriteFile(path, buff, 0606)

	if err != nil {
		logger.Error(err)
		return fmt.Errorf("unable to upload csv file")
	}

	return nil
}

func UploadClientsToFile(clients map[string]interface{}) {
	buff, err := json.Marshal(clients)
	if err != nil {
		logger.Error(err)
	}

	path, _ := filepath.Abs(clientsFilePath)

	err = os.WriteFile(path, buff, 0606)
	if err != nil {
		logger.Error(err)
	}
}

func GetUsers() map[string]string {
	data := getDB(usersFilePath)

	result := make(map[string]string)

	for k, v := range data {
		result[k] = v.(string)
	}

	return result
}

func GetClientsBytes() []byte {
	b, _ := getDBBytes(clientsFilePath)
	return b
}

func getDB(filePath string) map[string]interface{} {
	buff, err := getDBBytes(filePath)
	if err != nil {
		return nil
	}
	if len(buff) == 0 {
		return make(map[string]interface{})
	}

	var data map[string]interface{}

	err = json.Unmarshal(buff, &data)
	if err != nil {
		logger.Error(err)
		return nil
	}

	return data
}

func getDBBytes(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		_ = os.MkdirAll(filepath.Dir(filePath), 0606)
		file, _ = os.OpenFile(filePath, os.O_CREATE|os.O_RDONLY|os.O_TRUNC, 0606)

		data, _ := json.Marshal(map[string]string{})
		_, _ = file.Write(data)

		file.Close()
		return []byte{}, nil
	}
	defer file.Close()

	buff, err := io.ReadAll(file)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return buff, nil
}
