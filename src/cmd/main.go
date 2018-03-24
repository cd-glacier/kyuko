package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Name string `json:"name"`
}

type Response struct {
	Message string `json:message`
}

func kyukoHandler(ctx context.Context, event Event) (Response, error) {
	return Response{Message: fmt.Sprintf("hello %s!", event.Name)}, nil
}

func main() {
	lambda.Start(kyukoHandler)
}
