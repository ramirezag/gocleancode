package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gocleancode/utils"
	"io"
	"net/http"
	"os"
	"strconv"
)

func (handlers Handlers) UploadFile(w http.ResponseWriter, r *http.Request) {
	file, handle, err := r.FormFile("file")
	if err != nil {
		log.Error(err)
		jsonResponse(w, http.StatusBadRequest, Response{false, "Failed to save file!"})
		return
	}
	defer utils.CloseFile(file)
	generatedId, err := handlers.fileService.SaveFile(file, handle)
	if err != nil {
		log.Error(err)
		jsonResponse(w, http.StatusInternalServerError, Response{false, "Failed to save file!"})
	} else {
		jsonResponse(w, http.StatusCreated, Response{true, fmt.Sprintf("Created file with id %v.", generatedId)})
	}
}

func (handlers Handlers) GetFileById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileId := vars["fileId"]
	fileIdInt64, err := strconv.ParseInt(fileId, 0, 64)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, Response{false, "Unparseable fileId."})
		return
	}
	file, err := handlers.fileService.GetFileById(fileIdInt64)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, Response{false, "Failed to get file."})
		return
	}
	filePath := *file.FilePath
	actualFile, err := os.Open(filePath)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to open file %s.", filePath))
		jsonResponse(w, http.StatusInternalServerError, Response{false, "Failed to open file."})
		return
	}
	n, err := io.Copy(w, actualFile)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to copy file %s. Reason: %v", filePath, err))
		jsonResponse(w, http.StatusInternalServerError, Response{false, "Failed to get file."})
		return
	}
	w.Header().Set("Content-Type", *file.ContentType)
	w.Header().Set("Content-Disposition", "inline") // Display in browser
	w.Header().Set("Content-Length", fmt.Sprintf("%d", n))
}

func (handlers Handlers) DeleteFileById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileId := vars["fileId"]
	fileIdInt64, err := strconv.ParseInt(fileId, 0, 64)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, Response{false, "Unparseable fileId."})
		return
	}
	err = handlers.fileService.DeleteFileById(fileIdInt64)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, Response{false, "Failed to delete file with id " + fileId})
	} else {
		jsonResponse(w, http.StatusOK, Response{true, "Successfully deleted file with id " + fileId})
	}
}
