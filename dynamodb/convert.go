package dynamodb

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Converts a DynamoDB Attribute Value into a Map
func ConvertAVToMap(in map[string]events.DynamoDBAttributeValue) (map[string]interface{}, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	av := make(map[string]*dynamodb.AttributeValue)
	if err := json.Unmarshal(data, &av); err != nil {
		return nil, err
	}

	var out map[string]interface{}
	if err := dynamodbattribute.UnmarshalMap(av, &out); err != nil {
		return nil, err
	}

	return out, nil
}
