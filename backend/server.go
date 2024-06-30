package main

import (
	"fmt"
	mx "github.com/gorilla/mux"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/cf-api-tools"
	"net/http"
	"strconv"
)

func panicLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("panic")
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)

			}
		}()
		next.ServeHTTP(w, r)
	})
}

func setAdminData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("setting admin...")

	api = cf_api_tools.NewClient()

	apiKey := r.URL.Query().Get("key")
	apiSecret := r.URL.Query().Get("secret")
	handle := r.URL.Query().Get("handle")
	password := r.URL.Query().Get("password")

	api.SetApiKey(apiKey)
	api.SetApiSecret(apiSecret)
	api.SetHandle(handle)
	api.SetPassword(password)

	w.Write([]byte("Ok"))

}

func getGroups(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getting groups...")

	groups, client := api.GetGroupsList(nil)

	authClient = client

	data := cf_api_tools.EntitiesToJSON(groups)

	w.Write(data)
}

func getContests(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getting contests...")

	groupCode := r.URL.Query().Get("groupCode")
	if groupCode == "" {
		w.Write([]byte("groupCode field malformed or does not exist"))
		return
	}

	contests, client := api.GetContestsList(authClient, groupCode)
	authClient = client
	data := cf_api_tools.EntitiesToJSON(contests)

	w.Write(data)
}

func proceedProcess(w http.ResponseWriter, r *http.Request) {
	fmt.Println("proceeding...")

	groupCode := r.URL.Query().Get("groupCode")
	contestId, errId := strconv.ParseInt(r.URL.Query().Get("contestId"), 10, 64)
	count, errCount := strconv.Atoi(r.URL.Query().Get("count"))
	if errCount != nil {
		count = 0
	}

	if errId != nil {
		w.Write([]byte("contestId field malformed or does not exist"))
		return
	}
	if groupCode == "" {
		w.Write([]byte("groupCode field malformed or does not exist"))
		return
	}

	//taskWeightsString := r.URL.Query().Get("weights")
	//if taskWeightsString == "" {
	//	w.Write([]byte("weights field malformed or does not exist"))
	//	return
	//}

	//taskWeights := strings.Split(taskWeightsString, "-")
	//weights := make([]int, len(taskWeights))
	//for i, tw := range taskWeights {
	//	intW, err := strconv.Atoi(tw)
	//	if err != nil {
	//		w.Write([]byte("weights field malformed or does not exist"))
	//		return
	//	}
	//
	//	weights[i] = intW
	//}

	weights := []int{50, 50}

	data := api.GetStatistics(nil, groupCode, contestId, count, weights)

	w.Write(data)
}

// var logger *zap.Logger
var api *cf_api_tools.Client
var authClient *http.Client

func main() {
	//logger, _ = zap.NewProduction()
	//defer logger.Sync()

	fmt.Println("Start adding routes...")

	mux := mx.NewRouter()

	mux.HandleFunc("/api/setAdmin", setAdminData)

	mux.HandleFunc("/api/getGroups", getGroups)
	mux.HandleFunc("/api/getContests", getContests)
	mux.HandleFunc("/api/proceed", proceedProcess)

	http.Handle("/", mux)

	PORT := 8080

	fmt.Printf("All routes are added. Start polling port :%d...\n", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	if err != nil {
		fmt.Println("Error caught: ", err)
	} else {
		fmt.Println("Server finished work properly")
	}

}
