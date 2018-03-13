package graphql

import "github.com/contributor-ninja/infra/protocol"

type userResolver struct {
	s protocol.User
}

func (r *userResolver) Login() string {
	return r.s.Login
}

func (r *userResolver) AvatarUrl() string {
	return r.s.AvatarURL
}

func (r *userResolver) IsConnected() bool {
	return true
}
