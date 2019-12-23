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
	h, err := os.UserHomeDir()
	if err != nil {
		h = "/"
	}
	return h
}

func Prompt(hint string) string {
	fmt.Print(hint)
	var input string
	_, _ = fmt.Scanln(&input)
	return input
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

// ConvertKVStringsToMap converts ["key=value"] to {"key":"value"}
// Credit to https://github.com/docker/cli/blob/ebca1413117a3fcb81c89d6be226dcec74e5289f/opts/parse.go#L41
func ConvertKVStringsToMap(values []string) map[string]string {
	result := make(map[string]string, len(values))
	for _, value := range values {
		kv := strings.SplitN(value, "=", 2)
		if len(kv) == 1 {
			result[kv[0]] = ""
		} else {
			result[kv[0]] = kv[1]
		}
	}

	return result
}

// ConvertKVStringsToMapWithNil converts ["key=value"] to {"key":"value"}
// but set unset keys to nil - meaning the ones with no "=" in them.
// We use this in cases where we need to distinguish between
//   FOO=  and FOO
// where the latter case just means FOO was mentioned but not given a value
// Credit to: https://github.com/docker/cli/blob/ebca1413117a3fcb81c89d6be226dcec74e5289f/opts/parse.go#L60
func ConvertKVStringsToMapWithNil(values []string) map[string]*string {
	result := make(map[string]*string, len(values))
	for _, value := range values {
		kv := strings.SplitN(value, "=", 2)
		if len(kv) == 1 {
			result[kv[0]] = nil
		} else {
			result[kv[0]] = &kv[1]
		}
	}

	return result
}
