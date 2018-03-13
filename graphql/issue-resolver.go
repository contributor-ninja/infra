package graphql

import "github.com/contributor-ninja/infra/protocol"

type issueResolver struct {
	s protocol.Issue
}

func (r *issueResolver) Title() string {
	return r.s.Title
}

func (r *issueResolver) Body() string {
	return r.s.Body
}

func (r *issueResolver) HtmlUrl() string {
	return r.s.HTMLURL
}

func (r *issueResolver) AvatarUrl() string {
	return "https://avatars2.githubusercontent.com/u/1253363?s=200"
}

func (r *issueResolver) Project() *projectResolver {
	return &projectResolver{r.s.Project}
}

func (r *issueResolver) User() *userResolver {
	return &userResolver{r.s.User}
}
