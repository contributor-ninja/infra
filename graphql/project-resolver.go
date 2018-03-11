package graphql

import (
	"github.com/contributor-ninja/infra/protocol"
	gographql "github.com/graph-gophers/graphql-go"
)

type projectResolver struct {
	s protocol.GitHubProject
}

func (r *projectResolver) ID() gographql.ID {
	return gographql.ID(r.s.Id)
}

func (r *projectResolver) Org() string {
	return r.s.Org
}

func (r *projectResolver) Name() string {
	return r.s.Name
}
