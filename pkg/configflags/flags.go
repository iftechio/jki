package configflags

import (
	"fmt"
	"path/filepath"

	"github.com/iftechio/jki/pkg/registry"
	"github.com/iftechio/jki/pkg/utils"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ConfigFlags struct {
	configPath string
	registry   string
	kubeconfig string
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

func (f *ConfigFlags) KubeClient() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", f.kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
	homedir := utils.HomeDir()
	flags.StringVar(&f.configPath, "jkiconfig", filepath.Join(homedir, ".jki.yaml"), "Config path")
	flags.StringVarP(&f.registry, "registry", "r", "", "The desired registry. If not set, use the `default-registry` in config.")
	flags.StringVarP(&f.kubeconfig, "kubeconfig", "", filepath.Join(homedir, ".kube", "config"), "The path to kubeconfig. If not set `~/.kube/config` will be used")
}

func New() *ConfigFlags {
	return &ConfigFlags{}
}
