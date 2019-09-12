package registry

type RegistryInterface interface {
	CreateRepoIfNotExists(repo string) error
	Domain() string
	GetAuthToken() (string, error)
	GetLatestTag(repo string) (string, error)
	Verify() error
}
