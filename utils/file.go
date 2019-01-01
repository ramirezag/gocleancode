package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
)

func CloseFile(file io.Closer) {
	err := file.Close()
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to close file. %v", err))
	}
}
