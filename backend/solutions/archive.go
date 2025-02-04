package solutions

import (
	"errors"
	"fmt"
	arch "github.com/mholt/archiver/v3"
	"github.com/xyproto/unzip"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/db"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/entities"
	"gitlab.pg.innopolis.university/n.solomennikov/choosetwooption/backend/logger"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var resultZipName = "./result.zip"
var SolutionArchiveError = errors.New("unable to proceed operations with archive, check the correctness of .zip and try later")

func GetResultZipName(userId string) string {
	name, ext := getFileNameAndExtension(resultZipName)
	base := filepath.Dir(resultZipName)
	return filepath.Join(base, name+"_"+userId+ext)
}

func getFileNameAndExtension(filePath string) (string, string) {
	fileName := filepath.Base(filePath)
	ext := filepath.Ext(filePath)
	name := strings.TrimSuffix(fileName, ext)
	return name, ext
}

func unzipArchive(src, dest string) error {
	sp, _ := filepath.Abs(src)
	dp, _ := filepath.Abs(dest)

	err := unzip.Extract(sp, dp)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("unable to unzip the archive")
	}

	return nil
}

func moveFile(src, dst string, wg *sync.WaitGroup) {
	defer wg.Done()

	_ = os.Rename(src, dst)
}

func ParseSubmissions(dir string, authors map[int64]entities.User) error {
	dirStructure = make(map[string]map[string][]ArchiveObject)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			name, _ := getFileNameAndExtension(path)
			subId, _ := strconv.ParseInt(name, 10, 64)
			if author, ok := authors[subId]; ok {
				addToDirStructure(path, author)
			}
		}
		return nil
	})

	if err != nil {
		logger.Error(fmt.Errorf("error walking the directory: %w", err))
		return SolutionArchiveError
	}

	return nil

}

func MakeSolutionsArchive(srcArchive string, userId string, authors map[int64]entities.User) error {
	dest := "./temp_" + userId
	unarchived := "./solutions_" + userId

	defer func() {
		go func() {
			_ = os.RemoveAll(dest)
			_ = os.RemoveAll(srcArchive)
			_ = os.RemoveAll(unarchived)
		}()
	}()

	err := unzipArchive(srcArchive, dest)
	if err != nil {
		return err
	}

	err = ParseSubmissions(dest, authors)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	finalDir, _ := filepath.Abs(unarchived)

	handleEmailMap := db.GetUsers()

	_ = os.Mkdir(finalDir, 0755)
	for taskName, submissions := range dirStructure {
		curr := filepath.Join(finalDir, taskName)

		_ = os.Mkdir(curr, 0755)

		for lang, fileObjects := range submissions {
			langDir := filepath.Join(curr, lang)
			_ = os.Mkdir(langDir, 0755)
			for _, obj := range fileObjects {
				_, ext := getFileNameAndExtension(obj.Path)

				userEmail, ok := handleEmailMap[obj.Handle]
				if !ok {
					userEmail = obj.Handle
				}

				wg.Add(1)
				go moveFile(obj.Path,
					filepath.Join(langDir, userEmail+ext),
					&wg)
			}
		}
	}

	wg.Wait()

	err = arch.Archive([]string{unarchived}, GetResultZipName(userId))
	if err != nil {
		logger.Error(fmt.Errorf("failed to zip folder: %v", err))
		return SolutionArchiveError
	}

	logger.Logger().Info("Folder successfully zipped!")

	return nil
}
