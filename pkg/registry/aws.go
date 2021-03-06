package registry

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/docker/docker/api/types"
)

type AWSRegistry struct {
	Region              string `json:"region"`
	AccountID           string `json:"account_id"`
	AccessKey           string `json:"access_key"`
	SecretAccessKey     string `json:"secret_access_key"`
	LifecyclePolicyText string `json:"lifecycle_policy_text"`
}

var _ innerInterface = (*AWSRegistry)(nil)

func newAWSSession(region, accessKey, secretKey string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
}

func (r *AWSRegistry) CreateRepoIfNotExists(repo string) error {
	sess, err := newAWSSession(r.Region, r.AccessKey, r.SecretAccessKey)
	if err != nil {
		return err
	}

	ecrSvc := ecr.New(sess)
	input := ecr.DescribeRepositoriesInput{
		RepositoryNames: aws.StringSlice([]string{repo}),
	}
	output, err := ecrSvc.DescribeRepositories(&input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != ecr.ErrCodeRepositoryNotFoundException {
				return err
			}
		} else {
			return err
		}
	}
	if len(output.Repositories) != 0 {
		return nil
	}

	createInput := ecr.CreateRepositoryInput{
		RepositoryName: aws.String(repo),
	}
	_, err = ecrSvc.CreateRepository(&createInput)
	if err != nil {
		return err
	}

	var policy string
	if len(r.LifecyclePolicyText) != 0 {
		policy = r.LifecyclePolicyText
	} else {
		policy = defaultLifecyclePolicy
	}
	policyInput := ecr.PutLifecyclePolicyInput{
		RepositoryName:      aws.String(repo),
		LifecyclePolicyText: aws.String(policy),
	}
	_, err = ecrSvc.PutLifecyclePolicy(&policyInput)
	if err != nil {
		return err
	}
	return nil
}

func (r *AWSRegistry) Prefix() string {
	domain := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", r.AccountID, r.Region)
	if strings.HasPrefix(r.Region, "cn-") {
		return domain + ".cn"
	}
	return domain
}

func (r *AWSRegistry) MatchImage(image string) bool {
	domain := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", r.AccountID, r.Region)
	if strings.HasPrefix(r.Region, "cn-") {
		domain += ".cn"
	}
	return strings.HasPrefix(image, domain)
}

func (r *AWSRegistry) Host() string {
	return r.Prefix()
}

func (r *AWSRegistry) GetLatestTag(repo string) (tag string, err error) {
	sess, err := newAWSSession(r.Region, r.AccessKey, r.SecretAccessKey)
	if err != nil {
		return
	}

	ecrSvc := ecr.New(sess)
	input := &ecr.DescribeImagesInput{
		RepositoryName: &repo,
	}
	var output *ecr.DescribeImagesOutput
	output, err = ecrSvc.DescribeImages(input)
	if err != nil {
		return
	}
	if len(output.ImageDetails) == 0 {
		err = fmt.Errorf("repo has no image")
		return
	}
	// The results returned are unordered
	sort.Slice(output.ImageDetails, func(i, j int) bool {
		return output.ImageDetails[i].ImagePushedAt.Before(*output.ImageDetails[j].ImagePushedAt)
	})

	detail := output.ImageDetails[len(output.ImageDetails)-1]
	if len(detail.ImageTags) == 0 {
		err = fmt.Errorf("image has no tag")
		return
	}

	tag = *detail.ImageTags[0]
	return
}

func (r *AWSRegistry) Verify() error {
	tocheck := []struct {
		name, value string
	}{
		{
			name:  "region",
			value: r.Region,
		},
		{
			name:  "account_id",
			value: r.AccountID,
		},
		{
			name:  "access_key",
			value: r.AccessKey,
		},
		{
			name:  "secret_access_key",
			value: r.SecretAccessKey,
		},
	}
	for _, c := range tocheck {
		if len(c.value) == 0 {
			return fmt.Errorf("%s cannot be empty", c.name)
		}
	}
	return nil
}

func (r *AWSRegistry) GetAuthConfig() (auth types.AuthConfig, err error) {
	sess, err := newAWSSession(r.Region, r.AccessKey, r.SecretAccessKey)
	if err != nil {
		return
	}

	ecrSvc := ecr.New(sess)
	output, err := ecrSvc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return
	}
	if len(output.AuthorizationData) < 1 {
		return auth, fmt.Errorf("missing token from ecr")
	}

	token := output.AuthorizationData[0]
	encodedToken := aws.StringValue(token.AuthorizationToken)
	data, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		return auth, fmt.Errorf("decode ecr token: %s", err)
	}
	parts := strings.Split(string(data), ":")
	auth.Username, auth.Password = parts[0], parts[1]
	auth.ServerAddress = r.Prefix()
	return auth, nil
}
