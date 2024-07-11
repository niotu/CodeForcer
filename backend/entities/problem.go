package entities

type Submission struct {
	Index          string
	Points         float64
	SubmissionId   int64
	ProgramLang    string
	SubmissionTime int64
	Late           bool
}

type Problem struct {
	Name      string
	Index     string
	MaxPoints float64
}
