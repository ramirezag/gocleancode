package services

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gocleancode/config"
	"gocleancode/db"
	"gocleancode/repository"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type FileService interface {
	SaveFile(file multipart.File, handle *multipart.FileHeader) (int64, error)
	GetFileById(id int64) (repository.File, error)
	DeleteFileById(fileId int64) error
}

type fileService struct {
	db      db.Db
	repo    repository.FileRepo
	config  config.Configuration
	fileDir string
}

func NewFileService(db db.Db, repo repository.FileRepo, config config.Configuration) FileService {
	fileDir := config.UploadDir
	if !strings.HasPrefix(fileDir, os.TempDir()) {
		workingDir, _ := os.Getwd()
		fileDir = workingDir + "/" + config.UploadDir
	}
	if _, err := os.Stat(config.UploadDir); os.IsNotExist(err) {
		fmt.Println("Creating upload file dir: " + fileDir)
		os.MkdirAll(fileDir, os.ModePerm)
	}
	return fileService{db, repo, config, fileDir}
}

func (f fileService) SaveFile(multiPartFile multipart.File, fileHeader *multipart.FileHeader) (int64, error) {
	var generatedId int64
	t := time.Now()
	formattedDate := t.Format(time.RFC3339) // Eg output -> 2018-12-06T05:46:29+09:00
	fileName := fileHeader.Filename
	filePath := f.fileDir + "/" + formattedDate + "_" + fileName
	log.Info("Saving file " + filePath)
	data, err := ioutil.ReadAll(multiPartFile)
	if err != nil {
		return generatedId, err
	}
	err = ioutil.WriteFile(filePath, data, 0666)
	if err != nil {
		return generatedId, err
	}
	contentType := strings.ToLower(fileHeader.Header.Get("Content-Type"))
	now := time.Now()
	file := repository.File{FileName: &fileName, FilePath: &filePath, ContentType: &contentType, CreatedDt: &now}
	generatedId, err = f.repo.SaveFile(file)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to save file %v to DB.", file.FilePath))
		os.Remove(filePath) // If an error encountered while saving, delete the created file.
	} else {
		log.Debug(fmt.Sprintf("Successfully saved file %v to DB. Generated id is %d.", file.FilePath, generatedId))
	}
	return generatedId, err
}

func (f fileService) DeleteFileById(fileId int64) error {
	file, err := f.repo.GetFileById(fileId)
	if err != nil {
		log.Error(err)
		return err
	}
	return f.db.Transact(func(tx *sql.Tx) error {
		err = f.repo.TxDeleteFileById(fileId, tx)
		if err != nil {
			log.Error(err)
			return err
		} else {
			// Delete the item
			err = os.Remove(*file.FilePath)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to file %v. Reason: %v", file.FilePath, err))
				return err
			}
		}
		log.Info(fmt.Sprintf("Successfully deleted file with id %v", fileId))
		return nil
	})
}

func (f fileService) GetFileById(id int64) (repository.File, error) {
	return f.repo.GetFileById(id)
}
