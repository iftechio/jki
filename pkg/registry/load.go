package registry

import (
	"fmt"
	"io/ioutil"
	"os"

	"sigs.k8s.io/yaml"
)

func LoadRegistries(configPath string) (defaultRegistry string, registries map[string]*Registry, err error) {
	type Config struct {
		Registries      []*Registry `json:"registries"`
		DefaultRegistry string      `json:"default-registry"`
	}

	f, err := os.Open(configPath)
	if err != nil {
		return "", nil, fmt.Errorf("open file: %s", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return "", nil, fmt.Errorf("decode yaml: %s", err)
	}

	nReg := len(config.Registries)
	if nReg == 0 {
		return "", nil, fmt.Errorf("no registries found")
	}

	defReg := config.DefaultRegistry
	if len(defReg) == 0 {
		defReg = config.Registries[0].Name
	}

	regs := make(map[string]*Registry, nReg)
	for i, reg := range config.Registries {
		if len(reg.Name) == 0 && nReg > 1 {
			return "", nil, fmt.Errorf("name of registry %d cannot be empty", i)
		}
		regs[reg.Name] = reg
	}
	if _, exist := regs[defReg]; !exist {
		return "", nil, fmt.Errorf("default registry not found: %s", defReg)
	}
	return defReg, regs, nil
}
