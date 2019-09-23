package braze

import (
	"ar/braze/dynamodb"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

type Payload struct {
	APIKey     string                   `json:"api_key"`
	Attributes []map[string]interface{} `json:"attributes"`
}

func HandleRequest(ctx context.Context, e events.DynamoDBEvent) {
	payload := GetPayload(e)

	res, err := TrackUsers(payload)
	if err != nil {
		fmt.Println(err.Error())

		if res != nil {
			defer res.Body.Close()
			resBody, _ := ioutil.ReadAll(res.Body)
			fmt.Println(string(resBody))
		}
	}
}

func GetPayload(e events.DynamoDBEvent) *Payload {
	payload := new(Payload)
	payload.APIKey = os.Getenv("BRAZE_API_KEY")
	payload.Attributes = []map[string]interface{}{}

	for _, record := range e.Records {
		// ignore delete events
		if record.EventName == "REMOVE" {
			continue
		}

		attribute, err := GetAttribute(record)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		payload.Attributes = append(payload.Attributes, attribute)
	}

	return payload
}

func GetAttribute(record events.DynamoDBEventRecord) (map[string]interface{}, error) {
	fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

	oldImage, err := dynamodb.ConvertAVToMap(record.Change.OldImage)
	if err != nil {
		return nil, err
	}

	newImage, err := dynamodb.ConvertAVToMap(record.Change.NewImage)
	if err != nil {
		return nil, err
	}

	return ChangeForBraze(oldImage, newImage), nil
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

func TrackUsers(payload *Payload) (*http.Response, error) {
	url := "https://rest.fra-01.braze.eu/users/track"

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		return res, fmt.Errorf("response was not 201 Status Created, was status %d", res.StatusCode)
	}

	return res, nil
}
