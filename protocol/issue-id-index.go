package protocol

import "github.com/aws/aws-sdk-go/aws"

var (
	IssueIdIndexTable = aws.String("IssueIndices")
)

type IssueIdIndex struct {
	Name string `json:"name"` // primary key
	Ids  []int  `json:"ids"`
}

func MakeEmptyIssueIdIndex(name string) IssueIdIndex {
	return IssueIdIndex{
		Name: name,
		Ids:  make([]int, 0),
	}
}
