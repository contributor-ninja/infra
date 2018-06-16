package protocol

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
)

var (
	GitHubProjectTable = aws.String("GitHubProjects")
)

const (
	LANGUAGE_JS = "js"
)

type GitHubProject struct {
	Id       string   `json:"id"` // composition of org/name
	Org      string   `json:"org"`
	Name     string   `json:"name"`
	Labels   []string `json:"labels"`
	Language string   `json:"language"`
}

func MakeGitHubProject(org, name string) GitHubProject {
	id := fmt.Sprintf("%s/%s", org, name)

	labels := []string{
		"good first issue", // default GitHub label
	}

	return GitHubProject{
		Id: id,

		Org:  org,
		Name: name,

		Labels:   labels,
		Language: LANGUAGE_JS,
	}
}

func (i GitHubProject) GetId() string {
	return fmt.Sprintf("%s/%s", i.Org, i.Name)
}

func (i GitHubProject) HasLabel(label string) bool {
	for _, l := range i.Labels {
		if l == label {
			return true
		}
	}

	return false
}
