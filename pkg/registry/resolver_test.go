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
		if matches[1] != tC.region {
			t.Fatalf("wrong region: expect: %s, get: %s", tC.region, matches[1])
		}
	}
}
