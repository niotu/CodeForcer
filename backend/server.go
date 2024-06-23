package main

import (
	"fmt"
	mx "github.com/gorilla/mux"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/backend/entities"
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

func setApiKey(w http.ResponseWriter, r *http.Request) {
	api = entities.NewClient()

	apiKey := r.URL.Query().Get("key")
	apiSecret := r.URL.Query().Get("secret")

	api.SetApiKey(apiKey)
	api.SetApiSecret(apiSecret)

	w.Write([]byte("Ok"))

}

func getGroups(w http.ResponseWriter, r *http.Request) {
	groups, client := entities.GetGroupsList(nil)

	authClient = client

	data := entities.EntitiesToJSON(groups)

	w.Write(data)
}

func getContests(w http.ResponseWriter, r *http.Request) {
	groupCode := r.URL.Query().Get("groupCode")

	contests, client := entities.GetContestsList(authClient, groupCode)
	authClient = client
	data := entities.EntitiesToJSON(contests)

	w.Write(data)
}

func proceedProcess(w http.ResponseWriter, r *http.Request) {
	groupCode := r.URL.Query().Get("groupCode")
	contestID, _ := strconv.Atoi(r.URL.Query().Get("contestID"))

	con := entities.Contest{
		Id:        contestID,
		GroupCode: groupCode,
		Problems:  nil,
	}

	data := api.ParseAndFormEntities(con)

	w.Write(data)
}

// var logger *zap.Logger
var api *entities.Client
var authClient *http.Client

func main() {
	//logger, _ = zap.NewProduction()
	//defer logger.Sync()

	mux := mx.NewRouter()

	mux.HandleFunc("/setApiKey", setApiKey)

	mux.HandleFunc("/getGroups", getGroups)
	mux.HandleFunc("/getContests", getContests)
	mux.HandleFunc("/proceed", proceedProcess)

	http.Handle("/", mux)

	http.ListenAndServe(":8080", nil)

}
