module github.com/iftechio/jki

go 1.13

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200309214505-aa6a9891b09c

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305

replace github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200227195959-4d242818bf55

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1160
	github.com/aws/aws-sdk-go v1.29.23
	github.com/containerd/console v0.0.0-20191219165238-8375c3424e4d
	github.com/docker/docker v1.14.0-0.20190319215453-e7b5f7dbe98c
	github.com/moby/buildkit v0.7.0-rc1.0.20200312194508-a1bf12f80604
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/tonistiigi/fsutil v0.0.0-20200225063759-013a9fe6aee2
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	google.golang.org/grpc v1.27.1
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/cli-runtime v0.18.2
	k8s.io/client-go v0.18.2
	sigs.k8s.io/yaml v1.2.0
)
