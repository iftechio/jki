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
			image:  "131061968068.dkr.ecr.ap-northeast-1.amazonaws.com/api-gateway-controller:master-636e3c1",
			region: "ap-northeast-1",
		},
		{
			re:     reAWSECR,
			image:  "804775010343.dkr.ecr.cn-north-1.amazonaws.com.cn/anti-spider:master-2d9b02f",
			region: "cn-north-1",
		},
		{
			re:     reAliCloud,
			image:  "registry.cn-hangzhou.aliyuncs.com/iftech/api-gateway-controller:master-636e3c1",
			region: "cn-hangzhou",
		},
	}
	for _, tC := range testCases {
		matches := tC.re.FindStringSubmatch(tC.image)
		if matches == nil {
			t.Fatalf("no matches")
		}
		if matches[1] != tC.region {
			t.Fatalf("wrong region: expect: %s, get: %s", tC.region, matches[1])
		}
	}
}
