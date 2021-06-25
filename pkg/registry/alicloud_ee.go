package registry

import (
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/docker/docker/api/types"
)

type AliCloudEERegistry struct {
	AliCloudRegistry
	InstanceId   string `json:"instance_id"`
	InstanceHost string `json:"instance_host"`
}

var _ innerInterface = (*AliCloudEERegistry)(nil)

func (r *AliCloudEERegistry) CreateRepoIfNotExists(repo string) error {
	return nil
}

func (r *AliCloudEERegistry) MatchImage(image string) bool {
	if r.InstanceHost != "" {
		return strings.HasPrefix(image, r.InstanceHost)
	}
	endpoints, err := r.getEndpoints()
	if err != nil {
		return false
	}
	for _, endpoint := range endpoints {
		if strings.HasPrefix(image, endpoint.Domain) {
			return true
		}
	}
	return false
}

func (r *AliCloudEERegistry) Prefix() string {
	return fmt.Sprintf("%s/%s", r.Host(), r.Namespace)
}

func (r *AliCloudEERegistry) Host() string {
	if r.InstanceHost != "" {
		return r.InstanceHost
	}
	endpoints, _ := r.getEndpoints()
	host := ""
	for _, endpoint := range endpoints {
		if endpoint.Type == "USER" {
			host = endpoint.Domain
			break
		}
		host = endpoint.Domain
	}

	return host
}

func (r *AliCloudEERegistry) GetLatestTag(repo string) (tag string, err error) {
	client, err := r.getClient()
	if err != nil {
		return
	}

	req := cr_ee.CreateListRepoTagRequest()
	req.Domain = fmt.Sprintf("cr.%s.aliyuncs.com", r.Region)

	req.PageNo = requests.Integer(rune(1))
	req.PageSize = requests.Integer(rune(1))
	req.InstanceId = r.InstanceId

	req.RepoId, err = r.getRepoIdWithRepoName(repo)
	if err != nil {
		return
	}

	var resp *cr_ee.ListRepoTagResponse
	resp, err = client.ListRepoTag(req)
	if err != nil {
		return
	}
	if !resp.IsSuccess() {
		err = fmt.Errorf("cannot list repo tags")
		return
	}

	if resp.TotalCount == "0" {
		err = fmt.Errorf("repo has no image")
		return
	}

	if len(resp.Images) == 0 {
		err = fmt.Errorf("image has no tag")
		return
	}

	tag = resp.Images[0].Tag
	return
}

func (r *AliCloudEERegistry) GetAuthConfig() (auth types.AuthConfig, err error) {
	auth.ServerAddress = r.Prefix()
	if len(r.Username) != 0 && len(r.Password) != 0 {
		auth.Username, auth.Password = r.Username, r.Password
		return auth, nil
	}
	client, err := r.getClient()
	if err != nil {
		return
	}
	req := cr_ee.CreateGetAuthorizationTokenRequest()
	req.Domain = fmt.Sprintf("cr.%s.aliyuncs.com", r.Region)
	req.InstanceId = r.InstanceId
	resp, err := client.GetAuthorizationToken(req)
	if err != nil {
		return auth, fmt.Errorf("get token: %s", err)
	}
	auth.Username, auth.Password = resp.TempUsername, resp.AuthorizationToken
	return auth, nil
}

func (r *AliCloudEERegistry) getRepoIdWithRepoName(repoName string) (repoId string, err error) {
	client, err := r.getClient()
	if err != nil {
		return
	}
	req := cr_ee.CreateGetRepositoryRequest()
	req.Domain = fmt.Sprintf("cr.%s.aliyuncs.com", r.Region)
	req.RepoName = repoName
	req.RepoNamespaceName = r.Namespace
	req.InstanceId = r.InstanceId
	resp, err := client.GetRepository(req)
	if err != nil {
		return
	}
	return resp.RepoId, nil
}

func (r *AliCloudEERegistry) getEndpoints() (endpoint []cr_ee.Endpoints, err error) {
	client, err := r.getClient()
	if err != nil {
		return
	}

	req := cr_ee.CreateGetInstanceEndpointRequest()
	req.Domain = fmt.Sprintf("cr.%s.aliyuncs.com", r.Region)
	req.InstanceId = r.InstanceId
	req.EndpointType = "internet"
	req.ModuleName = "Registry"

	resp, err := client.GetInstanceEndpoint(req)
	if err != nil {
		return
	}
	if !resp.IsSuccess() {
		err = fmt.Errorf("cannot get instance endpoints %s", err)
		return
	}

	if len(resp.Domains) == 0 {
		err = fmt.Errorf("instance not endpoints")
		return
	}
	endpoint = resp.Domains
	return
}

func (r *AliCloudEERegistry) getClient() (client *cr_ee.Client, err error) {
	if len(r.AccessKey) != 0 && len(r.SecretAccessKey) != 0 {
		client, err = cr_ee.NewClientWithAccessKey(r.Region, r.AccessKey, r.SecretAccessKey)
	} else {
		err = fmt.Errorf("dont support username and password")
	}
	if err != nil {
		err = fmt.Errorf("create cr client: %s", err)
		return
	}
	return
}
