package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/g-hyoga/kyuko/src/data"
	"github.com/g-hyoga/kyuko/src/kyuko"
	"github.com/g-hyoga/kyuko/src/twitter"
)

var (
	T_CONSUMER_KEY        = os.Getenv("T_CONSUMER_KEY")
	T_CONSUMER_SECRET     = os.Getenv("T_CONSUMER_SECRET")
	T_ACCESS_TOKEN        = os.Getenv("T_ACCESS_TOKEN")
	T_ACCESS_TOKEN_SECRET = os.Getenv("T_ACCESS_TOKEN_SECRET")

	I_CONSUMER_KEY        = os.Getenv("I_CONSUMER_KEY")
	I_CONSUMER_SECRET     = os.Getenv("I_CONSUMER_SECRET")
	I_ACCESS_TOKEN        = os.Getenv("I_ACCESS_TOKEN")
	I_ACCESS_TOKEN_SECRET = os.Getenv("I_ACCESS_TOKEN_SECRET")
)

func kyukoHandler(ctx context.Context) (data.Response, error) {
	kyukoData := []data.KyukoData{}

	iClient, err := twitter.NewTwitterClient(I_CONSUMER_KEY, I_CONSUMER_SECRET, I_ACCESS_TOKEN, I_ACCESS_TOKEN_SECRET)
	if err != nil {
		return data.Response{Data: kyukoData, Error: err}, err
	}

	kyukoData, err = kyuko.Exec(1, iClient)
	if err != nil {
		return data.Response{Data: kyukoData, Error: err}, err
	}

	return data.Response{Data: kyukoData, Error: nil}, nil
}

func main() {
	lambda.Start(kyukoHandler)
}
