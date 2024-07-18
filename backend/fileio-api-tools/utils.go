package fileio_api_tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func createMultipartBody(filePath string) (*multipart.Writer, bytes.Buffer) {
	file, _ := os.Open(filePath)
	defer file.Close()

	// Create a buffer to hold the multipart form data
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	// Add the file
	fileWriter, _ := multipartWriter.CreateFormFile("file", "submissions.zip")

	_, _ = io.Copy(fileWriter, file)

	// Close the multipart writer to set the terminating boundary
	multipartWriter.Close()

	return multipartWriter, requestBody
}

func StoreFile(filePath string) (string, error) {
	zipError := fmt.Errorf("error storing zip file")

	mw, body := createMultipartBody(filePath)

	req, _ := http.NewRequest("POST", "https://file.io/?expires=1w&maxDownloads=1&autoDelete=true", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Errorf("error making request to %s: %w", req.URL.String(), err))
		return "", zipError
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(err)
		return "", zipError
	}

	var jsonResponse map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&jsonResponse)

	return jsonResponse["link"].(string), nil
}
