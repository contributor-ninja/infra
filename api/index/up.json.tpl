{
  "name": "index",
  "profile": "contributorninja-infra",
  "regions": [
    "us-east-1"
  ],
  "error_pages": {
    "disable": true
  },
  "proxy": {
    "attempts": 1
  },
  "stages": {
    "production": {
      "domain": "apiv2-index.contributor.ninja"
    },
    "staging": {
      "domain": "staging-apiv2-index.contributor.ninja"
    }
  },
  "cors": {
    "enable": true,
    "allowed_origins": [
      "https://contributor.ninja",
      "https://apiv2-index.contributor.ninja",
      "https://staging-apiv2-index.contributor.ninja",
      "http://localhost:3000"
    ],
    "allowed_methods": ["HEAD", "GET", "POST", "PUT", "OPTIONS"],
    "allow_credentials": true,
    "allowed_headers": ["*"]
  },
  "environment": {
    "DB_AWS_ACCESS_KEY_ID": "${DB_AWS_ACCESS_KEY_ID}",
    "DB_AWS_ACCESS_KEY": "${DB_AWS_ACCESS_KEY}",
    "GRAPHQL_CORS_ORIGIN": "https://contributor.ninja"
  }
}
