package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gocleancode/handlers"
	"gocleancode/repository"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	goFilePath "path/filepath"
	"strings"
	"testing"
	"time"
)

var uploadDir = os.TempDir() + "file_test/"

func TestMain(m *testing.M) {
	log.Info("Begin Test Suite")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		fmt.Println("Creating dir: " + uploadDir)
		os.MkdirAll(uploadDir, os.ModePerm)
	}
	exitCode := m.Run()
	log.Info("End Test Suite")
	os.Exit(exitCode)
}

func TestUploadFile(t *testing.T) {
	fileService, appHandlers := createHandlers()
	workingDir, _ := os.Getwd()
	testFile := workingDir + "/dragonball.jpg"
	fileContents := "This is a test."
	req, err := newfileUploadRequest("/files", testFile, fileContents)
	if err != nil {
		t.Errorf("Failed to create POST upload request %v.", err)
	}
	expectedFile, expectedHandle, err := req.FormFile("file")
	expectedGeneratedId := int64(1)
	fileService.On("SaveFile", expectedFile, expectedHandle).Return(expectedGeneratedId, nil).Once()

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// When
	appHandlers.ServeHTTP(rr, req)
	// Then
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %d want %d", status, http.StatusOK)
	}
	// Check the response body is what we expect.
	expectedResponse := handlers.Response{true, fmt.Sprintf("Created file with id %d.", expectedGeneratedId)}
	actualResponse := handlers.Response{}
	err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
		return
	}
	assert.Equal(t, expectedResponse, actualResponse)
	fileService.AssertCalled(t, "SaveFile", expectedFile, expectedHandle)
}

func TestGetFileById(t *testing.T) {
	fileService, appHandlers := createHandlers()
	fileId := int64(1)
	url := fmt.Sprintf("/files/%d", fileId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	expectedRespBody := []byte("hello world")
	expectedRespContentType := "image/jpg"
	fileName := "fname"
	filePath := uploadDir + "TestGetFileById.txt"
	createdDt := time.Now()
	file := repository.File{FileName: &fileName, FilePath: &filePath, ContentType: &expectedRespContentType, CreatedDt: &createdDt}
	fileService.On("GetFileById", fileId).Return(file, nil).Once()

	err = ioutil.WriteFile(filePath, expectedRespBody, 0666)
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
		return
	}
	// When
	appHandlers.ServeHTTP(rr, req)
	// Then
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Check the response body is what we expect.
	headers := rr.Header()
	assert.Equal(t, expectedRespBody, rr.Body.Bytes())
	assert.Equal(t, expectedRespContentType, headers.Get("Content-Type"))
	assert.Equal(t, "inline", headers.Get("Content-Disposition"))
	assert.Equal(t, fmt.Sprintf("%d", len(expectedRespBody)), headers.Get("Content-Length"))
	fileService.AssertCalled(t, "GetFileById", fileId)
}

func TestDeleteFileById(t *testing.T) {
	// Given
	fileService, appHandlers := createHandlers()
	fileId := int64(1)
	url := fmt.Sprintf("/files/%d", fileId)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	fileService.On("DeleteFileById", fileId).Return(nil).Once()
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// When
	appHandlers.ServeHTTP(rr, req)
	// Then
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %d want %d", status, http.StatusOK)
	}
	// Check the response body is what we expect.
	expectedResponse := handlers.Response{true, fmt.Sprintf("Successfully deleted file with id %d", fileId)}
	actualResponse := handlers.Response{}
	err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
		return
	}
	assert.Equal(t, expectedResponse, actualResponse)
	fileService.AssertCalled(t, "DeleteFileById", fileId)
}

func newfileUploadRequest(uri string, filePath string, fileContents string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", goFilePath.Base(filePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, strings.NewReader(fileContents))
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
