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

		SetRemovedKeysToNil(record.Change.OldImage, record.Change.NewImage)

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

		// Rename uuid to external_id
		out["external_id"] = out["uuid"]
		delete(out, "uuid")

		payload.Attributes = append(payload.Attributes, out)

		//jsonString, err := json.Marshal(out)
		//fmt.Println(err, string(jsonString))
	}

	return payload
}

func SetRemovedKeysToNil(old map[string]events.DynamoDBAttributeValue, new map[string]events.DynamoDBAttributeValue) {
	for k := range old {
		if _, ok := new[k]; !ok {
			new[k] = events.NewNullAttribute()
		}
	}
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
