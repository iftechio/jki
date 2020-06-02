package registry

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
)

// DockerHubRegistry can be any registry which is compatible with Docker Hub.
type DockerHubRegistry struct {
	Namespace string `json:"namespace"`
	Server    string `json:"server"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

var _ innerInterface = (*DockerHubRegistry)(nil)

func (r *DockerHubRegistry) CreateRepoIfNotExists(repo string) error {
	return nil
}

func (r *DockerHubRegistry) Prefix() string {
	if len(r.Server) == 0 {
		return r.Username
	}
	if len(r.Namespace) > 0 {
		return fmt.Sprintf("%s/%s", r.Server, r.Namespace)
	}
	return r.Server
}

func (r *DockerHubRegistry) MatchImage(image string) bool {
	var prefixes []string
	if len(r.Server) == 0 {
		prefixes = append(prefixes, fmt.Sprintf("%s/", r.Username), fmt.Sprintf("docker.io/%s/", r.Username))
	} else {
		prefixes = append(prefixes, fmt.Sprintf("%s/%s", r.Server, r.Namespace))
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(image, prefix) {
			return true
		}
	}
	return false
}

func (r *DockerHubRegistry) Host() string {
	if len(r.Server) == 0 {
		return "registry-1.docker.io"
	}
	return r.Server
}

func (r *DockerHubRegistry) GetLatestTag(repo string) (string, error) {
	return "latest", nil
}

func (r *DockerHubRegistry) Verify() error {
	if len(r.Username) == 0 || len(r.Password) == 0 {
		return fmt.Errorf("empty username or password")
	}
	return nil
}

func (r *DockerHubRegistry) GetAuthConfig() (types.AuthConfig, error) {
	var addr string
	if len(r.Server) == 0 {
		// However (for legacy reasons) the “official” Docker, Inc. hosted registry
		// must be specified with both a “https://” prefix and a “/v1/” suffix
		// even though Docker will prefer to use the v2 registry API.

		// See also https://docs.docker.com/engine/api/v1.40/#operation/ImageBuild
		addr = "https://index.docker.io/v1/"
	} else {
		addr = r.Server
	}
	auth := types.AuthConfig{
		ServerAddress: addr,
		Username:      r.Username,
		Password:      r.Password,
	}
	return auth, nil
}
