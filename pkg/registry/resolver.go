package registry

type Resolver struct {
	registries      map[string]*Registry
	defaultRegistry string
}

func (r *Resolver) ResolveRegistryByImage(img string) (Interface, error) {
	for _, reg := range r.registries {
		if reg.MatchImage(img) {
			return reg, nil
		}
	}
	// may be public image
	return &Registry{}, nil
}

func NewResolver(configPath string) (*Resolver, error) {
	defReg, regs, err := LoadRegistries(configPath)
	if err != nil {
		return nil, err
	}
	r := Resolver{
		defaultRegistry: defReg,
		registries:      regs,
	}
	return &r, nil
}
