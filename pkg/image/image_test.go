package image

import (
	"testing"
)

func TestFromString(t *testing.T) {
	testCases := []struct {
		image  string
		domain string
		repo   string
		tag    string
	}{
		{
			image:  "registry.cn-hangzhou.aliyuncs.com/foo/a-service:master-foo",
			domain: "registry.cn-hangzhou.aliyuncs.com/foo",
			repo:   "a-service",
			tag:    "master-foo",
		},
		{
			image:  "registry.cn-hangzhou.aliyuncs.com/bar/a-service",
			domain: "registry.cn-hangzhou.aliyuncs.com/bar",
			repo:   "a-service",
			tag:    "latest",
		},
		{
			image:  "123.dkr.ecr.cn-north-1.amazonaws.com.cn/nginx:latest",
			domain: "123.dkr.ecr.cn-north-1.amazonaws.com.cn",
			repo:   "nginx",
			tag:    "latest",
		},
		{
			image:  "456.dkr.ecr.cn-north-1.amazonaws.com/nginx:alpine",
			domain: "456.dkr.ecr.cn-north-1.amazonaws.com",
			repo:   "nginx",
			tag:    "alpine",
		},
		{
			image:  "nginx:alpine",
			domain: "",
			repo:   "nginx",
			tag:    "alpine",
		},
		{
			image:  "foo/nginx:alpine",
			domain: "foo",
			repo:   "nginx",
			tag:    "alpine",
		},
	}
	for _, tC := range testCases {
		img := FromString(tC.image)
		if img.Domain != tC.domain {
			t.Errorf("wrong domain, got: %s, expected: %s", img.Domain, tC.domain)
		}
		if img.Repo != tC.repo {
			t.Errorf("wrong repo, got: %s, expected: %s", img.Repo, tC.repo)
		}
		if img.Tag != tC.tag {
			t.Errorf("wrong tag, got: %s, expected: %s", img.Tag, tC.tag)
		}
	}
}
