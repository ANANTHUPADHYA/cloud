package v1

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/cloud/internal/app/files-manager/constants"
	fileModels "github.com/cloud/internal/app/files-manager/models"
	"github.com/cloud/internal/app/files-manager/services"
	userConsts "github.com/cloud/internal/app/user-manager/constants"
	userModels "github.com/cloud/internal/app/user-manager/models"
	usrSvc "github.com/cloud/internal/app/user-manager/services"
	"github.com/cloud/internal/pkg/models"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	indexHashKey             = ":user"
)

// FilesRouter holds the dependencies for files router
type FilesRouter struct {
	FileService services.FileService
	UserService usrSvc.UserService
}

// CreateFileRouter return a routing object
func CreateFileRouter(
	fileService services.FileService,
	userService usrSvc.UserService,
) *FilesRouter {
	return &FilesRouter{
		FileService: fileService,
		UserService: userService,
	}
}

func (fr *FilesRouter) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()
	dynamoQueryparams := &dynamodb.QueryInput{
		TableName:              aws.String(userConsts.UsersTableName),
		KeyConditionExpression: aws.String("#skey = :user"),
		ExpressionAttributeNames: map[string]*string{
			"#skey": aws.String("skey"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			indexHashKey: {
				S: aws.String(userConsts.TypeUsersForSortKey),
			},
		},
	}
	dbQuery := models.DatabaseQuery{
		QueryParams: dynamoQueryparams,
	}
	usersResp, err := fr.UserService.GetUsers(ctx, dbQuery)
	if err != nil {
		errRes := models.ErrorResponse{
			Message:         fmt.Sprintf("Failed to get users. Error: %s", err.Message),
			ErrorStatusCode: err.ErrorStatusCode,
		}
		c.JSON(err.ErrorStatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, usersResp)
}

func (fr *FilesRouter) GetAllFiles(c *gin.Context) {
	ctx := c.Request.Context()
	dynamoQueryparams := &dynamodb.QueryInput{
		TableName:              aws.String(userConsts.UsersTableName),
		KeyConditionExpression: aws.String("#skey = :user"),
		ExpressionAttributeNames: map[string]*string{
			"#skey": aws.String("skey"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			indexHashKey: {
				S: aws.String(userConsts.TypeUsersForSortKey),
			},
		},
	}
	dbQuery := models.DatabaseQuery{
		QueryParams: dynamoQueryparams,
	}
	usersResp, err := fr.UserService.GetUsers(ctx, dbQuery)
	if err != nil {
		errRes := models.ErrorResponse{
			Message:         fmt.Sprintf("Failed to get users. Error: %s", err.Message),
			ErrorStatusCode: err.ErrorStatusCode,
		}
		c.JSON(err.ErrorStatusCode, errRes)
		return
	}
	filesResp := []map[string]fileModels.FileInfo{}
	for _, value := range usersResp.Members {
		filesResp = append(filesResp, value.FileInfo)
	}

	c.JSON(http.StatusOK, filesResp)
}

// UploadFile uploads a file attachment to aws s3 bucket
func (fr *FilesRouter) UploadFile(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param(constants.UserIDKey)
	fileHeader, err := c.FormFile(constants.FileKey)
	if err != nil {
		errRes := models.ErrorResponse{
			Message:         fmt.Sprintf("Failed to get form data from key file. Error: %v", err),
			ErrorStatusCode: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		errRes := models.ErrorResponse{
			Message:         fmt.Sprintf("Failed to open the file. Error: %v", err),
			ErrorStatusCode: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	fileName := fileHeader.Filename
	fileName, err = fr.FileService.CheckValidFileName(ctx, fileName)
	if err != nil {
		errMsg := fmt.Sprintf("Please check the File name. Error: %v", err)
		errRes := models.ErrorResponse{
			Message:              errMsg,
			ErrorStatusCode:      http.StatusBadRequest,
			RecommendationAction: []string{"Please provide file name in alphanumeric format"},
		}
		c.JSON(errRes.ErrorStatusCode, errRes)
		return
	}

	fileSizeBytes := fileHeader.Size
	fileSize := CalculateFileSize(fileSizeBytes)
	log.Printf("File size %s", fileSize)
	if fileSize > 10 {
		errFileRes := models.ErrorResponse{
			Message: fmt.Sprintf("File size (%d MB) greater than 10 MB.", fileSize),
			RecommendationAction: []string{fmt.Sprint("File size should be lesser than 10 MB")},
			ErrorStatusCode: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errFileRes)
	}
	createdAt := time.Now().Format(time.RFC3339)
	fileInfo := fileModels.FileInfo{FileName: fileName, UpdatedAt: createdAt, CreatedAt: createdAt}
	user, errResp := fr.FileService.UploadFile(ctx, userID, fileInfo, f)
	if errResp != nil {
		errRes := models.ErrorResponse{
			Message:         fmt.Sprintf("Failed to upload file %s for User %s. Error: %v", fileName, userID, errResp),
			ErrorStatusCode: http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	updatedUser, errResp := fr.UserService.UpdateUser(ctx, user)
	if errResp != nil {
		errRes := models.ErrorResponse{
			Message:         fmt.Sprintf("Upload successful but failed to update user in DB. Error: %s", errResp.Message),
			ErrorStatusCode: errResp.ErrorStatusCode,
		}
		c.JSON(errResp.ErrorStatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// UpdateFileDescription updates a file description
func (fr *FilesRouter) UpdateFileDescription(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param(constants.UserIDKey)
	queryMap := c.Request.URL.Query()
	if _, ok := queryMap[constants.FileKey]; !ok {
		errRes := models.ErrorResponse{
			Message:         fmt.Sprintf("Failed to update file description for User %s. Expected file key in query.", userID),
			ErrorStatusCode: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	var updatefileName string
	updateFiles := queryMap[constants.FileKey]
	validFileList := []string{}
	for _, updatefileName = range updateFiles {
		if updatefileName == "" {
			errRes := models.ErrorResponse{
				Message:         fmt.Sprintf("Failed to update file description for User %s. Expected file name in query.", userID),
				ErrorStatusCode: http.StatusBadRequest,
			}
			c.JSON(http.StatusBadRequest, errRes)
			return
		}
		validFileName, err := fr.FileService.CheckValidFileName(ctx, updatefileName)
		if err != nil {
			errMsg := fmt.Sprintf("Please check the File name. Error: %v", err)
			errRes := models.ErrorResponse{
				Message:              errMsg,
				ErrorStatusCode:      http.StatusBadRequest,
				RecommendationAction: []string{"Please provide file name in alphanumeric format"},
			}
			c.JSON(errRes.ErrorStatusCode, errRes)
			return
		}
		validFileList = append(validFileList, validFileName)
	}
	updateDescription := fileModels.UpdateFileInfo{}
	body, error := ioutil.ReadAll(c.Request.Body)
	if error != nil {
		errRes := models.ErrorResponse{
			Message:              fmt.Sprintf("Invalid request body. Error: %s", error.Error()),
			RecommendationAction: []string{"Check file description update request body"},
			ErrorStatusCode:      http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	error = json.Unmarshal(body, &updateDescription)
	if error != nil {
		msg := fmt.Sprintf("Error while unmarshalling the body with updateUserFileInfo model. Error: %s", error.Error())
		errRes := models.ErrorResponse{
			Message:              "Request body json validation issue. " + msg,
			RecommendationAction: []string{"Check json request body"},
			ErrorStatusCode:      http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}
	checkResp := fr.FileService.CheckValidDescription(ctx, updateDescription.Description)
	if checkResp != nil {
		checkErrResp := models.ErrorResponse{
			Message:              checkResp.Message,
			RecommendationAction: checkResp.RecommendationAction,
			ErrorStatusCode:      checkResp.ErrorStatusCode,
		}
		c.JSON(checkErrResp.ErrorStatusCode, checkErrResp)
		return
	}
	updatedFile, updateErrResp := fr.FileService.UpdateUserFileDescription(ctx, userID, validFileList, updateDescription)
	if updateErrResp != nil {
		updateRes := models.ErrorResponse{
			Message:              fmt.Sprintf("unable to update the file description. Error: %s", updateErrResp.Message),
			RecommendationAction: updateErrResp.RecommendationAction,
			ErrorStatusCode:      updateErrResp.ErrorStatusCode,
		}
		c.JSON(updateErrResp.ErrorStatusCode, updateRes)
		return
	}
	updatedUserResp, errResp := fr.UserService.UpdateUser(ctx, updatedFile)
	if errResp != nil {
		errRes := models.ErrorResponse{
			Message:         fmt.Sprintf("Failed to update file description in DB. Error: %s", errResp.Message),
			ErrorStatusCode: errResp.ErrorStatusCode,
		}
		c.JSON(errResp.ErrorStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, updatedUserResp)
}

// DownloadFile downloads a file from aws s3 bucket by presignedURL
func (fr *FilesRouter) DownloadFile(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param(constants.UserIDKey)
	queryMap := c.Request.URL.Query()
	var fileName string
	if _, ok := queryMap[constants.FileKey]; !ok {
		errRes := models.ErrorResponse{
			Message:         fmt.Sprintf("Failed to download file for user %s. Expected file key in query.", userID),
			ErrorStatusCode: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	fileName = queryMap[constants.FileKey][0]
	if fileName == "" {
		errRes := models.ErrorResponse{
			Message:         fmt.Sprintf("Failed to download file for user %s. Expected file name in query.", userID),
			ErrorStatusCode: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	presignedURL, errResp := fr.FileService.DownloadFile(ctx, userID, fileName)
	if errResp != nil {
		errRes := models.ErrorResponse{
			Message:              fmt.Sprintf("unable to download the file attachment. Error: %s", errResp.Message),
			RecommendationAction: errResp.RecommendationAction,
			ErrorStatusCode:      errResp.ErrorStatusCode,
		}
		c.JSON(errResp.ErrorStatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, presignedURL)
}

// DeleteFile deletes a file attachment from aws s3 bucket
func (fr *FilesRouter) DeleteFile(c *gin.Context) {
	ctx := c.Request.Context()
	delUserID := c.Param(constants.UserIDKey)
	delQueryMap := c.Request.URL.Query()
	var delFileName string
	if _, ok := delQueryMap[constants.FileKey]; !ok {
		errRes := models.ErrorResponse{
			Message:         fmt.Sprintf("Failed to delete file for User %s. Expected file key in query.", delUserID),
			ErrorStatusCode: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	delFiles := delQueryMap[constants.FileKey]
	validFileList := []string{}
	for _, delFileName = range delFiles {
		if delFileName == "" {
			errRes := models.ErrorResponse{
				Message:         fmt.Sprintf("Failed to delete file for User %s. Expected file name in query.", delUserID),
				ErrorStatusCode: http.StatusBadRequest,
			}
			c.JSON(http.StatusBadRequest, errRes)
			return
		}
		validFileList = append(validFileList, delFileName)
	}
	var updatedUser userModels.UserDynamo
	var errResp *models.ErrorResponse
	for _, delFileName := range validFileList {
		updatedUser, delErrResp := fr.FileService.DeleteFile(ctx, delUserID, delFileName)
		if delErrResp != nil {
			delRes := models.ErrorResponse{
				Message:              fmt.Sprintf("unable to delete the quote attachment %s. Error: %s", delFileName, delErrResp.Message),
				RecommendationAction: delErrResp.RecommendationAction,
				ErrorStatusCode:      delErrResp.ErrorStatusCode,
			}
			c.JSON(delErrResp.ErrorStatusCode, delRes)
			return
		}
		updatedUser, errResp = fr.UserService.UpdateUser(ctx, updatedUser)
		if errResp != nil {
			errRes := models.ErrorResponse{
				Message:         fmt.Sprintf("Delete file is successful but failed to update user in DB. Error: %s", errResp.Message),
				ErrorStatusCode: errResp.ErrorStatusCode,
			}
			c.JSON(errResp.ErrorStatusCode, errRes)
			return
		}
	}
	c.JSON(http.StatusOK, updatedUser.User)
}

//CalculateFileSize returns uploaded file size
func CalculateFileSize(fileSizeBytes int64) int64 {
	fileSizeInKB := fileSizeBytes / 1024
	fileSizeInMB := fileSizeInKB / 1024
	return fileSizeInMB
}
