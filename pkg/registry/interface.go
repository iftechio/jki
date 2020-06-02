package registry

import (
	"github.com/docker/docker/api/types"
)

type innerInterface interface {
	CreateRepoIfNotExists(repo string) error
	Prefix() string
	Host() string
	GetAuthConfig() (types.AuthConfig, error)
	GetLatestTag(repo string) (string, error)
	MatchImage(image string) bool
	Verify() error
}

type Interface interface {
	innerInterface
	GetAuthToken() (string, error)
}
