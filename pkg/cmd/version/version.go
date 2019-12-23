package version

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/iftechio/jki/pkg/factory"
	"github.com/iftechio/jki/pkg/info"

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

func NewCmdVersion(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			vi := versionInfo{
				GitCommit: info.GitCommit,
				BuildDate: info.BuildDate,
				Version:   info.Version,
				GoVersion: runtime.Version(),
				Compiler:  runtime.Compiler,
				Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			}
			enc := json.NewEncoder(os.Stdout)
			_ = enc.Encode(&vi)
		},
	}
	return cmd
}
