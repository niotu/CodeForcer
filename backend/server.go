package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	mx "github.com/gorilla/mux"
	"github.com/joho/godotenv"
	cfapitools "gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/cf-api-tools"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/db"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var clients sync.Map // Concurrent map to store clients
var clientKeyToId sync.Map

func getClient(userID string) *cfapitools.Client {
	client, ok := clients.Load(userID)
	if !ok {
		return nil
	}
	return client.(*cfapitools.Client)
}

func getIdByClient(client *cfapitools.Client) string {
	key := client.DecodeApiKey()
	fmt.Println(key)

	id, ok := clientKeyToId.Load(key)
	if !ok {
		return ""
	}
	return id.(string)
}

func setClient(userID string, client *cfapitools.Client) {
	key := client.DecodeApiKey()
	clientKeyToId.Store(key, userID)
	clients.Store(userID, client)
}

func panicLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Logger().Error("",
					zap.Any("code", err))
				_, _ = w.Write(statusFailedResponse("server error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func infoLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger().Info("Request:",
			zap.String("Method", r.Method),
			zap.String("Path", r.URL.Path),
			zap.String("Query", r.URL.Query().Encode()))
		next.ServeHTTP(w, r)
	})
}

func setAdminData(w http.ResponseWriter, r *http.Request) {
	apiKey := r.URL.Query().Get("key")
	apiSecret := r.URL.Query().Get("secret")
	//handle := r.URL.Query().Get("handle")
	//password := r.URL.Query().Get("password")
	//userId := "123"
	userId := uuid.New().String()

	if isEmptyParams(apiSecret, apiKey) {
		_, _ = w.Write(statusFailedResponse(EmptyParamsErrorMsg))
		return
	}

	client, err := cfapitools.NewClient(apiKey, apiSecret)
	if err != nil {
		_, _ = w.Write(statusFailedResponse(err.Error()))
		return
	}

	if id := getIdByClient(client); id != "" {
		userId = id
	} else {
		//userId := uuid.New().String()
		setClient(userId, client)
	}

	logger.Logger().Info("Setting admin:",
		//zap.String("Handle", handle),
		//zap.String("Password", password),
		zap.String("ApiKey", apiKey),
		zap.String("ApiSecret", apiSecret),
		zap.String("UserID", userId),
	)

	jsonResp, _ := json.Marshal(&struct {
		Status string `json:"status"`
		Id     string `json:"id"`
	}{"OK", userId})

	_, _ = w.Write(jsonResp)
}

func getGroups(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")

	client := getClient(userID)

	if !validateAndWrite(w, client, userID) {
		return
	}

	logger.Logger().Info("Getting groups:",
		//zap.String("Key", client.k),
		zap.String("UserID", userID))

	groups, err := client.GetGroupsList()
	if err != nil {
		_, _ = w.Write(statusFailedResponse(err.Error()))
		return
	}

	data := statusOKResponse(groups)
	_, _ = w.Write(data)
}

func getContests(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	groupCode := r.URL.Query().Get("groupCode")

	client := getClient(userID)

	if !validateAndWrite(w, client, userID, groupCode) {
		return
	}

	logger.Logger().Info("Getting contests:",
		//zap.String("Handle", client.Handle),
		zap.String("UserID", userID),
		zap.String("GroupCode", groupCode))

	contests, err := client.GetContestsList(groupCode)
	if err != nil {
		_, _ = w.Write(statusFailedResponse(err.Error()))
		return
	}

	data := statusOKResponse(contests)
	_, _ = w.Write(data)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	groupCode := r.URL.Query().Get("groupCode")
	contestId, errId := strconv.ParseInt(r.URL.Query().Get("contestId"), 10, 64)

	client := getClient(userID)

	if !validateAndWrite(w, client, userID, groupCode, r.URL.Query().Get("contestId")) {
		return
	}
	if errId != nil {
		_, _ = w.Write(statusFailedResponse("Incorrect contest ID"))
		return
	}

	logger.Logger().Info("Getting tasks:",
		//zap.String("Handle", client.Handle),
		zap.String("UserID", userID),
		zap.String("GroupCode", groupCode),
		zap.Int64("ContestID", contestId))

	data, err := client.GetContestData(groupCode, contestId)
	if err != nil {
		_, _ = w.Write(statusFailedResponse(err.Error()))
		return
	}

	_, _ = w.Write(statusOKResponse(data))

}

func proceedProcess(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	groupCode := r.URL.Query().Get("groupCode")
	contestId, errId := strconv.ParseInt(r.URL.Query().Get("contestId"), 10, 64)

	count, errCount := strconv.Atoi(r.URL.Query().Get("count"))
	if errCount != nil {
		count = 0
	}

	weightsString := r.URL.Query().Get("weights")
	headersString := r.URL.Query().Get("headers")

	penalty, errPenalty := strconv.Atoi(r.URL.Query().Get("penalty"))
	if errPenalty != nil {
		penalty = 0
	}

	mode := r.URL.Query().Get("mode")
	if mode == "last" {
		mode = cfapitools.LastSolutionMode
	} else {
		mode = cfapitools.BestSolutionMode
	}

	lateSubmTimeString := r.URL.Query().Get("late")
	lateSubmissionSeconds, err := strconv.ParseInt(lateSubmTimeString, 10, 64)
	if err != nil {
		_, _ = w.Write(statusFailedResponse("Incorrect Unix format of `late` parameter"))
		return
	}

	client := getClient(userID)

	if !validateAndWrite(w, client, userID, groupCode, r.URL.Query().Get("contestId")) {
		return
	}
	if errId != nil {
		_, _ = w.Write(statusFailedResponse("Incorrect contest ID"))
		return
	}

	weights, err := parseWeights(weightsString)
	if err != nil {
		_, _ = w.Write(statusFailedResponse("Incorrect weights"))
		return
	}

	var headers []string
	if headersString == "" {
		headers = []string{}
	} else {
		headers = strings.Split(headersString, "-")
	}

	srcZip := "./submissions.zip"
	err = getZipFile(r, srcZip)
	if err != nil {
		_, _ = w.Write(statusFailedResponse(err.Error()))
		return
	}

	logger.Logger().Info("Proceeding:",
		//zap.String("Handle", client.Handle),
		zap.String("UserID", userID),
		zap.String("GroupCode", groupCode),
		zap.Int64("ContestID", contestId))

	extraParams := cfapitools.ParsingParameters{
		TasksWeights:          weights,
		ExtraHeaders:          headers,
		LatePenalty:           penalty,
		LateTime:              lateSubmissionSeconds,
		SubmissionParsingMode: mode,
	}

	data, err := client.GetStatistics(groupCode, contestId, count, extraParams, srcZip)
	if err != nil {
		_, _ = w.Write(statusFailedResponse(err.Error()))
		return
	}

	_, _ = w.Write(statusOKResponse(data))
	createMultipart(w, statusOKResponse(data))
}

func uploadUsers(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		logger.Error(err)
		_, _ = w.Write(statusFailedResponse("could not parse multipart form"))
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Error(err)
		_, _ = w.Write(statusFailedResponse("could not get uploaded file"))
		return
	}
	defer file.Close()

	buff, err := io.ReadAll(file)
	if err != nil {
		logger.Error(err)
		_, _ = w.Write(statusFailedResponse("could not read uploaded file"))
		return
	}

	err = db.UploadUsersToFile(buff)
	if err != nil {
		_, _ = w.Write(statusFailedResponse(err.Error()))
		return
	}

	logger.Logger().Info("Users file downloaded successfully",
		zap.Int64("size", handler.Size),
		zap.String("name", handler.Filename))

	_, _ = w.Write(statusOKResponse("file uploaded successfully"))
}

func main() {
	logger.Init()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error(err.Error())
		}
	}(logger.Logger())

	_ = godotenv.Load()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	HOST := os.Getenv("HOST")
	if HOST == "" {
		HOST = "8080"
	}

	logger.Logger().Info("Server started. Adding routes.",
		zap.String("Host", HOST),
		zap.String("Port", PORT))

	mux := mx.NewRouter()
	mux.Use(panicLogMiddleware)
	mux.Use(infoLogMiddleware)

	mux.HandleFunc("/api/setAdmin", setAdminData).Methods(http.MethodGet)
	mux.HandleFunc("/api/getTasks", getTasks).Methods(http.MethodGet)
	mux.HandleFunc("/api/getGroups", getGroups).Methods(http.MethodGet)
	mux.HandleFunc("/api/getContests", getContests).Methods(http.MethodGet)
	mux.HandleFunc("/api/proceed", proceedProcess).Methods(http.MethodPost)
	mux.HandleFunc("/api/uploadUsers", uploadUsers).Methods(http.MethodPost)

	http.Handle("/", mux)

	logger.Logger().Info("All routes are added. Start polling.",
		zap.String("Host", HOST),
		zap.String("Port", PORT))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil); err != nil {
		logger.Error(fmt.Errorf("HTTP Server error: %w", err))
	} else {
		logger.Logger().Info("Server finished work properly")
	}
}
