package graphql

import (
	"errors"

	"github.com/apex/log"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/contributor-ninja/infra/dynamodb"
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

type addGitHubProjectArgs struct {
	Org  string
	Name string
}

func (r *Resolver) AddProject(args addGitHubProjectArgs) (*projectResolver, error) {
	project := protocol.MakeGitHubProject(args.Org, args.Name)

	query := dynamodb.MakeFindQuery(project)
	resp, findQueryErr := r.DynamodbClient.Query(&query)

	if findQueryErr != nil {
		return nil, errors.New("could not find project")
	}

	if len(resp.Items) > 0 {
		return nil, errors.New("project already exists")
	}

	/*
		Insert new project
	*/
	av, err := dynamodbattribute.MarshalMap(project)

	if err != nil {
		log.WithError(err).Fatal("could not MarshalMap")
	}

	input := &awsdynamodb.PutItemInput{
		Item:      av,
		TableName: protocol.GitHubProjectTable,
	}

	_, putErr := r.DynamodbClient.PutItem(input)

	if putErr != nil {
		log.WithError(putErr).Fatal("could not send response")
	}

	log.
		WithFields(log.Fields{
			"id": project.Id,
		}).
		Info("added item in index")

	return &projectResolver{project}, nil
}
