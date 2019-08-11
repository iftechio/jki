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

/*
type NormalRegistry struct {
	Name     string `json:"name" yaml:"name"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func (r *NormalRegistry) GetAuthToken() (string, error) {
	return toRegistryAuth(r.Username, r.Password)
}
*/

type Registry struct {
	Name     string            `json:"name" yaml:"name"`
	AliCloud *AliCloudRegistry `json:"aliyun" yaml:"aliyun"`
	AWS      *AWSRegistry      `json:"aws" yaml:"aws"`
}

func (r *Registry) Domain() string {
	switch {
	case r.AWS != nil:
		return r.AWS.Domain()
	case r.AliCloud != nil:
		return r.AliCloud.Domain()
	}
	return ""
}

func (r *Registry) GetAuthToken() (string, error) {
	switch {
	case r.AWS != nil:
		return r.AWS.GetAuthToken()
	case r.AliCloud != nil:
		return r.AliCloud.GetAuthToken()
	}
	return "", errUnknownRegistry
}

func (r *Registry) CreateRepoIfNotExists(repo string) error {
	switch {
	case r.AWS != nil:
		return r.AWS.CreateRepoIfNotExists(repo)
	case r.AliCloud != nil:
		return r.AliCloud.CreateRepoIfNotExists(repo)
	}
	return errUnknownRegistry
}