package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/contributor-ninja/infra/dynamodb"
	"github.com/contributor-ninja/infra/protocol"

	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	yaml "gopkg.in/yaml.v2"
)

type YAMLProject struct {
	Repo  string `yaml:"repo"`
	Label string `yaml:"label"`
}

func main() {
	svc, dynamoerr := dynamodb.NewClient()

	if dynamoerr != nil {
		panic(dynamoerr)
	}

	files, err := ioutil.ReadDir("./projects")

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if path.Ext(file.Name()) == ".yml" {
			fmt.Printf("----- %s -----\n", file.Name())
			data, err := ioutil.ReadFile("./projects/" + file.Name())

			if err != nil {
				panic(err)
			}

			p := make([]YAMLProject, 0)

			if err := yaml.Unmarshal([]byte(data), &p); err != nil {
				panic(err)
			}

			lang := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			insert(svc, p, lang)
		}
	}
}

func insert(svc *awsdynamodb.DynamoDB, projects []YAMLProject, lang string) {

	for _, yamlproject := range projects {
		splitedRepo := strings.Split(yamlproject.Repo, "/")
		org := splitedRepo[0]
		name := splitedRepo[1]

		project := protocol.MakeGitHubProject(org, name)
		project.Language = lang

		if !project.HasLabel(yamlproject.Label) {
			project.Labels = append(project.Labels, yamlproject.Label)
		}

		inputQuery := dynamodb.MakePutGitHubProject(project)
		_, putErr := svc.PutItem(&inputQuery)

		if putErr != nil {
			panic(putErr)
		}
	}

}
