package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, e events.DynamoDBEvent) {
	for _, record := range e.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

		Keys(record.Change.OldImage, record.Change.NewImage)

		data, err := json.Marshal(record.Change.NewImage)
		if err != nil {
			panic(err)
		}

		av := make(map[string]*dynamodb.AttributeValue)
		if err := json.Unmarshal(data, &av); err != nil {
			panic(err)
		}

		var out map[string]interface{}
		if err := dynamodbattribute.UnmarshalMap(av, &out); err != nil {
			panic(err)
		}

		jsonString, err := json.Marshal(out)
		fmt.Println(err, string(jsonString))
	}
}

func Keys(old map[string]events.DynamoDBAttributeValue, new map[string]events.DynamoDBAttributeValue) {
	for k := range old {
		if _, ok := new[k]; !ok {
			new[k] = events.NewNullAttribute()
		}
	}
}
