package db

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var usersFilePath = "./db/users.json"

func UploadUsersToFile(table []byte) error {
	reader := bytes.NewReader(table)

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

func GetUsers() map[string]string {
	file, err := os.Open(usersFilePath)
	if err != nil {
		file, _ = os.OpenFile(usersFilePath, os.O_CREATE|os.O_RDONLY|os.O_TRUNC, 0606)

		data, _ := json.Marshal(map[string]string{})
		_, _ = file.Write(data)

		file.Close()
		return make(map[string]string)
	}
	defer file.Close()

	buff, err := io.ReadAll(file)
	if err != nil {
		logger.Error(err)
		return nil
	}

	var data map[string]string

	err = json.Unmarshal(buff, &data)
	if err != nil {
		logger.Error(err)
		return nil
	}

	return data
}
