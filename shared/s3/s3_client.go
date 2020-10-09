package s3

import (
	"context"
	"kargo-back/shared/environment"
	"kargo-back/shared/logger"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var (
	// time in minutes that a presigned url is valid for
	presignedURLExpiration = environment.GetInt64("PRESIGNED_URL_EXPIRATION", 15)

	s3Client s3iface.S3API
)

func init() {
	sess := session.Must(session.NewSession())
	s3Client = s3.New(sess)
}

// GetGetPreSignedURL returns a presigned get url for the specified bucket and path
func GetGetPreSignedURL(ctx context.Context, bucket, path string) (string, error) {
	getObjectRequest, output := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})

	getObjectRequest.SetContext(ctx)

	defer logger.CloseOrLog(output.Body)

	urlString, err := getObjectRequest.Presign(time.Duration(presignedURLExpiration) * time.Minute)
	if err != nil {
		return "", err
	}

	return urlString, nil
}

// GetPutPreSignedURL returns a presigned put url for the specified bucket and path
func GetPutPreSignedURL(ctx context.Context, bucket, path string) (string, error) {
	putObjectRequest, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})

	putObjectRequest.SetContext(ctx)

	urlString, err := putObjectRequest.Presign(time.Duration(presignedURLExpiration) * time.Minute)
	if err != nil {
		return "", err
	}

	return urlString, nil
}
