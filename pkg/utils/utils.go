package utils

import (
	"bufio"
	"fmt"
	"io"
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

func ExtractBaseImages(input io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(input)
	var ret []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "FROM") {
			continue
		}
		i := strings.IndexRune(line, '#')
		if i != -1 {
			line = line[:i]
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid spec: %s", line)
		}
		ret = append(ret, parts[1])
	}
	return ret, nil
}
