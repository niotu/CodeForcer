package solutions

import (
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"strconv"
)

type ArchiveObject struct {
	Handle string
	Path   string
}

var dirStructure map[string]map[string][]ArchiveObject

func addToDirStructure(filePath string, author entities.User) {
	filename, _ := getFileNameAndExtension(filePath)

	subId, _ := strconv.ParseInt(filename, 10, 64)

	var taskName, programLang string

	for i, submission := range author.Solutions {
		if submission.SubmissionId == subId {
			programLang = author.Solutions[i].ProgramLang
			taskName = "Task " + submission.Index

			break
		}
	}

	if _, ok := dirStructure[taskName]; !ok {
		dirStructure[taskName] = make(map[string][]ArchiveObject)
	}
	dirStructure[taskName][programLang] =
		append(dirStructure[taskName][programLang], ArchiveObject{
			Handle: author.Handle,
			Path:   filePath,
		})
}
