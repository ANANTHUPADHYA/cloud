

package awss3

//go:generate mockgen -source ./aws_s3.go -package awss3 -destination ./aws_s3_mock.go

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	fileModels "github.com/ANANTHUPADHYA/cloud/internal/app/files-manager/models"
	"github.com/ANANTHUPADHYA/cloud/internal/pkg/models"
)

const (
	presignTime = 2
)

// IfAWSS3 holds the aws s3 functions
type IfAWSS3 interface {
	GenerateS3PresignedURL(ctx context.Context, userID string, filename string) (fileModels.DownloadFileInfo, error)
	UploadAttachmentTOS3Bucket(ctx context.Context, userID string, filename string, filereader io.Reader) error
	DeleteFileInS3(ctx context.Context, quoteID string, filename string) error
	//DeleteQuoteAttachmentFolderInS3(ctx context.Context, quoteID string, tenantID string) error
	GetAWSS3Session() (*session.Session, error)
}

// IfAWSS3SvcImpl returns the aws s3 service
type IfAWSS3SvcImpl interface {
	GetS3SVC() (*s3.S3, error)
}

type awsS3 struct {
	awsS3API s3iface.S3API
	awsCreds AWSCredentials
}

type s3ifimpl struct {
	awsCreds AWSCredentials
}

// NewAWSS3Impl returns aws s3 service implementation
func NewAWSS3Impl(creds AWSCredentials) IfAWSS3SvcImpl {
	return &s3ifimpl{awsCreds: creds}
}

//NewAWSS3 give aws s3 implementation of IfAWSS3
func NewAWSS3(awsS3APIImpl s3iface.S3API, creds AWSCredentials) IfAWSS3 {
	return &awsS3{awsS3API: awsS3APIImpl, awsCreds: creds}
}

// GetS3SVC returns the aws s3 svc
func (awss3 s3ifimpl) GetS3SVC() (*s3.S3, error) {
	ctx := context.Background()
	awsAccessKey, awsSecretKey, awsRegion, err := awss3.awsCreds.GetAwsCredDetails(ctx)
	if err != nil {
		return &s3.S3{}, fmt.Errorf("Error getting credentials")
	}

	creds := models.AwsCredentials{
		AccessKey: awsAccessKey,
		SecretKey: awsSecretKey,
	}

	sess, err := awss3.awsCreds.GetSession(awsRegion, creds)
	if err != nil {
		return &s3.S3{}, fmt.Errorf("Error getting AWS session")
	}
	svc := s3.New(sess)
	return svc, nil
}

func (awss3 awsS3) GenerateS3PresignedURL(ctx context.Context, quoteID string, fileName string) (fileModels.DownloadFileInfo, error) {
	downloadAttachInfo := fileModels.DownloadFileInfo{}
	awsS3BucketName := awss3.awsCreds.GetAwsS3BucketName(ctx)
	objectURL := "/" + "/" + quoteID + "/"
	oURL := objectURL + fileName
	req, _ := awss3.awsS3API.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(awsS3BucketName),
		Key:    aws.String(oURL),
	})
	urlStr, err := req.Presign(presignTime * time.Minute)
	if err != nil {
		return downloadAttachInfo, err
	}
	downloadAttachInfo.PresignedURL = urlStr
	return downloadAttachInfo, nil
}

func (awss3 awsS3) GetAWSS3Session() (*session.Session, error) {
	ctx := context.Background()
	awsAccessKey, awsSecretKey, awsRegion, err := awss3.awsCreds.GetAwsCredDetails(ctx)
	if err != nil {
		return &session.Session{}, fmt.Errorf("Error getting credentials")
	}

	creds := models.AwsCredentials{
		AccessKey: awsAccessKey,
		SecretKey: awsSecretKey,
	}

	sess, err := awss3.awsCreds.GetSession(awsRegion, creds)
	if err != nil {
		return &session.Session{}, fmt.Errorf("Error getting AWS session")
	}
	return sess, nil
}

func (awss3 awsS3) UploadAttachmentTOS3Bucket(ctx context.Context, userID string, filename string, filereader io.Reader) error {
	awsS3BucketName := awss3.awsCreds.GetAwsS3BucketName(ctx)
	objectURL := "/" + userID + "/" + filename
	sess, err := awss3.GetAWSS3Session()
	if err != nil {
		sessErr := fmt.Sprintf("Error while getting aws s3 session. Error: %s", err.Error())
		return fmt.Errorf(sessErr)
	}
	uploader := s3manager.NewUploader(sess)
	log.Print("Successful till this stage")
	// Upload the file to S3.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsS3BucketName),
		Key:    aws.String(objectURL),
		Body:   filereader,
	})
	if err != nil {
		sessErr := fmt.Sprintf("Error while uploading file %s for user %s. Error: %s", filename, userID, err.Error())
		return fmt.Errorf(sessErr)
	}
	return nil
}

func (awss3 awsS3) DeleteFileInS3(ctx context.Context, userID string, filename string) error {
	awsS3BucketName := awss3.awsCreds.GetAwsS3BucketName(ctx)
	objectURL := "/" + userID + "/" + filename
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(awsS3BucketName),
		Key:    aws.String(objectURL),
	}
	_, err := awss3.awsS3API.DeleteObject(input)
	if err != nil {
		delErr := fmt.Sprintf("Error while file %s for user  %s. Error: %s", filename, userID, err.Error())
		return fmt.Errorf(delErr)
	}
	return nil
}

func (awss3 awsS3) DeleteQuoteAttachmentFolderInS3(ctx context.Context, quoteID string, tenantID string) error {
	awsS3BucketName := awss3.awsCreds.GetAwsS3BucketName(ctx)
	objectURL := tenantID + "/" + quoteID
	var s3Objects []*s3.ObjectIdentifier
	resp, err := awss3.awsS3API.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(awsS3BucketName),
		Prefix: aws.String(objectURL)})
	if err != nil {
		delErr := fmt.Sprintf("Error while listing list items from S3 bucket. Error: %s", err.Error())
		return fmt.Errorf(delErr)
	}
	if len(resp.Contents) == 0 {
		delErr := fmt.Errorf("warning: Got empty list of quote attachments from S3 bucket for quote %s", quoteID)
		return delErr
	}
	for _, item := range resp.Contents {
		s3Objects = append(s3Objects, &s3.ObjectIdentifier{Key: item.Key})
	}
	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(awsS3BucketName),
		Delete: &s3.Delete{
			Objects: s3Objects,
			Quiet:   aws.Bool(false),
		},
	}
	_, err = awss3.awsS3API.DeleteObjects(input)
	if err != nil {
		delErr := fmt.Sprintf("Error while delete quote attachments for quote %s. Error: %s", quoteID, err.Error())
		return fmt.Errorf(delErr)
	}
	return nil
}
