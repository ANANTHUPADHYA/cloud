package models

import "github.com/aws/aws-sdk-go/service/dynamodb"

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
}

// AwsCredentials are the security details needed to access aws for a specific provider
type AwsCredentials struct {
	AccessKey   string `json:"accessKey,omitempty"`
	SecretKey   string `json:"secretKey,omitempty"`
	EndpointURL string `json:"endpoint,omitempty"`
}

// DatabaseQuery represents the database query type
type DatabaseQuery struct {
	Equal       QueryMap
	NotEqual    QueryMap
	Default     DefaultQuery
	QueryParams *dynamodb.QueryInput
}

// QueryMap is query map
type QueryMap map[string][]string

// DefaultQuery is default query for building expression
type DefaultQuery struct {
	Key   string
	Value string
}
