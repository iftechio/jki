package registry

type RegistryInterface interface {
	CreateRepoIfNotExists(repo string) error
	Prefix() string
	GetAuthToken() (string, error)
	GetLatestTag(repo string) (string, error)
	Verify() error
}
