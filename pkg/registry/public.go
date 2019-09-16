package registry

// PublicRegistry represents registries which does not need authentication.
type PublicRegistry struct {
}

var _ RegistryInterface = &PublicRegistry{}

func (p *PublicRegistry) CreateRepoIfNotExists(repo string) error {
	return nil
}

func (p *PublicRegistry) Prefix() string {
	return ""
}

func (p *PublicRegistry) GetAuthToken() (string, error) {
	return "", nil
}

func (p *PublicRegistry) GetLatestTag(repo string) (string, error) {
	return "latest", nil
}

func (p *PublicRegistry) Verify() error {
	return nil
}
