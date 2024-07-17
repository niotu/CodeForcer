package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	cfapitools "gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/cf-api-tools"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/solutions"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	EmptyParamsErrorMsg  = "Some parameters are empty"
	UserNotFoundErrorMsg = "User not found"
)

var NoFileProvided = errors.New("no file provided")

type EntitiesResponseObject interface{}

func statusFailedResponse(comment string) []byte {
	resp := struct {
		Status  string `json:"status"`
		Comment string `json:"comment"`
	}{Status: "FAILED", Comment: comment}

	jsonResp, _ := json.Marshal(&resp)

	return jsonResp
}

func statusOKResponse(obj EntitiesResponseObject) []byte {
	resp := struct {
		Status                 string `json:"status"`
		EntitiesResponseObject `json:"result"`
	}{Status: "OK", EntitiesResponseObject: obj}

	jsonResp, _ := json.Marshal(&resp)

	return jsonResp
}

func isEmptyParams(params ...string) bool {
	for _, p := range params {
		if p == "" {
			return true
		}
	}
	return false
}

func validateAndWrite(w http.ResponseWriter, client *cfapitools.Client, params ...string) bool {
	if isEmptyParams(params...) {
		_, _ = w.Write(statusFailedResponse(EmptyParamsErrorMsg))
		return false
	}

	if client == nil {
		_, _ = w.Write(statusFailedResponse(UserNotFoundErrorMsg))
		return false
	}

	return true
}

func parseWeights(weightsString string) ([]int, error) {
	weightsSplitted := strings.Split(weightsString, "-")
	var weights []int
	for _, s := range weightsSplitted {
		weight, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		weights = append(weights, weight)
	}

	return weights, nil
}

func createMultipart(w http.ResponseWriter, jsonData []byte, userId string) {
	AttachZipError := "unable to attach zip file"
	resultZipName := solutions.GetResultZipName(userId)

	defer func() {
		go func() {
			_ = os.Remove(resultZipName)
		}()
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	w.Header().Set("Content-Type", "multipart/mixed; boundary="+writer.Boundary())

	// Write JSON part
	jsonPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {"application/json"},
	})
	if err != nil {
		logger.Error(err)
		_, _ = w.Write(statusFailedResponse(AttachZipError))
		return
	}
	if _, err := jsonPart.Write(jsonData); err != nil {
		logger.Error(err)
		_, _ = w.Write(statusFailedResponse(AttachZipError))
		return
	}

	// Write ZIP part

	zipBytes, err := os.ReadFile(resultZipName)
	if err != nil {
		logger.Error(err)
		_, _ = w.Write(statusFailedResponse(AttachZipError))
		return
	}

	zipPart, err := writer.CreatePart(map[string][]string{
		"Content-Type":        {"application/zip"},
		"Content-Disposition": {fmt.Sprintf("attachment; filename=\"%s\"", resultZipName)},
	})
	if err != nil {
		logger.Error(err)
		_, _ = w.Write(statusFailedResponse(AttachZipError))
		return
	}
	if _, err := zipPart.Write(zipBytes); err != nil {
		logger.Error(err)
		_, _ = w.Write(statusFailedResponse(AttachZipError))
		return
	}

	writer.Close()

	_, _ = w.Write(body.Bytes())

}

func getZipFile(r *http.Request, srcZip string) error {
	err := r.ParseMultipartForm(15 << 20)
	if err != nil {
		logger.Logger().Info("NO FILE PROVIDED. IGNORE THIS STEP")
		return NoFileProvided
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Logger().Info("NO FILE PROVIDED. IGNORE THIS STEP")
		return NoFileProvided
	}
	defer file.Close()

	logger.Logger().Info("File downloaded successfully",
		zap.Int64("size", handler.Size),
		zap.String("name", handler.Filename))

	tempFile, err := os.Create(srcZip)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("could not create temp file")
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return fmt.Errorf("could not save uploaded file")
	}

	return nil
}
