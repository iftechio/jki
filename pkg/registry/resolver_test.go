package registry

import (
	"testing"
)

func TestResolveImage(t *testing.T) {
	t.Parallel()
	rs := Resolver{
		registries: map[string]*Registry{
			"r1": {
				AWS: &AWSRegistry{
					Region:    "cn-north-1",
					AccountID: "1234",
				},
			},
			"r2": {
				AWS: &AWSRegistry{
					Region:    "cn-north-1",
					AccountID: "4321",
				},
			},
			"r3": {
				AWS: &AWSRegistry{
					Region:    "us-east-1",
					AccountID: "123",
				},
			},
			"a1": {
				AliCloud: &AliCloudRegistry{
					Region:    "cn-hangzhou",
					Namespace: "ns1",
				},
			},
			"a2": {
				AliCloud: &AliCloudRegistry{
					Region:    "cn-hangzhou",
					Namespace: "ns-2",
				},
			},
			"a3": {
				AliCloud: &AliCloudRegistry{
					Region:    "cn-shanghai",
					Namespace: "ns_3",
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
		{
			image:  "registry.cn-hangzhou.aliyuncs.com/ns-2/busybox:latest",
			prefix: "registry.cn-hangzhou.aliyuncs.com/ns-2",
		},
		{
			image:  "registry.cn-shanghai.aliyuncs.com/ns_3/busybox:latest",
			prefix: "registry.cn-shanghai.aliyuncs.com/ns_3",
		},
		{
			image:  "registry.cn-hangzhou.aliyuncs.com/ns1/busybox:latest",
			prefix: "registry.cn-hangzhou.aliyuncs.com/ns1",
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
