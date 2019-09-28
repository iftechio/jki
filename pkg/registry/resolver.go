package registry

import (
	"regexp"
)

var (
	reAWSECR   = regexp.MustCompile(`(?P<AccountID>\d+)\.dkr\.ecr\.(?P<Region>[\w-]+)\.amazonaws\.com`)
	reAliCloud = regexp.MustCompile(`registry\.(?P<Region>[\w-]+)\.aliyuncs.com`)
)

type Resolver struct {
	registries      map[string]*Registry
	defaultRegistry string
}

func (r *Resolver) ResolveRegistryByImage(img string) (RegistryInterface, error) {
	if matches := reAWSECR.FindStringSubmatch(img); matches != nil {
		accountID := matches[1]
		region := matches[2]
		for _, reg := range r.registries {
			if reg.AWS == nil {
				continue
			}
			if reg.AWS.Region == region && reg.AWS.AccountID == accountID {
				return reg, nil
			}
		}
	} else if matches := reAliCloud.FindStringSubmatch(img); matches != nil {
		region := matches[1]
		for _, reg := range r.registries {
			if reg.AliCloud == nil {
				continue
			}
			if reg.AliCloud.Region == region {
				return reg, nil
			}
		}
	}
	// may be public image
	return &Registry{}, nil
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
