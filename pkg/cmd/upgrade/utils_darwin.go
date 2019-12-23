package upgrade

import (
	"fmt"
	"strings"
)

func validatePath(fp string) error {
	if strings.HasPrefix(fp, "/usr/local/Cellar/") {
		return fmt.Errorf("`jki` seems to be installed by homebrew, please upgrade jki using homebrew")
	}
	return nil
}
