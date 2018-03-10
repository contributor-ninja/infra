package dynamodb

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/contributor-ninja/infra/protocol"
)

var (
	dynamodbRegion = "us-east-1"

	id     = os.Getenv("DB_AWS_ACCESS_KEY_ID")
	secret = os.Getenv("DB_AWS_ACCESS_KEY")
	token  = ""
)

func NewClient() (*awsdynamodb.DynamoDB, error) {
	creds := credentials.NewStaticCredentials(id, secret, token)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(dynamodbRegion),
		Credentials: creds,
	})

	if err != nil {
		return nil, err
	}

	return awsdynamodb.New(sess), err
}

func MakeFindQuery(p protocol.GitHubProject) awsdynamodb.QueryInput {
	return awsdynamodb.QueryInput{
		TableName: protocol.GitHubProjectTable,

		KeyConditions: map[string]*awsdynamodb.Condition{
			"id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*awsdynamodb.AttributeValue{
					{
						S: aws.String(p.GetId()),
					},
				},
			},
		},
	}
}
