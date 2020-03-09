# Example of using local DynamoDB with GO

## Run DynamoDB localy using Docker

- Run docker container:

`docker run -p 8000:8000 amazon/dynamodb-local`

- Create a table:

`
aws dynamodb create-table \
    --table-name user_variations \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url http://localhost:8000
` 

- Create a GSI:


`
aws dynamodb update-table \
    --endpoint-url http://localhost:8000 \
    --table-name user_variations \
    --attribute-definitions \
        AttributeName=user_id,AttributeType=S \
        AttributeName=campaign_id,AttributeType=S \
    --global-secondary-index-updates \
    "[{\"Create\":{\"IndexName\": \"user_id_campaign_id-GSI\",\"KeySchema\":\
    [{\"AttributeName\":\"user_id\",\"KeyType\":\"HASH\"}, {\"AttributeName\":\"campaign_id\",\"KeyType\":\"RANGE\"}], \
    \"ProvisionedThroughput\": {\"ReadCapacityUnits\": 5, \"WriteCapacityUnits\": 5},\"Projection\":{\"ProjectionType\":\"ALL\"}}}]"
`

- Insert an item:


`
aws dynamodb put-item \
    --endpoint-url http://localhost:8000 \
    --table-name user_variations \
    --item '{"id": {"S": "1"}, "user_id": {"S": "1"}, "campaign_id": {"S": "1"}, "end_time": {"N": "10"}}'
`
