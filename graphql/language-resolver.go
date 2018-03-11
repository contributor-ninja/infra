package graphql

import "github.com/contributor-ninja/infra/protocol"

type languageResolver struct {
	s protocol.Language
}

func (r *languageResolver) Name() string {
	return r.s.Name
}
