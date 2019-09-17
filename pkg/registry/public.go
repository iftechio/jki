package registry

import (
	"github.com/docker/docker/api/types"
)

// PublicRegistry represents registries which does not need authentication.
type PublicRegistry struct {
}

var _ RegistryInterface = &PublicRegistry{}

func (r *PublicRegistry) CreateRepoIfNotExists(repo string) error {
	return nil
}

func (r *PublicRegistry) Prefix() string {
	return ""
}

func (r *PublicRegistry) GetAuthToken() (string, error) {
	return "", nil
}

func (r *PublicRegistry) GetLatestTag(repo string) (string, error) {
	return "latest", nil
}

func (r *PublicRegistry) Verify() error {
	return nil
}

func (r *PublicRegistry) GetAuthConfig() (types.AuthConfig, error) {
	// However (for legacy reasons) the “official” Docker, Inc. hosted registry
	// must be specified with both a “https://” prefix and a “/v1/” suffix
	// even though Docker will prefer to use the v2 registry API.

	// See also https://docs.docker.com/engine/api/v1.40/#operation/ImageBuild
	auth := types.AuthConfig{
		ServerAddress: "https://index.docker.io/v1/",
	}
	return auth, nil
}
