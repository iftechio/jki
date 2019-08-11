package git

import (
	"os/exec"
	"strings"
)

func getOutput(name string, arg ...string) (string, error) {
	data, err := exec.Command(name, arg...).Output()
	return strings.TrimSpace(string(data)), err
}

func HasChanges() bool {
	err := exec.Command("git", "diff", "--quiet").Run()
	return err != nil
}

func GetOriginURL() (string, error) {
	return getOutput("git", "config", "--get", "remote.origin.url")
}

func GetCurrentBranch() (string, error) {
	return getOutput("git", "rev-parse", "--abbrev-ref", "HEAD")
}

func GetAbbrevCommitHash() (string, error) {
	return getOutput("git", "rev-parse", "--short", "HEAD")
}

func GetTagOfCommit(commitHash string) (string, error) {
	return getOutput("git", "describe", "--exact-match", "--tags", commitHash)
}
