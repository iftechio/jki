package factory

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"

	"github.com/iftechio/jki/pkg/configflags"
	"github.com/iftechio/jki/pkg/registry"

	"github.com/docker/docker/client"
	"k8s.io/client-go/kubernetes"
)

type Factory interface {
	genericclioptions.RESTClientGetter
	NewBuilder() *resource.Builder
	DockerClient() (*client.Client, error)
	LoadRegistries() (defReg string, registries map[string]*registry.Registry, err error)
	ToResolver() (*registry.Resolver, error)
	KubeClient() (*kubernetes.Clientset, error)
}

type factoryImpl struct {
	*configflags.ConfigFlags
}

func (f *factoryImpl) DockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

func (f *factoryImpl) NewBuilder() *resource.Builder {
	return resource.NewBuilder(f.ConfigFlags)
}

func New(cflags *configflags.ConfigFlags) Factory {
	return &factoryImpl{
		ConfigFlags: cflags,
	}
}
