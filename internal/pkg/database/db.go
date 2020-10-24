package database

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/ANANTHUPADHYA/cloud/internal/app/user-manager/constants"
	"github.com/ANANTHUPADHYA/cloud/internal/app/user-manager/models"
	dbModels "github.com/ANANTHUPADHYA/cloud/internal/pkg/models"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
)

const (
	awsAccess        = "DB_AWS_ACCESS_KEY_ID"
	awsSecret        = "DB_AWS_SECRET_ACCESS_KEY"
	awsRegion        = "DB_AWS_REGION"
	dynamodbEndpoint = "DYNAMODB_ENDPOINT_URL"
)

type awsCreds struct{}

// AWSServiceSessions - holds the function to get the AWS credentials
type AWSServiceSessions interface {
	GetAwsCredDetails() (string, string, string, string, error)
	GetDynamodbSVC(httpClient *http.Client) (*dynamodb.DynamoDB, error)
}

//NewAWSCredsImpl give dynamodb API implementation
func NewAWSCredsImpl() AWSServiceSessions {
	return &awsCreds{}
}

func getSession(id string, awsCredentials dbModels.AwsCredentials, httpClient *http.Client) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(awsCredentials.EndpointURL),
		Region:      aws.String(id),
		Credentials: credentials.NewStaticCredentials(awsCredentials.AccessKey, awsCredentials.SecretKey, ""),
		HTTPClient:  httpClient,
	})
	return sess, err
}

// GetAwsCredDetails get aws credential details
func (awsCreds *awsCreds) GetAwsCredDetails() (string, string, string, string, error) {
	awsAccessKey, set := os.LookupEnv(awsAccess)
	if !set {
		err := fmt.Errorf("AWS access key is not set in ENV %s", awsAccess)
		return "", "", "", "", err
	}
	awsSecretKey, set := os.LookupEnv(awsSecret)
	if !set {
		err := fmt.Errorf("AWS secret key is not set in ENV %s", awsSecret)
		return awsAccessKey, "", "", "", err
	}
	awsRegion, set := os.LookupEnv(awsRegion)
	if !set {
		err := fmt.Errorf("AWS region is not set in ENV %s", awsRegion)
		return awsAccessKey, awsSecretKey, "", "", err
	}
	dynamoDBEndpoint, set := os.LookupEnv(dynamodbEndpoint)
	if !set {
		err := fmt.Errorf("dynamoDB endpoint is not set in ENV %s", dynamodbEndpoint)
		return awsAccessKey, awsSecretKey, awsRegion, "", err
	}

	return awsAccessKey, awsSecretKey, awsRegion, dynamoDBEndpoint, nil
}

// GetDynamodbSVC returns the dynamodb svc
func (awsCreds *awsCreds) GetDynamodbSVC(httpClient *http.Client) (*dynamodb.DynamoDB, error) {
	awsAccessKey, awsSecretKey, awsRegion, dynamoDBEndpoint, err := awsCreds.GetAwsCredDetails()

	if err != nil {
		return &dynamodb.DynamoDB{}, err
	}
	creds := dbModels.AwsCredentials{
		AccessKey:   awsAccessKey,
		SecretKey:   awsSecretKey,
		EndpointURL: dynamoDBEndpoint,
	}

	sess, err := getSession(awsRegion, creds, httpClient)
	if err != nil {
		return &dynamodb.DynamoDB{}, err
	}
	svc := dynamodb.New(sess)

	return svc, nil
}

type UsersDynamoDBAPI interface {
	CreateUserInDynamoDB(ctx context.Context, input models.UserDynamo, condition string) (models.UserDynamo, error)
	GetUserInDynamoDB(ctx context.Context, pkey string, skey string) (models.UserDynamo, error)
	GetUsersInDynamoDB(ctx context.Context, query dbModels.DatabaseQuery) ([]models.UserDynamo, error)
	GetUserCredentials(ctx context.Context, userEmail string) (models.CredIsAdmin , error)
	DeleteUserInDynamoDB(ctx context.Context, pkey string, skey string) error
}

type userDynamodbImpl struct {
	usrSvc dynamodbiface.DynamoDBAPI
}

func NewUsersDBImpl(usrSvc dynamodbiface.DynamoDBAPI) userDynamodbImpl {
	return userDynamodbImpl{
		usrSvc: usrSvc,
	}
}

func (dbImpl userDynamodbImpl) CreateUserInDynamoDB(ctx context.Context, userInput models.UserDynamo, condition string) (models.UserDynamo, error) {
	var userOutput models.UserDynamo

	av, err := dynamodbattribute.MarshalMap(userInput)
	if err != nil {
		return userOutput, err
	}
	log.Printf("Item is printed %v", av)
	var expr *string
	if condition != "" {
		expr = aws.String(condition)
	}

	input := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(constants.UsersTableName),
		ConditionExpression: expr,
	}

	_, err = dbImpl.usrSvc.PutItem(input)
	if err != nil {
		return userOutput, err
	}

	return userInput, nil
}

// GetUserInDynamoDB gets user from db id
func (dbImpl userDynamodbImpl) GetUserInDynamoDB(ctx context.Context, pkey string, skey string) (models.UserDynamo, error) {
	user := models.UserDynamo{}
	key := map[string]*dynamodb.AttributeValue{
		constants.UsersTablePrimaryKey: {
			S: aws.String(pkey),
		},
		constants.UsersTableSortKey: {
			S: aws.String(skey),
		},
	}
	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(constants.UsersTableName),
	}
	// Make the DynamoDB Query API call
	result, err := dbImpl.usrSvc.GetItem(input)
	if err != nil {
		return user, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (dbImpl userDynamodbImpl) GetUserCredentials(ctx context.Context, userEmail string) (models.CredIsAdmin , error) {
	userCreds := models.CredIsAdmin{}
	filt := expression.Name("EmailAddress").Equal(expression.Value(userEmail))

	// Get back the title, year, and rating
	proj := expression.NamesList(expression.Name("EmailAddress"), expression.Name("Password"), expression.Name("IsAdmin"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		return userCreds, err
	}
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:     expr.Filter(),
		ProjectionExpression: expr.Projection(),
		TableName:                 aws.String(constants.UsersTableName),
	}

	// Make the DynamoDB Query API call
	result, err := dbImpl.usrSvc.Scan(input)
	if err != nil {
		return userCreds, err
	}

	if len(result.Items) == 0 {
		return userCreds, errors.New("No user found with this email")
	}
	if len(result.Items) > 1 {
		return userCreds, errors.New("More than one items found for the provided userEmail")
	}
	fmt.Printf("Result %v", result.Items)
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &userCreds)
	if err != nil {
		return userCreds, err
	}

	return userCreds, nil
}


// DeleteUserInDynamoDB gets user from db id
func (dbImpl userDynamodbImpl) DeleteUserInDynamoDB(ctx context.Context, pkey string, skey string) error {
	key := map[string]*dynamodb.AttributeValue{
		constants.UsersTablePrimaryKey: {
			S: aws.String(pkey),
		},
		constants.UsersTableSortKey: {
			S: aws.String(skey),
		},
	}
	input := &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(constants.UsersTableName),
	}
	// Make the DynamoDB Query API call
	_, err := dbImpl.usrSvc.DeleteItem(input)
	if err != nil {
		return err
	}
	return nil
}

// GetUsersInDynamoDBUsingQuery gets users from dynamo db
func (dbImpl *userDynamodbImpl) GetUsersInDynamoDB(ctx context.Context, query dbModels.DatabaseQuery) ([]models.UserDynamo, error) {
	listUsers := []models.UserDynamo{}

	if query.Default.Key == "" || query.Default.Value == "" {
		query.Default = dbModels.DefaultQuery{
			Key:   constants.UsersTableSortKey,
			Value: constants.TypeUsersForSortKey,
		}
	}
	input, err := buildQueryDynamoDB(ctx, query, constants.UsersTableName)
	if err != nil {
		return listUsers, errors.Wrap(err, "Error building query expression")
	}

	for {
		// Make the DynamoDB Query API call
		result, err := dbImpl.usrSvc.Scan(input)
		if err != nil {
			return listUsers, err
		}
		user := []models.UserDynamo{}
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &user)
		if err != nil {
			return listUsers, err
		}
		listUsers = append(listUsers, user...)
		if result.LastEvaluatedKey == nil {
			break
		}
		input.ExclusiveStartKey = result.LastEvaluatedKey
	}

	return listUsers, nil
}

func buildQueryDynamoDB(ctx context.Context, query dbModels.DatabaseQuery, tableName string) (*dynamodb.ScanInput, error) {
	// Query with sort key as "order" by default, if not specified otherwise
	if query.Default.Key == "" || query.Default.Value == "" {
		query.Default = dbModels.DefaultQuery{
			Key:   constants.UsersTableSortKey,
			Value: constants.TypeUsersForSortKey,
		}
	}

	filter := expression.Name(query.Default.Key).Equal(expression.Value(query.Default.Value))
	// Set the equal filters
	for key, vals := range query.Equal {
		var exprVals []expression.OperandBuilder
		for _, val := range vals {
			exprVals = append(exprVals, expression.Value(val))
		}
		filter = filter.And(expression.Name(key).In(exprVals[0], exprVals[1:]...))
	}

	// set the not equal filters
	for key, vals := range query.NotEqual {
		var exprVals []expression.OperandBuilder
		for _, val := range vals {
			exprVals = append(exprVals, expression.Value(val))
		}
		filter = filter.And(expression.Name(key).In(exprVals[0], exprVals[1:]...).Not())
	}
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return &dynamodb.ScanInput{}, errors.Wrap(err, "error building expression")
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
	}

	return params, nil
}