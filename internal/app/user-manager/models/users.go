package models

import fileModels "github.com/cloud/internal/app/files-manager/models"

// User represents the users
type Users struct {
	Members []User `json:"members"`
}

type DynamoKeys struct {
	PKey string `json:PKey`
	SKey string `json:SKey`
}

type Credentials struct {
	EmailAddress string                         `json:emailAddress`
	Password     []byte
}

type CredIsAdmin struct {
	EmailAddress string                         `json:emailAddress`
	Password     []byte
	IsAdmin      bool                           `json:isAdmin`
}

type User struct {
	UserID       string                         `json:userID`
	FirstName    string                         `json:firstName`
	LastName     string                         `json:lastName,omitempty`
	IsAdmin      bool                           `json:isAdmin`
	FileInfo     map[string]fileModels.FileInfo `json:"files,omitempty"`
	Credentials
}

type UserDynamo struct {
	DynamoKeys
	User
}

type UserInput struct {
	FirstName    string `json:firstName`
	LastName     string `json:lastName, omitempty`
	EmailAddress string `json:emailAddress`
	IsAdmin      bool   `json:isAdmin`
	Password     string `json:password`
}

type UserInputLogin struct {
	EmailAddress string `json:emailAddress`
	Password string `json:password`
}

type UserOutput struct {
	FirstName string `json:firstName`
	LastName  string `json:lastName, omitempty`
}

// DefaultQuery is default query for building expression
type DefaultQuery struct {
	Key   string
	Value string
}
