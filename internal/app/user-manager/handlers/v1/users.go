package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ANANTHUPADHYA/cloud/internal/app/files-manager/constants"
	"github.com/ANANTHUPADHYA/cloud/internal/app/user-manager/models"
	"github.com/ANANTHUPADHYA/cloud/internal/app/user-manager/services"
	errModels "github.com/ANANTHUPADHYA/cloud/internal/pkg/models"
	"github.com/ANANTHUPADHYA/cloud/internal/pkg/utils"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

type UMSRest struct {
	UserService services.UserService
}

func CreateUMSRouter(userService services.UserService) *UMSRest {
	return &UMSRest{
		UserService: userService,
	}
}

func (ur *UMSRest) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var userInput models.UserInput
	err := c.BindJSON(&userInput)
	if err != nil {
		errRes := errModels.ErrorResponse{
			Message:              "Invalid request body",
			RecommendationAction: []string{"Check user create request body"},
			ErrorStatusCode:      http.StatusBadRequest,
		}
		fmt.Printf("Error while binding json to user input %s", err)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 8)
	fmt.Printf("jas  %v", hashedPassword)
	if err != nil {
		errRes := errModels.ErrorResponse{
			Message:              "Invalid password entered",
			RecommendationAction: []string{"Check password parameter"},
			ErrorStatusCode:      http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errRes)
	}
	userUUID := utils.GenerateUUID()
	userDynamo := models.UserDynamo{
		DynamoKeys: models.DynamoKeys{
			PKey: userUUID,
			SKey: "user",
		},
		User: models.User{
			UserID:       userUUID,
			FirstName:    userInput.FirstName,
			LastName:     userInput.LastName,
			IsAdmin:      userInput.IsAdmin,
			Credentials: models.Credentials{
				EmailAddress: userInput.EmailAddress,
				Password:     hashedPassword,
			},
		},
	}
	userCreateResp, errCreate := ur.UserService.CreateUser(ctx, userDynamo)
	if errCreate != nil {
		errMsg := fmt.Sprintf("Error creating user. %s", errCreate.Message)
		errCode := errCreate.ErrorStatusCode
		if errCode == 0 {
			errCode = http.StatusInternalServerError
		}
		errResp := errModels.ErrorResponse{
			Message:         errMsg,
			ErrorStatusCode: errCode,
		}
		fmt.Printf("Error while creating user %s", err)
		c.JSON(errCode, errResp)
		return
	}
	//fmt.Printf("Successfully created user. %s with ID %s", userCreateResp.FirstName, userCreateResp.UserID)
	c.JSON(http.StatusCreated, userCreateResp.User)
	return
}

func (ur *UMSRest) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var credInput models.UserInputLogin
	err := c.BindJSON(&credInput)
	if err != nil {
		errRes := errModels.ErrorResponse{
			Message:              "Invalid request body",
			RecommendationAction: []string{"Check user credentials. Should have username and password"},
			ErrorStatusCode:      http.StatusBadRequest,
		}
		fmt.Printf("Error while binding json to user creds %s", err)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	userCredsIsAdmin, errValidate := ur.UserService.GetAndValidateCredentials(ctx, credInput)
	if errValidate != nil {
		errMsg := fmt.Sprintf("Error validating user %s", errValidate.Message)
		errCode := errValidate.ErrorStatusCode
		if errCode == 0 {
			errCode = http.StatusInternalServerError
		}
		errResp := errModels.ErrorResponse{
			Message:         errMsg,
			ErrorStatusCode: errCode,
		}
		c.JSON(errCode, errResp)
		return
	}
	c.JSON(http.StatusOK, userCredsIsAdmin)
	return
}

// GetUser gets the user from db
func (ur *UMSRest) GetUser(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param(constants.UserIDKey)

	userResp, err := ur.UserService.GetUser(ctx, userID)
	if err != nil {
		errRes := errModels.ErrorResponse{
			Message:         fmt.Sprintf("Failed to get user. Error: %v", err.Message),
			ErrorStatusCode: http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}
	c.JSON(http.StatusOK, userResp)
}