package tests

import (
	"ar/braze/braze"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Test1(t *testing.T) {
	dat, err := ioutil.ReadFile("./test.json")
	check(err)

	var e events.DynamoDBEvent
	err = json.Unmarshal(dat, &e)
	check(err)

	payload := braze.GetPayload(e)
	b, err := json.MarshalIndent(payload, "", "  ")
	fmt.Println(string(b))
}
