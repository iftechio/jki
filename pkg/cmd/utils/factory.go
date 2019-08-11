package utils

import (
	"github.com/iftechio/jki/pkg/cmd/configflags"
	"github.com/iftechio/jki/pkg/registry"

	"github.com/docker/docker/client"
)

type Factory interface {
	DockerClient() (*client.Client, error)
	LoadRegistries() (defReg string, registries map[string]*registry.Registry, err error)
	ToResolver() (*registry.Resolver, error)
}

type factoryImpl struct {
	configFlags *configflags.ConfigFlags
}

func (f *factoryImpl) DockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

func (f *factoryImpl) LoadRegistries() (defReg string, registries map[string]*registry.Registry, err error) {
	return f.configFlags.LoadRegistries()
}

func (f *factoryImpl) ToResolver() (*registry.Resolver, error) {
	return f.configFlags.ToResolver()
}

func NewFactory(cflags *configflags.ConfigFlags) Factory {
	return &factoryImpl{cflags}
}
