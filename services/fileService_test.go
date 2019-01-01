package services_test

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gocleancode/config"
	mockDb "gocleancode/db/mocks"
	"gocleancode/repository"
	mockRepos "gocleancode/repository/mocks"
	"gocleancode/services"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"
	"testing"
	"time"
)

var uploadDir = os.TempDir() + "fileService_test/"

func createFileService() (*mockDb.Db, *mockRepos.FileRepo, services.FileService) {
	fileRepo := &mockRepos.FileRepo{}
	db := &mockDb.Db{}
	appConfig := config.Configuration{UploadDir: uploadDir}
	fileService := services.NewFileService(db, fileRepo, appConfig)
	return db, fileRepo, fileService
}

func TestSaveFile(t *testing.T) {
	type MockFile struct {
		io.Reader
		io.ReaderAt
		io.Seeker
		io.Closer
	}
	fileContents := "This is a test."
	fileToSave := &MockFile{}
	fileToSave.Reader = strings.NewReader(fileContents)
	fileName := "TestSaveFile.txt"
	contentType := "text/plain"
	header := textproto.MIMEHeader{}
	header.Add("Content-Type", contentType)
	fileHeader := &multipart.FileHeader{Filename: fileName, Header: header}
	_, fileRepo, fileService := createFileService()
	fileParamMatcher := mock.MatchedBy(func(f repository.File) bool {
		fileNameMatched := *f.FileName == fileName
		filePath := *f.FilePath
		filePathMatched := strings.HasPrefix(filePath, uploadDir) && strings.HasSuffix(filePath, fileName)
		contentTypeMatched := *f.ContentType == contentType
		return fileNameMatched && filePathMatched && contentTypeMatched
	})
	expectedGeneratedId := int64(8)
	fileRepo.On("SaveFile", fileParamMatcher).Return(expectedGeneratedId, nil).Once()
	actualGeneratedId, err := fileService.SaveFile(fileToSave, fileHeader)
	assert.Nil(t, err)
	assert.Equal(t, expectedGeneratedId, actualGeneratedId)
	fileRepo.AssertCalled(t, "SaveFile", fileParamMatcher)
}

func TestDeleteFileById(t *testing.T) {
	// Given
	db, fileRepo, fileService := createFileService()
	fileId := int64(1)
	fileName := "fname"
	filePath := uploadDir + "TestDeleteFileById.txt"
	contentType := "contentType"
	createdDt := time.Now()
	file := repository.File{FileName: &fileName, FilePath: &filePath, ContentType: &contentType, CreatedDt: &createdDt}
	fileRepo.On("GetFileById", fileId).Return(file, nil).Once()
	tx := &sql.Tx{}
	db.On("Transact", mock.Anything).Return(func(f func(*sql.Tx) error) error {
		return f(tx)
	}).Once()
	fileRepo.On("TxDeleteFileById", fileId, tx).Return(nil).Once()
	_, err := os.Create(filePath)
	if err != nil {
		t.Errorf("Expected no error in creating file, but got %s instead", err)
	}
	// When
	err = fileService.DeleteFileById(fileId)
	// Then
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
		return
	}
	fileRepo.AssertCalled(t, "GetFileById", fileId)
	db.AssertCalled(t, "Transact", mock.Anything)
	fileRepo.AssertCalled(t, "TxDeleteFileById", fileId, mock.AnythingOfType("*sql.Tx"))
	_, err = os.Stat(filePath)
	assert.True(t, os.IsNotExist(err))
}

func TestGetFileById(t *testing.T) {
	id := int64(1)
	_, fileRepo, fileService := createFileService()
	expectedFile := repository.File{}
	fileRepo.On("GetFileById", id).Return(expectedFile, nil, nil).Once()
	actualFile, err := fileService.GetFileById(id)
	assert.Nil(t, err)
	assert.Equal(t, expectedFile, actualFile)
	fileRepo.AssertCalled(t, "GetFileById", id)
}
