package registry

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func LoadRegistries(configPath string) (defaultRegistry string, registries map[string]*Registry, err error) {
	type Config struct {
		Registries      []*Registry `yaml:"registries"`
		DefaultRegistry string      `yaml:"default-registry"`
	}

	f, err := os.Open(configPath)
	if err != nil {
		return "", nil, fmt.Errorf("open file: %s", err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	var config Config
	err = dec.Decode(&config)
	if err != nil {
		return "", nil, fmt.Errorf("decode yaml: %s", err)
	}

	nreg := len(config.Registries)
	if nreg == 0 {
		return "", nil, fmt.Errorf("no registries found")
	}

	defReg := config.DefaultRegistry
	if len(defReg) == 0 {
		defReg = config.Registries[0].Name
	}

	regs := make(map[string]*Registry, nreg)
	for i, reg := range config.Registries {
		if len(reg.Name) == 0 && nreg > 1 {
			return "", nil, fmt.Errorf("name of registry %d cannot be empty", i)
		}
		regs[reg.Name] = reg
	}
	if _, exist := regs[defReg]; !exist {
		return "", nil, fmt.Errorf("default registry not found: %s", defReg)
	}
	return defReg, regs, nil
}
