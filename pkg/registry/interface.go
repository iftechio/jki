package registry

import (
	"github.com/docker/docker/api/types"
)

type RegistryInterface interface {
	CreateRepoIfNotExists(repo string) error
	Prefix() string
	GetAuthConfig() (types.AuthConfig, error)
	GetAuthToken() (string, error)
	GetLatestTag(repo string) (string, error)
	Verify() error
}
