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
	groups, client := api.GetGroupsList(nil)

	authClient = client

	data := cf_api_tools.EntitiesToJSON(groups)

	w.Write(data)
}

func getContests(w http.ResponseWriter, r *http.Request) {
	groupCode := r.URL.Query().Get("groupCode")

	contests, client := api.GetContestsList(authClient, groupCode)
	authClient = client
	data := cf_api_tools.EntitiesToJSON(contests)

	w.Write(data)
}

func proceedProcess(w http.ResponseWriter, r *http.Request) {
	groupCode := r.URL.Query().Get("groupCode")
	contestID, _ := strconv.ParseInt(r.URL.Query().Get("contestID"), 10, 64)

	data := api.GetStatistics(nil, groupCode, contestID)

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

	mux.HandleFunc("/setAdmin", setAdminData)

	mux.HandleFunc("/getGroups", getGroups)
	mux.HandleFunc("/getContests", getContests)
	mux.HandleFunc("/proceed", proceedProcess)

	http.Handle("/", mux)

	fmt.Println("All routes are added. Start polling port :8080...")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error caught: ", err)
	} else {
		fmt.Println("Server finished work properly")
	}

}
