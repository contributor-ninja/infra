{
  "name": "crawler",
  "description": "",
  "memory": 128,
  "timeout": 5,
  "role": "arn:aws:iam::161663720085:role/crawler_lambda_function",
  "environment": {
    "DB_AWS_ACCESS_KEY_ID": "${DB_AWS_ACCESS_KEY_ID}",
    "DB_AWS_ACCESS_KEY": "${DB_AWS_ACCESS_KEY}",
    "GITHUB_TOKEN": "${GITHUB_TOKEN}",
  }
}
