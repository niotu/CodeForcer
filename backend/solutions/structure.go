package solutions

import (
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"strconv"
)

var dirStructure map[string]map[string][]string

func addToDirStructure(filePath string, author entities.User) {
	filename, _ := getFileNameAndExtension(filePath)

	subId, _ := strconv.ParseInt(filename, 10, 64)

	if _, ok := dirStructure[author.Handle]; !ok {
		dirStructure[author.Handle] = make(map[string][]string)
	}

	for i, submission := range author.Solutions {
		if submission.SubmissionId == subId {
			programLang := author.Solutions[i].ProgramLang

			dirStructure[author.Handle][programLang] =
				append(dirStructure[author.Handle][programLang], filePath)
			break
		}
	}
}
