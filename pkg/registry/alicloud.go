package registry

import (
	"fmt"
)

type AliCloudRegistry struct {
	Region          string `json:"region" yaml:"region"`
	Namespace       string `json:"namespace" yaml:"namespace"`
	Username        string `json:"username" yaml:"username"`
	Password        string `json:"password" yaml:"password"`
	AccessKey       string `json:"access_key" yaml:"access_key"`
	SecretAccessKey string `json:"secret_access_key" yaml:"secret_access_key"`
}

func (r *AliCloudRegistry) CreateRepoIfNotExists(repo string) error {
	return nil
}

func (r *AliCloudRegistry) Domain() string {
	return fmt.Sprintf("registry.%s.aliyuncs.com/%s", r.Region, r.Namespace)
}

func (r *AliCloudRegistry) GetAuthToken() (string, error) {
	// TODO: get token using access_key and secret_access_key
	if len(r.Username) != 0 && len(r.Password) != 0 {
		return toRegistryAuth(r.Username, r.Password)
	}
	return "", fmt.Errorf("not implmented")
}
