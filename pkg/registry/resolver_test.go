package registry

import (
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		re     *regexp.Regexp
		image  string
		region string
	}{
		{
			re:     reAWSECR,
			image:  "1234.dkr.ecr.ap-northeast-1.amazonaws.com/foo1:master-636e3c1",
			region: "ap-northeast-1",
		},
		{
			re:     reAWSECR,
			image:  "4321.dkr.ecr.cn-north-1.amazonaws.com.cn/foo2:master-2d9b02f",
			region: "cn-north-1",
		},
		{
			re:     reAliCloud,
			image:  "registry.cn-hangzhou.aliyuncs.com/bar/foo1:master-636e3c1",
			region: "cn-hangzhou",
		},
	}
	for _, tC := range testCases {
		matches := tC.re.FindStringSubmatch(tC.image)
		if matches == nil {
			t.Fatalf("no matches")
		}
		var gotRegion string
		if tC.re == reAWSECR {
			gotRegion = matches[2]
		} else if tC.re == reAliCloud {
			gotRegion = matches[1]
		}
		if gotRegion != tC.region {
			t.Fatalf("wrong region: expect: %s, get: %s", tC.region, gotRegion)
		}
	}
}

func TestResolveImage(t *testing.T) {
	t.Parallel()
	rs := Resolver{
		registries: map[string]*Registry{
			"r1": &Registry{
				AWS: &AWSRegistry{
					Region:    "cn-north-1",
					AccountID: "1234",
				},
			},
			"r2": &Registry{
				AWS: &AWSRegistry{
					Region:    "cn-north-1",
					AccountID: "4321",
				},
			},
			"r3": &Registry{
				AWS: &AWSRegistry{
					Region:    "us-east-1",
					AccountID: "123",
				},
			},
		},
	}

	testCases := []struct {
		image  string
		prefix string
	}{
		{
			image:  "1234.dkr.ecr.cn-north-1.amazonaws.com.cn/foo1:dev",
			prefix: "1234.dkr.ecr.cn-north-1.amazonaws.com.cn",
		},
		{
			image:  "4321.dkr.ecr.cn-north-1.amazonaws.com.cn/foo1:latest",
			prefix: "4321.dkr.ecr.cn-north-1.amazonaws.com.cn",
		},
		{
			image:  "123.dkr.ecr.us-east-1.amazonaws.com/foo1:dev",
			prefix: "123.dkr.ecr.us-east-1.amazonaws.com",
		},
	}
	for _, tC := range testCases {
		ri, err := rs.ResolveRegistryByImage(tC.image)
		if err != nil {
			t.Fatal(err)
		}
		if ri.Prefix() != tC.prefix {
			t.Fatalf("got: %s, expected: %s", ri.Prefix(), tC.prefix)
		}
	}

}
