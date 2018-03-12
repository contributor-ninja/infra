package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/apex/log"
	loghandlers "github.com/apex/log/handlers/json"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	muxhandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	gographql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

	"github.com/contributor-ninja/infra/api"
	"github.com/contributor-ninja/infra/dynamodb"
	"github.com/contributor-ninja/infra/graphql"
)

var (
	port       = os.Getenv("PORT")
	corsOrigin = os.Getenv("GRAPHQL_CORS_ORIGIN")

	page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.css" rel="stylesheet" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/es6-promise/4.1.1/es6-promise.auto.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/query", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}

			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)
)

func init() {
	log.SetHandler(loghandlers.Default)
}

func main() {
	addr := ":" + port

	r := mux.NewRouter()

	s, schemaErr := graphql.GetSchema("./schema.graphql")
	if schemaErr != nil {
		log.WithError(schemaErr).Fatal("could not get GraphQL schema")
	}

	svc, err := dynamodb.NewClient()

	if err != nil {
		log.WithError(err).Fatal("connection to dynamodb failed")
	}

	schema := gographql.MustParseSchema(s, &graphql.Resolver{svc})
	handlers := Handlers{svc}

	/*
		Handlers
	*/

	r.HandleFunc("/status", handlers.getStatusHandler)

	r.Handle("/query", &relay.Handler{Schema: schema})

	r.Handle("/graphiql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	originsOk := muxhandlers.AllowedOrigins([]string{corsOrigin})
	methodsOk := muxhandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	routerWithCors := muxhandlers.CORS(originsOk, methodsOk)(r)

	if err := http.ListenAndServe(addr, routerWithCors); err != nil {
		log.WithError(err).Fatal("error listening")
	}
}

/*
	Handler methods
*/
type Handlers struct {
	dynamodbClient *awsdynamodb.DynamoDB
}

func (h Handlers) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	_, listTablesErr := h.dynamodbClient.ListTables(&awsdynamodb.ListTablesInput{})

	if listTablesErr != nil {
		log.WithError(listTablesErr).Fatal("connection to dynamodb failed")
	}

	res := api.Response{"ok"}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		log.WithError(err).Fatal("could not send response")
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(jsonBytes); err != nil {
		log.WithError(err).Fatal("could not send response")
	}
}
