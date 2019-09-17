package registry

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/docker/docker/api/types"
)

type AliCloudRegistry struct {
	Region          string `json:"region" yaml:"region"`
	Namespace       string `json:"namespace" yaml:"namespace"`
	Username        string `json:"username" yaml:"username"`
	Password        string `json:"password" yaml:"password"`
	AccessKey       string `json:"access_key" yaml:"access_key"`
	SecretAccessKey string `json:"secret_access_key" yaml:"secret_access_key"`
}

var _ RegistryInterface = &AliCloudRegistry{}

func (r *AliCloudRegistry) CreateRepoIfNotExists(repo string) error {
	return nil
}

func (r *AliCloudRegistry) Prefix() string {
	return fmt.Sprintf("registry.%s.aliyuncs.com/%s", r.Region, r.Namespace)
}

func (r *AliCloudRegistry) GetAuthToken() (string, error) {
	auth, err := r.GetAuthConfig()
	if err != nil {
		return "", err
	}
	return toRegistryAuth(auth.Username, auth.Password)
}

func (r *AliCloudRegistry) GetLatestTag(repo string) (tag string, err error) {
	var client *cr.Client
	if len(r.AccessKey) != 0 && len(r.SecretAccessKey) != 0 {
		client, err = cr.NewClientWithAccessKey(r.Region, r.AccessKey, r.SecretAccessKey)
	} else {
		err = fmt.Errorf("dont support username and password")
	}
	if err != nil {
		err = fmt.Errorf("create cr client: %s", err)
		return
	}

	req := cr.CreateGetRepoTagsRequest()
	req.Domain = fmt.Sprintf("cr.%s.aliyuncs.com", r.Region)
	req.RepoNamespace = r.Namespace
	req.RepoName = repo

	var rawResp *cr.GetRepoTagsResponse
	rawResp, err = client.GetRepoTags(req)
	if err != nil {
		return
	}
	var resp struct {
		Data struct {
			Total int `json:"total"`
			Page  int `json:"page"`
			Tags  []struct {
				Status      string `json:"status"`
				Digest      string `json:"digest"`
				ImageId     string `json:"imageId"`
				ImageCreate int    `json:"imageCreate"`
				Tag         string `json:"tag"`
				ImageSize   int    `json:"imageSize"`
			} `json:"tags"`
		} `json:"data"`
	}
	err = json.Unmarshal(rawResp.GetHttpContentBytes(), &resp)
	if err != nil {
		return
	}

	if resp.Data.Total == 0 {
		err = fmt.Errorf("repo has no image")
		return
	}
	if len(resp.Data.Tags) == 0 {
		err = fmt.Errorf("image has no tag")
		return
	}

	tag = resp.Data.Tags[0].Tag
	return
}

func (r *AliCloudRegistry) Verify() error {
	isNotEmpty := func(s string) bool {
		return len(s) != 0
	}

	tocheck := []struct {
		name, value string
	}{
		{
			name:  "region",
			value: r.Region,
		},
		{
			name:  "namespace",
			value: r.Namespace,
		},
	}
	for _, c := range tocheck {
		if !isNotEmpty(c.value) {
			return fmt.Errorf("%s cannot be empty", c.name)
		}
	}

	if !((isNotEmpty(r.Username) && isNotEmpty(r.Password)) || (isNotEmpty(r.AccessKey) && isNotEmpty(r.SecretAccessKey))) {
		return fmt.Errorf("neither username and password nor access_key and secret_access_key are specified")
	}

	return nil
}

func (r *AliCloudRegistry) GetAuthConfig() (auth types.AuthConfig, err error) {
	auth.ServerAddress = fmt.Sprintf("registry.%s.aliyuncs.com", r.Region)
	if len(r.Username) != 0 && len(r.Password) != 0 {
		auth.Username, auth.Password = r.Username, r.Password
		return auth, nil
	}
	if len(r.AccessKey) != 0 && len(r.SecretAccessKey) != 0 {
		type GetAuthTokenResponse struct {
			Data struct {
				AuthorizationToken string `json:"authorizationToken"`
				UserName           string `json:"tempUserName"`
			} `json:"data"`
		}
		client, err := cr.NewClientWithAccessKey(r.Region, r.AccessKey, r.SecretAccessKey)
		if err != nil {
			return auth, fmt.Errorf("create cr client: %s", err)
		}
		req := cr.CreateGetAuthorizationTokenRequest()
		req.Domain = fmt.Sprintf("cr.%s.aliyuncs.com", r.Region)
		rawResp, err := client.GetAuthorizationToken(req)
		if err != nil {
			return auth, fmt.Errorf("get token: %s", err)
		}
		var resp GetAuthTokenResponse
		err = json.Unmarshal(rawResp.GetHttpContentBytes(), &resp)
		if err != nil {
			return auth, err
		}
		auth.Username, auth.Password = resp.Data.UserName, resp.Data.AuthorizationToken
		return auth, nil
	}
	return auth, fmt.Errorf("neither username and password nor access_key and secret_access_key are specified")
}
