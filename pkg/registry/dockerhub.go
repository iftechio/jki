package registry

import (
	"fmt"

	"github.com/docker/docker/api/types"
)

// DockerHubRegistry can be any registry which is compatible with Docker Hub.
type DockerHubRegistry struct {
	Namespace string `json:"namespace" yaml:"namespace"`
	Server    string `json:"server" yaml:"server"`
	Username  string `json:"username" yaml:"username"`
	Password  string `json:"password" yaml:"password"`
}

var _ innerRegistryInterface = (*DockerHubRegistry)(nil)

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
