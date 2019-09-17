package main

import (
	"ar/braze/braze"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(braze.HandleRequest)
}
