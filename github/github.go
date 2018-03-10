package github

import (
	"context"
	"os"

	gogithub "github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/contributor-ninja/infra/protocol"
)

var (
	token = os.Getenv("GITHUB_TOKEN")
)

type Client struct {
	githubClient *gogithub.Client
}

func NewClient() *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	githubClient := gogithub.NewClient(tc)

	return &Client{githubClient}
}

func (c *Client) FetchIssues(p protocol.GitHubProject) ([]protocol.Issue, error) {
	issues := make([]protocol.Issue, 0)

	opts := &gogithub.IssueListByRepoOptions{
		Labels: p.Labels,
	}

	res, _, err := c.githubClient.Issues.ListByRepo(context.Background(), p.Org, p.Name, opts)

	if err != nil {
		return issues, err
	}

	for _, item := range res {
		issues = append(issues, protocol.Issue{
			Id:        int(*item.ID),
			Title:     *item.Title,
			Body:      *item.Title,
			CreatedAt: *item.CreatedAt,
			HTMLURL:   *item.HTMLURL,

			User: protocol.User{
				Login:     *item.User.Login,
				AvatarURL: *item.User.AvatarURL,
			},
		})
	}

	return issues, nil
}
