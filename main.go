package main

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type UserVariation struct {
	ID         string `dynamodbav:"id"`
	UserID     string `dynamodbav:"user_id"`
	CampaignID string `dynamodbav:"campaign_id"`
	EndTime    uint   `dynamodbav:"end_time"`
}

type DynamoDBInstance struct {
	ctx context.Context
	db  *dynamodb.DynamoDB
}

func NewDynamoDBInstance(ctx context.Context, endpoint, region string) *DynamoDBInstance {
	cfg := &aws.Config{
		Endpoint: aws.String(endpoint),
		Region:   aws.String(region), // this should be same as in the ~/.aws/config
	}

	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess, cfg)

	return &DynamoDBInstance{
		ctx: ctx,
		db:  db,
	}
}

// Example:
// result, _ := queryUsingGSI(repo, "user_variations", "user_id_campaign_id-GSI", "#user_id = :userId AND #campaign_id = :campaignId", "#endTime > :endTime")
func queryUsingGSI(instance *DynamoDBInstance, tableName, indexName, keyCondExpr, filterExpr string) ([]map[string]*dynamodb.AttributeValue, error) {
	limit := int64(1)
	falseBool := false

	exprAttrNames := map[string]*string{
		"#user_id":     aws.String("user_id"),
		"#campaign_id": aws.String("campaign_id"),
		"#endTime":     aws.String("end_time"),
	}

	exprAttrValues := map[string]*dynamodb.AttributeValue{
		":userId": {
			S: aws.String("1"),
		},
		":campaignId": {
			S: aws.String("1"),
		},
		":endTime": {
			N: aws.String(strconv.Itoa(0)),
		},
	}

	queryInput := dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(indexName),
		KeyConditionExpression:    aws.String(keyCondExpr),
		FilterExpression:          aws.String(filterExpr),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
		Limit:                     &limit,
		ScanIndexForward:          &falseBool,
	}

	result, err := instance.db.Query(&queryInput)
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

func batchWrite(instance *DynamoDBInstance, tableName string, item map[string]*dynamodb.AttributeValue) (*dynamodb.BatchWriteItemOutput, error) {
	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			tableName: {
				{
					PutRequest: &dynamodb.PutRequest{
						Item: item,
					},
				},
			},
		},
	}

	return instance.db.BatchWriteItem(input)
}

func main() {
	// instance := NewDynamoDBInstance(context.Background(), "http://localhost:8000", "eu-central-1")

}
