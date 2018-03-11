package graphql

import (
	"strconv"

	"github.com/apex/log"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	gographql "github.com/graph-gophers/graphql-go"

	"github.com/contributor-ninja/infra/dynamodb"
	"github.com/contributor-ninja/infra/protocol"
)

type columnResolver struct {
	dynamodbClient *awsdynamodb.DynamoDB
	s              protocol.Column
}

func (r *columnResolver) ID() gographql.ID {
	return gographql.ID(r.s.Id)
}

func (r *columnResolver) Language() *languageResolver {
	return &languageResolver{r.s.Language}
}

// Fetch the issues given the column name
func (r *columnResolver) Issues() []*issueResolver {
	resolvers := make([]*issueResolver, 0)

	batchOpts := dynamodb.MakeBatchGetItemByIssueIndex(r.s.IssueIdIndex)

	resp, getErr := r.dynamodbClient.BatchGetItem(&batchOpts)

	if getErr != nil {
		log.WithError(getErr).Fatal("could not get issues")
	}

	issueTableName := *protocol.IssueTable
	tableResp := resp.Responses[issueTableName]

	for _, item := range tableResp {
		id, err := strconv.Atoi(*item["id"].N)

		if err != nil {
			log.WithError(err).Fatal("could not decode response")
		}

		issue := protocol.Issue{
			Id:      id,
			Title:   *item["title"].S,
			Body:    *item["body"].S,
			HTMLURL: *item["html_url"].S,

			Project: protocol.GitHubProject{
				Name: "test",
				Org:  "test",
			},

			// FIXME(sven): can't query user, it's a reserved work
			// TODO(sven): rename user in the model
			User: protocol.User{
				Login:     "xtuc",
				AvatarURL: "foo.png",
			},
		}

		resolvers = append(resolvers, &issueResolver{issue})
	}

	return resolvers
}
