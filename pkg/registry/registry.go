package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
)

var (
	errUnknownRegistry = fmt.Errorf("unknown registry")
)

func toRegistryAuth(user, passwd string) (string, error) {
	authConfig := types.AuthConfig{
		Username: user,
		Password: passwd,
	}
	data, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

type Registry struct {
	Name     string            `json:"name" yaml:"name"`
	AliCloud *AliCloudRegistry `json:"aliyun" yaml:"aliyun"`
	AWS      *AWSRegistry      `json:"aws" yaml:"aws"`
}

var _ RegistryInterface = &Registry{}

func (r *Registry) registryInterface() RegistryInterface {
	if r.AliCloud != nil {
		return r.AliCloud
	}
	if r.AWS != nil {
		return r.AWS
	}
	panic(errUnknownRegistry)
}

func (r *Registry) Prefix() string {
	return r.registryInterface().Prefix()
}

func (r *Registry) GetAuthToken() (string, error) {
	return r.registryInterface().GetAuthToken()
}

func (r *Registry) CreateRepoIfNotExists(repo string) error {
	return r.registryInterface().CreateRepoIfNotExists(repo)
}

func (r *Registry) GetLatestTag(repo string) (string, error) {
	return r.registryInterface().GetLatestTag(repo)
}

func (r *Registry) Verify() error {
	return r.registryInterface().Verify()
}

func (r *Registry) GetAuthConfig() (types.AuthConfig, error) {
	return r.registryInterface().GetAuthConfig()
}
