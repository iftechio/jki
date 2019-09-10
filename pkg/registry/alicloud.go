package registry

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
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

func (r *AliCloudRegistry) Domain() string {
	return fmt.Sprintf("registry.%s.aliyuncs.com/%s", r.Region, r.Namespace)
}

func (r *AliCloudRegistry) GetAuthToken() (string, error) {
	if len(r.Username) != 0 && len(r.Password) != 0 {
		return toRegistryAuth(r.Username, r.Password)
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
			return "", fmt.Errorf("create cr client: %s", err)
		}
		req := cr.CreateGetAuthorizationTokenRequest()
		req.Domain = fmt.Sprintf("cr.%s.aliyuncs.com", r.Region)
		rawResp, err := client.GetAuthorizationToken(req)
		if err != nil {
			return "", fmt.Errorf("get token: %s", err)
		}
		var resp GetAuthTokenResponse
		err = json.Unmarshal(rawResp.GetHttpContentBytes(), &resp)
		if err != nil {
			return "", err
		}
		return toRegistryAuth(resp.Data.UserName, resp.Data.AuthorizationToken)
	}
	return "", fmt.Errorf("neither username and password nor access key and secret access key are specified")
}
