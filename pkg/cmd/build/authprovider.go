package build

import (
	"sync"

	"github.com/moby/buildkit/session/auth"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/iftechio/jki/pkg/registry"
)

type AuthProvider struct {
	registries map[string]*registry.Registry
	cache      sync.Map // for concurrent use
}

func NewAuthProvider(registries map[string]*registry.Registry) *AuthProvider {
	ap := &AuthProvider{
		registries: make(map[string]*registry.Registry, len(registries)),
	}
	for _, reg := range registries {
		host := reg.Host()
		ap.registries[host] = reg
	}
	return ap
}

func (ap *AuthProvider) Credentials(ctx context.Context, req *auth.CredentialsRequest) (*auth.CredentialsResponse, error) {
	reg, ok := ap.registries[req.Host]
	if !ok {
		// may be public image
		return &auth.CredentialsResponse{}, nil
	}
	cache, ok := ap.cache.Load(req.Host)
	if ok {
		return cache.(*auth.CredentialsResponse), nil
	}
	c, err := reg.GetAuthConfig()
	if err != nil {
		return nil, err
	}
	resp := &auth.CredentialsResponse{
		Username: c.Username,
		Secret:   c.Password,
	}
	ap.cache.Store(req.Host, resp)
	return resp, nil
}

func (ap *AuthProvider) Register(s *grpc.Server) {
	auth.RegisterAuthServer(s, ap)
}
