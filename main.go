package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gocleancode/config"
	ivdnDb "gocleancode/db"
	"gocleancode/handlers"
	"gocleancode/repository"
	ivdnService "gocleancode/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	// Set logrus config that will be used accross the app
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func main() {
	server := NewServer()
	// Graceful shutdown of the server on kill or CTRL+C
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		sig := <-gracefulStop
		log.Info(fmt.Sprintf("Caught sig: %+v", sig))
		err := server.Shutdown(nil)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to shutdown the server - %v", err))
		}
	}()
	log.Info("Server started at " + server.Addr)
	log.Error(server.ListenAndServe())
}

func NewServer() *http.Server {
	appConfig := config.New()
	if appConfig.DbUser == "" || appConfig.DbPass == "" || appConfig.DbName == "" {
		panic("Database configs are required.")
	}
	mysqlDb, err := ivdnDb.NewMysqlDb(appConfig)
	if err != nil {
		panic("Failed to connect to database.")
	}
	fileRepo := repository.NewFileRepo(mysqlDb)
	fileService := ivdnService.NewFileService(mysqlDb, fileRepo, appConfig)
	appHandlers := handlers.NewHandlers(fileService)
	serverUrl := fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port)
	server := &http.Server{Addr: serverUrl, Handler: appHandlers}
	server.RegisterOnShutdown(func() {
		err := mysqlDb.Close()
		if err != nil {
			log.Error(fmt.Sprintf("Failed to close %v db. %v", mysqlDb.DataSourceName, err))
		}
		log.Info("Server stopped")
	})
	return server
}
