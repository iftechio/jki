package utils

import (
	"log"
	"os"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func HomeDir() string {
	return os.Getenv("HOME")
}
