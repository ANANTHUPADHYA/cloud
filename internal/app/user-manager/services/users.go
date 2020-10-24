package services

import (
	"context"
	"fmt"
	"github.com/cloud/internal/app/user-manager/constants"
	"github.com/cloud/internal/app/user-manager/models"
	"github.com/cloud/internal/pkg/database"
	commonModels "github.com/cloud/internal/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type UserService interface {
	CreateUser(ctx context.Context, input models.UserDynamo) (models.UserDynamo, *commonModels.ErrorResponse)
	GetAndValidateCredentials(ctx context.Context, userCred models.UserInputLogin) (models.CredIsAdmin, *commonModels.ErrorResponse)
	GetUser(ctx context.Context, userID string) (models.UserDynamo, *commonModels.ErrorResponse)
	GetUsers(ctx context.Context, query commonModels.DatabaseQuery)  (models.Users, *commonModels.ErrorResponse)
	UpdateUser(ctx context.Context, user models.UserDynamo) (models.UserDynamo, *commonModels.ErrorResponse)
	DeleteUser(ctx context.Context, userID string) (models.UserDynamo, *commonModels.ErrorResponse)
	GetAndValidateUser(ctx context.Context, userID string) (models.UserDynamo, *commonModels.ErrorResponse)
}

type UserManager struct {
	UserSvc database.UsersDynamoDBAPI
}

func NewUserService(dbSvc database.UsersDynamoDBAPI) UserService {
	return &UserManager{
		UserSvc: dbSvc,
	}
}

func (qm *UserManager) GetAndValidateCredentials(ctx context.Context, userCred models.UserInputLogin) (models.CredIsAdmin, *commonModels.ErrorResponse) {
	userCredResp, err := qm.UserSvc.GetUserCredentials(ctx, userCred.EmailAddress)
	if err != nil {
		return userCredResp, &commonModels.ErrorResponse{
			Message:         fmt.Sprintf("Error getting user credentials from database. %s", err.Error()),
			ErrorStatusCode: http.StatusInternalServerError,
		}
	}
	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(userCredResp.Password), []byte(userCred.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		return userCredResp, &commonModels.ErrorResponse{
			Message: fmt.Sprintf("Error validating password"),
			RecommendationAction: []string{fmt.Sprintf("Enter the correct password")},
			ErrorStatusCode: http.StatusUnauthorized,
		}
	}

	return userCredResp, nil
}

func (qm *UserManager) GetUsers(ctx context.Context, dbQuery commonModels.DatabaseQuery) (models.Users, *commonModels.ErrorResponse) {
	var users models.Users
	usersResp, err := qm.UserSvc.GetUsersInDynamoDB(ctx, dbQuery)
	if err != nil {
		return users, &commonModels.ErrorResponse{
			Message:         fmt.Sprintf("Error getting users from database. %s", err.Error()),
			ErrorStatusCode: http.StatusInternalServerError,
		}
	}

	for _, userDynamo := range usersResp {
		users.Members = append(users.Members, userDynamo.User)
	}
	return users, nil
}

func (um *UserManager) CreateUser(ctx context.Context, input models.UserDynamo) (models.UserDynamo, *commonModels.ErrorResponse) {
	var err error
	var userCreateResp models.UserDynamo

	userCreateResp, err = um.UserSvc.CreateUserInDynamoDB(ctx, input, "")
	if err != nil {
		return userCreateResp, &commonModels.ErrorResponse{
			Message:         fmt.Sprintf("Error creating user. %s", err.Error()),
			ErrorStatusCode: http.StatusInternalServerError,
		}
	}

	return userCreateResp, nil
}

// GetUser gets the user from dynamo db
func (um *UserManager) GetUser(ctx context.Context, userID string) (models.UserDynamo, *commonModels.ErrorResponse) {
	userResp, err := um.UserSvc.GetUserInDynamoDB(ctx, userID, constants.TypeUsersForSortKey)
	if err != nil {
		return models.UserDynamo{}, &commonModels.ErrorResponse{
			Message:         fmt.Sprintf("Erro getting user. %s", err.Error),
			ErrorStatusCode: http.StatusInternalServerError,
		}
	}
	if userResp.UserID == "" {
		return models.UserDynamo{}, &commonModels.ErrorResponse{
			Message:         fmt.Sprintf("The given user %s doesn't exist in the dynamo db", userID),
			ErrorStatusCode: http.StatusBadRequest,
		}
	}
	return userResp, nil
}

// GetAndValidateUser will get the user from db and validate
func (um *UserManager) GetAndValidateUser(ctx context.Context, userID string) (models.UserDynamo, *commonModels.ErrorResponse) {
	log.Printf("entering get and validate user mentsh")
	user, err := um.UserSvc.GetUserInDynamoDB(ctx, userID, constants.TypeUsersForSortKey)
	log.Printf("User details %v %v", user, err)
	if err != nil {
		return user, &commonModels.ErrorResponse{
			Message:         fmt.Sprintf("Error getting user. %s", err.Error()),
			ErrorStatusCode: http.StatusInternalServerError,
		}
	}
	if user.UserID == "" {
		return user, &commonModels.ErrorResponse{
			Message:         fmt.Sprint("User ID is empty"),
			ErrorStatusCode: http.StatusBadRequest,
		}
	}
	return user, nil
}

// UpdateUser updates user
func (um *UserManager) UpdateUser(ctx context.Context, userInput models.UserDynamo) (models.UserDynamo, *commonModels.ErrorResponse) {
	err := um.UserSvc.DeleteUserInDynamoDB(ctx, userInput.UserID, "user")
	if err != nil {
		return models.UserDynamo{}, &commonModels.ErrorResponse{
			Message:         fmt.Sprintf("Error occured while updating the user. %s", err.Error()),
			ErrorStatusCode: http.StatusInternalServerError,
		}
	}
	userUpdateResp, err := um.UserSvc.CreateUserInDynamoDB(ctx, userInput, "")
	if err != nil {
		return userUpdateResp, &commonModels.ErrorResponse{
			Message:         fmt.Sprintf("Error updating user. %s", err.Error()),
			ErrorStatusCode: http.StatusInternalServerError,
		}
	}
	return userUpdateResp, nil
}

// DeleteUser deletes the user from dynamo db
func (um *UserManager) DeleteUser(ctx context.Context, userID string) (models.UserDynamo, *commonModels.ErrorResponse) {
	userResp := models.UserDynamo{}
	err := um.UserSvc.DeleteUserInDynamoDB(ctx, userID, constants.TypeUsersForSortKey)
	if err != nil {
		return models.UserDynamo{}, &commonModels.ErrorResponse{
			Message:         fmt.Sprintf("Error deleting user. %s", err.Error),
			ErrorStatusCode: http.StatusInternalServerError,
		}
	}
	return userResp, nil
}