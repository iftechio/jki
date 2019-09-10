package registry

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type AWSRegistry struct {
	Region              string `json:"region" yaml:"region"`
	AccountID           string `json:"account_id" yaml:"account_id"`
	AccessKey           string `json:"access_key" yaml:"access_key"`
	SecretAccessKey     string `json:"secret_access_key" yaml:"secret_access_key"`
	LifecyclePolicyText string `json:"lifecycle_policy_text" yaml:"lifecycle_policy_text"`
}

var _ RegistryInterface = &AWSRegistry{}

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

func (r *AWSRegistry) Domain() string {
	domain := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", r.AccountID, r.Region)
	if strings.HasPrefix(r.Region, "cn-") {
		return domain + ".cn"
	}
	return domain
}

func (r *AWSRegistry) GetAuthToken() (string, error) {
	sess, err := newAWSSession(r.Region, r.AccessKey, r.SecretAccessKey)
	if err != nil {
		return "", err
	}

	ecrSvc := ecr.New(sess)
	output, err := ecrSvc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return "", err
	}
	if len(output.AuthorizationData) < 1 {
		return "", fmt.Errorf("missing token")
	}
	token := output.AuthorizationData[0]
	encodedToken := aws.StringValue(token.AuthorizationToken)
	data, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		return "", fmt.Errorf("decode ecr token: %s", err)
	}
	parts := strings.Split(string(data), ":")
	return toRegistryAuth(parts[0], parts[1])
}
