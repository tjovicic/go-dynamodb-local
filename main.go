package main

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func queryUsingGSI(tableName, indexName, keyCondExpr, filterExpr string,
	exprAttrNames map[string]*string, exprAttrValue map[string]*dynamodb.AttributeValue) ([]map[string]*dynamodb.AttributeValue, error) {
	limit := int64(1)
	falseBool := false
	queryInput := dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(indexName),
		KeyConditionExpression:    aws.String(keyCondExpr),
		FilterExpression:          aws.String(filterExpr),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValue,
		Limit:                     &limit,
		ScanIndexForward:          &falseBool,
	}

	cfg := &aws.Config{
		Endpoint: aws.String("http://localhost:8000"),
		Region:   aws.String("eu-central-1"), // this should be same sa in the ~/.aws/config
	}
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess, cfg)

	result, err := db.Query(&queryInput)
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

func main() {
	expressionAttributeNames := map[string]*string{
		"#user_id":     aws.String("user_id"),
		"#campaign_id": aws.String("campaign_id"),
		"#endTime":     aws.String("end_time"),
	}

	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
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

	fmt.Println(queryUsingGSI(
		"user_variations", "user_id_campaign_id-GSI", "#user_id = :userId AND #campaign_id = :campaignId", "#endTime > :endTime",
		expressionAttributeNames, expressionAttributeValues,
	))
}
