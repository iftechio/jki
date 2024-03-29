package configflags

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/iftechio/jki/pkg/registry"
	"github.com/iftechio/jki/pkg/utils"
)

type ConfigFlags struct {
	configPath  string
	registry    string
	platform    string
	kubeconfig  string
	namespace   string
	konfigFlags *genericclioptions.ConfigFlags
}

func (f *ConfigFlags) ToRESTConfig() (*rest.Config, error) {
	return f.konfigFlags.ToRESTConfig()
}

func (f *ConfigFlags) Platform() string {
	if f.platform == "" {
		return runtime.GOARCH
	}
	return f.platform
}

func (f *ConfigFlags) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	return f.konfigFlags.ToDiscoveryClient()
}

func (f *ConfigFlags) ToRESTMapper() (meta.RESTMapper, error) {
	return f.konfigFlags.ToRESTMapper()
}

func (f *ConfigFlags) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return f.konfigFlags.ToRawKubeConfigLoader()
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
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func (f *ConfigFlags) ConfigPath() string {
	return f.configPath
}

func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
	homedir := utils.HomeDir()
	flags.StringVar(&f.configPath, "jkiconfig", filepath.Join(homedir, ".jki.yaml"), "Config path")
	flags.StringVarP(&f.registry, "registry", "r", "", "The desired registry. If not set, use the `default-registry` in config.")
	flags.StringVarP(&f.platform, "platform", "p", "", fmt.Sprintf("The desired platform. (default \"%s\")", runtime.GOARCH))
	flags.StringVar(f.konfigFlags.KubeConfig, "kubeconfig", filepath.Join(homedir, ".kube", "config"), "The path to kubeconfig. If not set `~/.kube/config` will be used")
	flags.StringVarP(f.konfigFlags.Namespace, "namespace", "n", "", "If present, the namespace scope for this CLI request")
}

func New() *ConfigFlags {
	return &ConfigFlags{
		konfigFlags: genericclioptions.NewConfigFlags(true),
	}
}
