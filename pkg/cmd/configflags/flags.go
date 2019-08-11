package configflags

import (
	"fmt"
	"os"
	"path"

	"github.com/iftechio/jki/pkg/registry"
	"github.com/spf13/pflag"
)

type ConfigFlags struct {
	configPath string
	registry   string
}

func (f *ConfigFlags) ToResolver() (*registry.Resolver, error) {
	return registry.NewResolver(f.configPath)
}

func (f *ConfigFlags) LoadRegistries() (defReg string, registries map[string]*registry.Registry, err error) {
	defReg, registries, err = registry.LoadRegistries(f.configPath)
	if len(f.registry) != 0 {
		if _, exist := registries[f.registry]; !exist {
			return "", nil, fmt.Errorf("registry not found: %s", f.registry)
		}
		defReg = f.registry
	}
	return
}

func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
	homedir := os.Getenv("HOME")
	flags.StringVar(&f.configPath, "jkiconfig", path.Join(homedir, ".jki.yaml"), "Config path")
	flags.StringVar(&f.registry, "registry", "", "The desired registry. If not set, use the `default-registry` in config.")
}

func New() *ConfigFlags {
	return &ConfigFlags{}
}
