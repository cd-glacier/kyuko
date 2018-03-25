package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/g-hyoga/kyuko/src/data"
	"github.com/g-hyoga/kyuko/src/kyuko"
	myS3 "github.com/g-hyoga/kyuko/src/s3"
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

	S3_BUCKET_NAME = os.Getenv("S3_BUCKET_NAME")
)

func kyukoHandler(ctx context.Context) (data.Response, error) {
	kyukoData := []data.KyukoData{}

	var s3 myS3.S3
	s3.GetClient(S3_BUCKET_NAME, endpoints.ApNortheast1RegionID)

	// Imadegawa
	iClient, err := twitter.NewTwitterClient(I_CONSUMER_KEY, I_CONSUMER_SECRET, I_ACCESS_TOKEN, I_ACCESS_TOKEN_SECRET)
	if err != nil {
		return data.Response{Data: kyukoData, Error: err}, err
	}

	kyukoData, err = kyuko.Exec(1, iClient, s3)
	if err != nil {
		return data.Response{Data: kyukoData, Error: err}, err
	}

	// Kyoutanabe
	tClient, err := twitter.NewTwitterClient(T_CONSUMER_KEY, T_CONSUMER_SECRET, T_ACCESS_TOKEN, T_ACCESS_TOKEN_SECRET)
	if err != nil {
		return data.Response{Data: kyukoData, Error: err}, err
	}

	kyukoData, err = kyuko.Exec(2, tClient, s3)
	if err != nil {
		return data.Response{Data: kyukoData, Error: err}, err
	}

	return data.Response{Data: kyukoData, Error: nil}, nil
}

func main() {
	lambda.Start(kyukoHandler)
}
