package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
)

var (
	ErrUnknownRegistry = fmt.Errorf("unknown registry")
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
	Name      string             `json:"name" yaml:"name"`
	AliCloud  *AliCloudRegistry  `json:"aliyun" yaml:"aliyun"`
	AWS       *AWSRegistry       `json:"aws" yaml:"aws"`
	DockerHub *DockerHubRegistry `json:"dockerhub" yaml:"dockerhub"`
}

var _ RegistryInterface = (*Registry)(nil)

var publicReg = &PublicRegistry{}

func (r *Registry) delegate() innerRegistryInterface {
	switch {
	case r.AliCloud != nil:
		return r.AliCloud
	case r.AWS != nil:
		return r.AWS
	case r.DockerHub != nil:
		return r.DockerHub
	default:
		return publicReg
	}
}

func (r *Registry) Prefix() string {
	return r.delegate().Prefix()
}

func (r *Registry) GetAuthToken() (string, error) {
	auth, err := r.GetAuthConfig()
	if err != nil {
		return "", err
	}
	return toRegistryAuth(auth.Username, auth.Password)
}

func (r *Registry) CreateRepoIfNotExists(repo string) error {
	return r.delegate().CreateRepoIfNotExists(repo)
}

func (r *Registry) GetLatestTag(repo string) (string, error) {
	return r.delegate().GetLatestTag(repo)
}

func (r *Registry) Verify() error {
	ri := r.delegate()
	if _, ok := ri.(*PublicRegistry); ok {
		return ErrUnknownRegistry
	}
	return r.delegate().Verify()
}

func (r *Registry) GetAuthConfig() (types.AuthConfig, error) {
	return r.delegate().GetAuthConfig()
}
