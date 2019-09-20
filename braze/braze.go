package braze

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Payload struct {
	APIKey     string                   `json:"api_key"`
	Attributes []map[string]interface{} `json:"attributes"`
}

func HandleRequest(ctx context.Context, e events.DynamoDBEvent) {
	payload := GetBrazeJSON(e)

	TrackUsers(payload)
}

func GetBrazeJSON(e events.DynamoDBEvent) *Payload {
	payload := new(Payload)
	payload.APIKey = os.Getenv("BRAZE_API_KEY")
	payload.Attributes = []map[string]interface{}{}

	for _, record := range e.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

		oldData, err := ConvertToMap(record.Change.OldImage)
		if err != nil {
			panic(err)
		}

		newData, err := ConvertToMap(record.Change.NewImage)
		if err != nil {
			panic(err)
		}

		out := ChangeForBraze(oldData, newData)

		payload.Attributes = append(payload.Attributes, out)
	}

	return payload
}

func ConvertToMap(old map[string]events.DynamoDBAttributeValue) (map[string]interface{}, error) {
	data, err := json.Marshal(old)
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

func ChangeForBraze(in map[string]interface{}, out map[string]interface{}) map[string]interface{} {
	uuid := in["uuid"].(string)
	for k := range in {
		_, exists := out[k]
		if exists {
			if in[k] == out[k] {
				delete(out, k)
			}
		} else {
			out[k] = nil
		}

	}

	// Rename uuid to external_id
	out["external_id"] = uuid

	return out
}

func TrackUsers(payload *Payload) {
	url := "https://rest.fra-01.braze.eu/users/track"

	jsonString, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))

	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
