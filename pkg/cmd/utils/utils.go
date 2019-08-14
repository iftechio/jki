package utils

import (
	"log"
	"os"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func HomeDir() string {
	return os.Getenv("HOME")
}

func ReplaceRegistry(frImg string, domain string) string {
	parts := strings.Split(frImg, "/")
	repoWithTag := parts[len(parts)-1]
	parts = strings.Split(repoWithTag, ":")
	if len(parts) == 1 {
		// missing colon
		repoWithTag += ":latest"
	} else {
	}

	toImg := domain + "/" + repoWithTag
	return toImg
}
