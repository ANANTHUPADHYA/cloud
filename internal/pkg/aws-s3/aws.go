package awss3

//go:generate mockgen -source ./aws.go -package awss3 -destination ./aws_mock.go

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/ANANTHUPADHYA/cloud/internal/pkg/http/transport"
	"github.com/ANANTHUPADHYA/cloud/internal/pkg/models"
	"github.com/ANANTHUPADHYA/cloud/internal/pkg/utils"
)

const (
	awsS3AccessKey         = "S3_AWS_ACCESS_KEY_ID"
	awsS3Secret            = "S3_AWS_SECRET_ACCESS_KEY"
	awsS3Region            = "S3_AWS_DEFAULT_REGION"
	awsS3BucketName        = "S3_TITAN_BUCKET_NAME"
	defaultAwsS3BucketName = "glc-quote-intg-attachments"
)

//AWSCredentials - holds the function to get the AWS credentials
type AWSCredentials interface {
	GetAwsCredDetails(ctx context.Context) (string, string, string, error)
	GetSession(id string, awsCredentials models.AwsCredentials) (*session.Session, error)
	GetAwsS3BucketName(context.Context) string
}

type awsCreds struct{}

// NewAWSCredsImpl creds implementation
func NewAWSCredsImpl() AWSCredentials {
	return &awsCreds{}
}

// GetAwsCredDetails - to get all aws credential details
func (awsCreds *awsCreds) GetAwsCredDetails(context.Context) (string, string, string, error) {
	ctx := context.Background()
	awsAccessKey, set := os.LookupEnv(awsS3AccessKey)
	if !set {
		err := fmt.Errorf("AWS access key is not set in ENV %s", awsS3AccessKey)
		return "", "", "", err
	}
	awsSecretKey, set := os.LookupEnv(awsS3Secret)
	if !set {
		err := fmt.Errorf("AWS secret key is not set in ENV %s", awsS3Secret)
		return awsAccessKey, "", "", err
	}
	awsRegion := utils.GetEnvOrDefault(awsS3Region, "us-east-1")
	return awsAccessKey, awsSecretKey, awsRegion, nil
}

// GetSession will return aws session
func (awsCreds *awsCreds) GetSession(id string, awsCredentials models.AwsCredentials) (*session.Session, error) {
	loggerTransport := transport.NewLoggerTransport(http.DefaultTransport)
	metricTransport := transport.NewMetricTransport(loggerTransport)

	// adding timeout explicitly to limit the wait time for each client call
	// keeping it as 15 seconds to avoid the request timing out too soon
	httpClient := &http.Client{
		Timeout:   15 * time.Second,
		Transport: metricTransport,
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(id),
		Credentials: credentials.NewStaticCredentials(awsCredentials.AccessKey, awsCredentials.SecretKey, ""),
		HTTPClient:  httpClient,
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// GetAwsCredDetails - to get all aws credential details
func (awsCreds *awsCreds) GetAwsS3BucketName(context.Context) string {
	s3BucketName := utils.GetEnvOrDefault(awsS3BucketName, defaultAwsS3BucketName)
	return s3BucketName
}
