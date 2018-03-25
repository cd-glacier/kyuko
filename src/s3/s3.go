package s3

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/g-hyoga/kyuko/src/data"
)

type S3 struct {
	client *awsS3.S3
	bucket *string
}

func (s3 *S3) GetClient(bucket string, region string) {
	if bucket == "" || region == "" {
		errors.New("S3 bucket name or region is missing")
	}

	s3.client = awsS3.New(session.New(), &aws.Config{
		Region: aws.String(region),
	})
	s3.bucket = aws.String(bucket)
}

func (s3 *S3) Put(key string, object []data.KyukoData) (*awsS3.PutObjectOutput, error) {
	jsonObj, err := toJson(object)
	if err != nil {
		return nil, err
	}

	input := &awsS3.PutObjectInput{
		Body:   aws.ReadSeekCloser(strings.NewReader(jsonObj)),
		Bucket: s3.bucket,
		Key:    aws.String(key),
	}

	result, err := s3.client.PutObject(input)
	if err != nil {
		return result, err
	}
	return result, err
}

func toJson(object []data.KyukoData) (string, error) {
	b, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
