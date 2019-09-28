package registry

import (
	"github.com/docker/docker/api/types"
)

// PublicRegistry represents registries which does not need authentication.
type PublicRegistry struct {
}

var _ innerRegistryInterface = (*PublicRegistry)(nil)

func (r *PublicRegistry) CreateRepoIfNotExists(repo string) error {
	return nil
}

func (r *PublicRegistry) Prefix() string {
	return ""
}

func (r *PublicRegistry) GetLatestTag(repo string) (string, error) {
	return "latest", nil
}

func (r *PublicRegistry) Verify() error {
	return nil
}

func (r *PublicRegistry) GetAuthConfig() (types.AuthConfig, error) {
	return types.AuthConfig{}, nil
}
