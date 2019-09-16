package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getOutput(dumpError bool, args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	data, err := cmd.Output()
	if err != nil && dumpError {
		if ee, ok := err.(*exec.ExitError); ok {
			c := strings.Join(args, " ")
			fmt.Fprintf(os.Stderr, "ERROR: `%s`: %s", c, string(ee.Stderr))
		}
	}
	return strings.TrimSpace(string(data)), err
}

func HasChanges() bool {
	cmd := exec.Command("git", "diff", "--quiet")
	cmd.Stderr = os.Stderr
	return cmd.Run() != nil
}

func GetOriginURL() (string, error) {
	return getOutput(true, "git", "config", "--get", "remote.origin.url")
}

func GetCurrentBranch() (string, error) {
	return getOutput(true, "git", "rev-parse", "--abbrev-ref", "HEAD")
}

func GetAbbrevCommitHash() (string, error) {
	return getOutput(true, "git", "rev-parse", "--short", "HEAD")
}

func GetTagOfCommit(commitHash string) (string, error) {
	return getOutput(false, "git", "describe", "--exact-match", "--tags", commitHash)
}
