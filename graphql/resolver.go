package graphql

import (
	"errors"
	"os"
	"strconv"

	"github.com/apex/invoke"
	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/contributor-ninja/infra/dynamodb"
	"github.com/contributor-ninja/infra/protocol"
)

var (
	invokeRegion = "us-east-1"

	id     = os.Getenv("DB_AWS_ACCESS_KEY_ID")
	secret = os.Getenv("DB_AWS_ACCESS_KEY")
	token  = ""
)

type Resolver struct {
	DynamodbClient *awsdynamodb.DynamoDB
}

var (
	defaultUser = protocol.User{
		Login: "xtuc",
	}
)

/*
  Queries
*/

func (r *Resolver) Dashboard() ([]*columnResolver, error) {
	resolvers := make([]*columnResolver, 0)

	for _, col := range protocol.DefaultDashboard.Columns {

		/*
			Fetch and decode IssueIdIndex for each Column
		*/
		query := dynamodb.MakeGetIssueIdIndexQuery(col.Id)
		resp, findQueryErr := r.DynamodbClient.Query(&query)

		if findQueryErr != nil {
			log.WithError(findQueryErr).Fatal("could not fetch IssueIdIndex")
		}

		if len(resp.Items) == 0 {
			log.
				WithField("id", col.Id).
				Info("IssueIdIndex not found for Column")

			continue
		}

		issueIdIndex := &protocol.IssueIdIndex{
			Name: col.Id,
			Ids:  make([]int, 0),
		}

		idsValueAtribute := *resp.Items[0]["ids"]

		for _, str := range idsValueAtribute.NS {
			intid, convErr := strconv.Atoi(*str)

			if convErr != nil {
				log.WithError(convErr).Fatal("could not decode items")
			}

			issueIdIndex.Ids = append(issueIdIndex.Ids, intid)

		}

		col.IssueIdIndex = issueIdIndex

		resolvers = append(resolvers, &columnResolver{
			s:              col,
			dynamodbClient: r.DynamodbClient,
		})
	}

	return resolvers, nil
}

func (r *Resolver) User() (*userResolver, error) {
	return &userResolver{defaultUser}, nil
}

/*
  Mutations
*/

type addGitHubProjectArgs struct {
	Org      string
	Name     string
	Labels   *[]*string
	Language string
}

func (r *Resolver) AddProject(args addGitHubProjectArgs) (*projectResolver, error) {
	project := protocol.MakeGitHubProject(args.Org, args.Name)

	for _, label := range *args.Labels {
		project.Labels = append(project.Labels, *label)
	}

	project.Language = args.Language

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

	/*
		Submit the indexing job
	*/
	creds := credentials.NewStaticCredentials(id, secret, token)

	client := lambda.New(session.New(&aws.Config{
		Region:      aws.String(invokeRegion),
		Credentials: creds,
	}))

	invokeErr := invoke.InvokeAsync(client, "crawler_crawler", project)

	if invokeErr != nil {
		log.WithError(invokeErr).Fatal("could not invoke crawler")
	}

	log.Info("crawler InvokeAsync")

	return &projectResolver{project}, nil
}
