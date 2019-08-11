package registry

import (
	"fmt"
	"regexp"
)

var (
	reAWSECR   = regexp.MustCompile(`\d+\.dkr\.ecr\.(?P<Region>[\w-]+)\.amazonaws\.com`)
	reAliCloud = regexp.MustCompile(`registry\.(?P<Region>[\w-]+)\.aliyuncs.com`)
	//reAliCloudVpc      = regexp.MustCompile(`registry-vpc\.[\w-]+\.aliyuncs.com`)
	//reAliCloudInternal = regexp.MustCompile(`registry-internal\.[\w-]+\.aliyuncs.com`)
)

const (
	RegistryAWS      = "AWS"
	RegistryAliCloud = "AliCloud"
)

type Resolver struct {
	registries      map[string]*Registry
	defaultRegistry string
}

func (r *Resolver) ResolveName(name string) (regKind string, registry *Registry, err error) {
	reg, exist := r.registries[name]
	if !exist {
		return "", nil, fmt.Errorf("not found: %s", name)
	}

	switch {
	case reg.AliCloud != nil:
		return RegistryAliCloud, reg, nil
	case reg.AWS != nil:
		return RegistryAWS, reg, nil
	default:
		return "", reg, nil
	}
}

func (r *Resolver) ResolveRegistryAuth(img string) (authToken string, err error) {
	if matches := reAWSECR.FindStringSubmatch(img); matches != nil {
		region := matches[1]
		for _, reg := range r.registries {
			if reg.AWS == nil {
				continue
			}
			if reg.AWS.Region == region {
				return reg.AWS.GetAuthToken()
			}
		}
	} else if matches := reAliCloud.FindStringSubmatch(img); matches != nil {
		region := matches[1]
		for _, reg := range r.registries {
			if reg.AliCloud == nil {
				continue
			}
			if reg.AliCloud.Region == region {
				return reg.AliCloud.GetAuthToken()
			}
		}
	}
	// may be public image
	return "", nil
}

func NewResolver(configPath string) (*Resolver, error) {
	defReg, regs, err := LoadRegistries(configPath)
	if err != nil {
		return nil, err
	}
	r := Resolver{
		defaultRegistry: defReg,
		registries:      regs,
	}
	return &r, nil
}
