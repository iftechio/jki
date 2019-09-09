package version

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	cmdutils "github.com/bario/jki/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

type versionInfo struct {
	GitCommit string
	BuildDate string
	Version   string
	GoVersion string
	Compiler  string
	Platform  string
}

func NewCmdVersion(f cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			info := versionInfo{
				GitCommit: gitCommit,
				BuildDate: buildDate,
				Version:   version,
				GoVersion: runtime.Version(),
				Compiler:  runtime.Compiler,
				Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			}
			enc := json.NewEncoder(os.Stdout)
			_ = enc.Encode(&info)
		},
	}
	return cmd
}
