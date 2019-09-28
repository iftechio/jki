package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadDockerIgnore reads a .dockerignore file in `contextDir` and returns the list of file patterns
// to ignore. Note this will trim whitespace from each line as well
// as use GO's "clean" func to get the shortest/cleanest path for each.
// Excerpt from https://github.com/docker/docker/blob/master/builder/dockerignore/dockerignore.go
func ReadDockerIgnore(contextDir string) ([]string, error) {
	f, err := os.Open(filepath.Join(contextDir, ".dockerignore"))
	switch {
	case os.IsNotExist(err):
		return nil, nil
	case err != nil:
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var excludes []string
	currentLine := 0

	utf8bom := []byte{0xEF, 0xBB, 0xBF}
	for scanner.Scan() {
		scannedBytes := scanner.Bytes()
		// We trim UTF8 BOM
		if currentLine == 0 {
			scannedBytes = bytes.TrimPrefix(scannedBytes, utf8bom)
		}
		pattern := string(scannedBytes)
		currentLine++
		// Lines starting with # (comments) are ignored before processing
		if strings.HasPrefix(pattern, "#") {
			continue
		}
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		// normalize absolute paths to paths relative to the context
		// (taking care of '!' prefix)
		invert := pattern[0] == '!'
		if invert {
			pattern = strings.TrimSpace(pattern[1:])
		}
		if len(pattern) > 0 {
			pattern = filepath.Clean(pattern)
			pattern = filepath.ToSlash(pattern)
			if len(pattern) > 1 && pattern[0] == '/' {
				pattern = pattern[1:]
			}
		}
		if invert {
			pattern = "!" + pattern
		}

		excludes = append(excludes, pattern)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading .dockerignore: %v", err)
	}
	return excludes, nil
}
