package main

// server.go is entry point of this HTTP sever program

import (
	"encoding/json"
	"errors"
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

var clients sync.Map       // Concurrent map to store clients ("userId":*cf_api_tools.Client)
var clientKeyToId sync.Map // Concurrent map to store client (reversed version of clients) ("apiKey":"userId")

func ConvertClientsToMap() map[string]interface{} {
	defaultMap := make(map[string]interface{})

	clients.Range(func(userId, client any) bool {
		key := userId.(string)
		defaultMap[key] = *client.(*cfapitools.Client)
		return true
	})

	return defaultMap
}

func SetUpClientsFromFile(data []byte) {
	defaultMap := make(map[string]cfapitools.Client)

	_ = json.Unmarshal(data, &defaultMap)
	for k, v := range defaultMap {
		client := v
		setClient(k, &client)
	}
}

// getClient returns pointer to cf_api_tools.Client object if found, nil otherwise
func getClient(userID string) *cfapitools.Client {
	client, ok := clients.Load(userID)
	if !ok {
		return nil
	}
	return client.(*cfapitools.Client)
}

// getIdByClient returns userID of given client if it is already existing
func getIdByClient(client *cfapitools.Client) string {
	key := client.DecodeApiKey()

	id, ok := clientKeyToId.Load(key)
	if !ok {
		return ""
	}
	return id.(string)
}

// setClient sets a pointer to the cf_api_tools.Client object to clients and updates database file
func setClient(userID string, client *cfapitools.Client) {
	key := client.DecodeApiKey()
	clientKeyToId.Store(key, userID)
	clients.Store(userID, client)

	db.UploadClientsToFile(ConvertClientsToMap())
}

// corsMiddleware is a middleware function that sets appropriate headers to http.ResponseWriter object
// to allow origins for CORS policy
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// panicLogMiddleware is a middleware function that returns corresponding message if server error occurred
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

// infoLogMiddleware is a middleware function that write to logs about every request to server
func infoLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger().Info("Request:",
			zap.String("Method", r.Method),
			zap.String("Path", r.URL.Path),
			zap.String("Query", r.URL.Query().Encode()))
		next.ServeHTTP(w, r)
	})
}

// setAdminData is a handler function for /setAdmin route.
// It gets `key` and `secret` provided in request as parameters and creates new userId
// if such user does not exist, or return existing userId otherwise.
// returns {"status": "OK", "userId":userId}
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
		zap.String("apiKey", apiKey),
		zap.String("apiSecret", apiSecret),
		zap.String("UserID", userId),
	)

	jsonResp, _ := json.Marshal(&struct {
		Status string `json:"status"`
		Id     string `json:"id"`
	}{"OK", userId})

	_, _ = w.Write(jsonResp)
}

// getGroups is a handler function for /getGroups route.
// It returns json object with all CodeForces groups of given user
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

// getContests is a handler function for /getContests route.
// It returns json object with all CodeForces contests of given group
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

// getTasks is a handler function for /getTasks route.
// It returns json object with all problems of given contest
// Requires `userID`, `groupCode`, and `contestId` parameters in query
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

// proceedProcess is a handler function for /proceed route.
// It handles all submissions of given contest, according to given parameters.
// It returns json object with all tasks' details (entities.Problem), students' submissions (entities.Submission),
// link to filled GoogleSheet, and its CSV version.
// Requires:
//
//	`userID`, `groupCode`, `contestId`,
//	`weights`: moodle weights of each problem in points, separated with "-";
//	`mode`: best/last - which submission of student need to handle;
//	`late`: late submission duration in hours;
//	`penalty`: penalty of late submission in percents;
//
// Not required:
//
//	`headers`: additional columns in google sheets table
//
// Not required attachment:
//
//	.zip file in multipart/form-data, key is "file": archive with all submissions of given contest downloaded from CF
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
	lateSubmissionHours, err := strconv.Atoi(lateSubmTimeString)
	if err != nil {
		_, _ = w.Write(statusFailedResponse("Incorrect format of `late` parameter"))
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

	// try to get attached zip file and write to file `srcZip` in root directory
	srcZip := fmt.Sprintf("./submissions_%s.zip", userID)
	getZipErr := getZipFile(r, srcZip) // if no file provided in multipart just continue without exiting with error
	if getZipErr != nil && !errors.Is(getZipErr, NoFileProvided) {
		_, _ = w.Write(statusFailedResponse(getZipErr.Error()))
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
		LateDurationSeconds:   int64(lateSubmissionHours) * 3600,
		SubmissionParsingMode: mode,
	}

	data, err := client.GetStatistics(groupCode, contestId, count, extraParams)
	if err != nil {
		_, _ = w.Write(statusFailedResponse(err.Error()))
		return
	}

	// if getZipErr is nil, then start handle submissions archive and write json and zip to multipart
	if !errors.Is(getZipErr, NoFileProvided) {
		err = cfapitools.GetSolutions(srcZip, userID, data)
		if err != nil {
			_, _ = w.Write(statusFailedResponse(err.Error()))
			return
		}

		createMultipart(w, statusOKResponse(data), userID)
	} else { // otherwise just write json to standard body
		_, _ = w.Write(statusOKResponse(data))
	}
}

// uploadUsers is a handler function to /uploadUsers route.
// It parses uploaded to multipart/form-data .csv file with handles and emails of all students
// and writes them to database file in json format
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
	// initialize logger
	logger.Init()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error(err.Error())
		}
	}(logger.Logger())

	// setUp clients from database file
	SetUpClientsFromFile(db.GetClientsBytes())

	// load environment variables
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

	// add middlewares
	mux := mx.NewRouter()
	mux.Use(panicLogMiddleware)
	mux.Use(corsMiddleware)
	mux.Use(infoLogMiddleware)

	// set up routes
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
