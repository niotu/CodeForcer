package entities

import (
	"encoding/json"
	"log"
	"os"
)

type User struct {
	Handle    string
	Solutions map[string]*Submission
}

func UserListToJSON(users []User) []byte {
	data, err := json.Marshal(users)
	if err != nil {
		log.Fatal(err)
	}

	file, _ := os.OpenFile("users.json", os.O_CREATE|os.O_TRUNC, 0606)
	file.Write(data)

	return data
}
