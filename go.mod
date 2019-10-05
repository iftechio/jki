module github.com/iftechio/jki

go 1.12

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190822205725-ed20165a37b4

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305

require (
	github.com/aliyun/alibaba-cloud-sdk-go v0.0.0-20190929091402-5711055976b5
	github.com/aws/aws-sdk-go v1.25.6
	github.com/containerd/console v0.0.0-20181022165439-0650fd9eeb50
	github.com/docker/docker v1.14.0-0.20190319215453-e7b5f7dbe98c
	github.com/moby/buildkit v0.6.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/api v0.0.0-20191005115622-2e41325d9e4b
	k8s.io/apimachinery v0.0.0-20191005115455-e71eb83a557c
	k8s.io/client-go v12.0.0+incompatible
)
