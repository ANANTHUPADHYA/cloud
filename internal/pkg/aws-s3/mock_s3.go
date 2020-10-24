
package awss3

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// mock implementation
// Define a mock struct to be used in your unit tests of myFunc.

type mockS3 struct {
	s3iface.S3API
}

func (ms3 mockS3) GetObjectRequest(*s3.GetObjectInput) (*request.Request, *s3.GetObjectOutput) {
	mockOperation := request.Operation{}
	mockHttpRequest := http.Request{Method: "GET", URL: &url.URL{Host: "mock-Host"}}
	return &request.Request{Operation: &mockOperation, HTTPRequest: &mockHttpRequest}, &s3.GetObjectOutput{}
}

func (ms3 mockS3) DeleteObject(*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return &s3.DeleteObjectOutput{}, nil
}

func (ms3 mockS3) ListObjectsV2(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	var listObjects []*s3.Object
	obj1 := "mock-key-1"
	listObjects = append(listObjects, &s3.Object{Key: &obj1})
	return &s3.ListObjectsV2Output{Contents: listObjects}, nil
}

func (ms3 mockS3) DeleteObjects(*s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error) {
	return &s3.DeleteObjectsOutput{}, nil
}

type mockS3Err struct {
	s3iface.S3API
}

func (ms3 mockS3Err) ListObjectsV2(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	var listObjects []*s3.Object
	obj1 := "mock-key-2"
	listObjects = append(listObjects, &s3.Object{Key: &obj1})
	return &s3.ListObjectsV2Output{Contents: listObjects}, nil
}

func (ms3 mockS3Err) DeleteObject(*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return &s3.DeleteObjectOutput{}, fmt.Errorf("error deleting the object in s3")
}

func (ms3 mockS3Err) DeleteObjects(*s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error) {
	return &s3.DeleteObjectsOutput{}, fmt.Errorf("error deleting the objects in s3")
}

type mockS3ListErr struct {
	s3iface.S3API
}

func (ms3 mockS3ListErr) ListObjectsV2(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	var listObjects []*s3.Object
	obj1 := "mock-key-3"
	listObjects = append(listObjects, &s3.Object{Key: &obj1})
	return &s3.ListObjectsV2Output{Contents: listObjects}, fmt.Errorf("Error while listing list items from S3 bucket")
}

type mockS3EmptyListErr struct {
	s3iface.S3API
}

func (ms3 mockS3EmptyListErr) ListObjectsV2(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	var listObjects []*s3.Object
	return &s3.ListObjectsV2Output{Contents: listObjects}, nil
}
