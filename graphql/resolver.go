package graphql

import (
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/contributor-ninja/infra/protocol"
)

type Resolver struct {
	DynamodbClient *awsdynamodb.DynamoDB
}

var (
	defaultDashboard = protocol.Dashboard{
		Columns: []protocol.Column{
			{
				Id:       "default-php",
				Language: protocol.Language{"php"},
				IssueIdIndex: protocol.IssueIdIndex{
					Ids: []int{291963812},
				},
			},

			{
				Id:       "default-js",
				Language: protocol.Language{"js"},
				IssueIdIndex: protocol.IssueIdIndex{
					Ids: []int{294631994, 270481545, 295238273},
				},
			},

			{
				Id:       "default-html",
				Language: protocol.Language{"html"},
				IssueIdIndex: protocol.IssueIdIndex{
					Ids: []int{301986252, 270481580},
				},
			},

			{
				Id:       "default-ruby",
				Language: protocol.Language{"ruby"},
				IssueIdIndex: protocol.IssueIdIndex{
					Ids: []int{304049912},
				},
			},
		},
	}
)

func (r *Resolver) Dashboard() ([]*columnResolver, error) {
	resolvers := make([]*columnResolver, 0)

	for _, col := range defaultDashboard.Columns {
		resolvers = append(resolvers, &columnResolver{
			s:              col,
			dynamodbClient: r.DynamodbClient,
		})
	}

	return resolvers, nil
}

func (r *Resolver) AddProject() *string {
	return nil
}
