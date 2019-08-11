package registry

import "testing"

func TestRegexp(t *testing.T) {
	t.Parallel()
	imgs := []string{
		"131061968068.dkr.ecr.ap-northeast-1.amazonaws.com/api-gateway-controller:master-636e3c1",
		"804775010343.dkr.ecr.cn-north-1.amazonaws.com.cn/anti-spider:master-2d9b02f",
	}
	regions := []string{"ap-northeast-1", "cn-north-1"}
	for i, img := range imgs {
		matches := reAWSECR.FindStringSubmatch(img)
		if matches == nil {
			t.Fatalf("no matches")
		}
		if matches[1] != regions[i] {
			t.Fatalf("wrong region: expect: %s, get: %s", regions[i], matches[1])
		}
	}

	imgs = []string{
		"registry.cn-hangzhou.aliyuncs.com/iftech/api-gateway-controller:master-636e3c1",
	}
	regions = []string{"cn-hangzhou"}
	for i, img := range imgs {
		matches := reAliCloud.FindStringSubmatch(img)
		if matches == nil {
			t.Fatalf("no matches")
		}
		if matches[1] != regions[i] {
			t.Fatalf("wrong region: expect: %s, get: %s", regions[i], matches[1])
		}
	}
}
