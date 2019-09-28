package registry

import (
	"github.com/docker/docker/api/types"
)

type innerRegistryInterface interface {
	CreateRepoIfNotExists(repo string) error
	Prefix() string
	GetAuthConfig() (types.AuthConfig, error)
	GetLatestTag(repo string) (string, error)
	Verify() error
}

type RegistryInterface interface {
	innerRegistryInterface
	GetAuthToken() (string, error)
}
