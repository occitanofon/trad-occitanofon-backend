package helpers

import (
	"log"
	"os"
)

func CreateLogFile(logPath string) *os.File {
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}
