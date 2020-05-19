package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func TestBatchWrite(t *testing.T) {
	t.Run("batch write", func(t *testing.T) {
		uv := UserVariation{
			ID:         "1",
			CampaignID: "1",
			UserID:     "1",
			EndTime:    1584023280,
		}

		item, err := dynamodbattribute.MarshalMap(uv)
		if err != nil {
			t.Fatal(err)
		}

		instance := NewDynamoDBInstance(context.Background(), "http://localhost:8000", "eu-central-1")

		_, err = batchWrite(instance, "user_variations", item)
		if err != nil {
			t.Fatal(err)
		}

		scanInput := &dynamodb.ScanInput{
			TableName: aws.String("user_variations"),
		}

		out, err := instance.db.Scan(scanInput)
		if err != nil {
			t.Fatal(err)
		}

		var want int64 = 1
		got := *out.Count
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}

		fmt.Println(out)
	})
}
