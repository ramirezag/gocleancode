package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gocleancode/services"
	"net/http"
)

type Handlers struct {
	fileService services.FileService
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewHandlers(fileService services.FileService) *mux.Router {
	handlers := Handlers{fileService}
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, http.StatusOK, Response{true, "UP"})
	})
	r.HandleFunc("/files", handlers.UploadFile).Methods("POST")
	r.HandleFunc("/files/{fileId}", handlers.GetFileById).Methods("GET")
	r.HandleFunc("/files/{fileId}", handlers.DeleteFileById).Methods("DELETE")
	return r
}

func jsonResponse(w http.ResponseWriter, code int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(500)
		log.Error(err)
	} else {
		w.WriteHeader(code)
		_, err := w.Write(jsonBytes)
		if err != nil {
			log.Error(fmt.Sprintf("Error encountered when writing %v to response, %v", response, err))
		}
	}
}
