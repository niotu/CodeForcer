package entities

type User struct {
	Handle    string
	Solutions map[string]*Submission
}
