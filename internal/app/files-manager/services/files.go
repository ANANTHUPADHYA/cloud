package services

import (
	"context"
	"fmt"
	fileModels "github.com/ANANTHUPADHYA/cloud/internal/app/files-manager/models"
	usrModels "github.com/ANANTHUPADHYA/cloud/internal/app/user-manager/models"
	"github.com/ANANTHUPADHYA/cloud/internal/app/user-manager/services"
	awss3pkg "github.com/ANANTHUPADHYA/cloud/internal/pkg/aws-s3"
	commonModels "github.com/ANANTHUPADHYA/cloud/internal/pkg/models"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	limitFileDescriptionChars = 400
)

// FileService - holds the functions used for file management
type FileService interface {
	UploadFile(ctx context.Context, userID string, fileInfo fileModels.FileInfo, filereader io.Reader) (usrModels.UserDynamo, *commonModels.ErrorResponse)
	DownloadFile(ctx context.Context, userID string, fileName string) (fileModels.DownloadFileInfo, *commonModels.ErrorResponse)
	DeleteFile(ctx context.Context, userID string, fileName string) (usrModels.UserDynamo, *commonModels.ErrorResponse)
	CheckValidFileName(ctx context.Context, fileName string) (string, error)
	CheckValidDescription(ctx context.Context, description string) *commonModels.ErrorResponse
	UpdateUserFileDescription(ctx context.Context, userID string, updateFiles []string, updateDescription fileModels.UpdateFileInfo) (usrModels.UserDynamo, *commonModels.ErrorResponse)
}

type FileManager struct {
	UserSvc services.UserService
	//UserDBSvc  database.UsersDynamoDBAPI
	AWSS3Svc awss3pkg.IfAWSS3
}

// NewFileService creates an instance of File Service
func NewFileService(userService services.UserService, awsS3Service awss3pkg.IfAWSS3) FileService {
	return &FileManager{
		UserSvc: userService,
		//UserDBSvc: userDBService,
		AWSS3Svc: awsS3Service,
	}
}

//CheckValidFileName validates the file name
func (fm *FileManager) CheckValidFileName(ctx context.Context, fileName string) (string, error) {
	validFileName, err := url.QueryUnescape(fileName)
	if err != nil {
		return validFileName, err
	}
	return validFileName, nil
}

// UploadFile uploads the file to aws s3
func (fm *FileManager) UploadFile(ctx context.Context, userID string, fileInfo fileModels.FileInfo, f io.Reader) (usrModels.UserDynamo, *commonModels.ErrorResponse) {
	log.Printf("User ID %s", userID)
	user := usrModels.UserDynamo{}

	user, err := fm.UserSvc.GetAndValidateUser(ctx, userID)
	if err != nil {
		log.Printf("error occure %v   and   %v", user, err)
		return user, err
	}
	log.Printf("no error %v   and   %v", user, err)
	uploadErr := fm.AWSS3Svc.UploadAttachmentTOS3Bucket(ctx, userID, fileInfo.FileName, f)
	if uploadErr != nil {
		return user, &commonModels.ErrorResponse{
			Message:         fmt.Sprintf("Error while uploading the file. %s", uploadErr.Error()),
			ErrorStatusCode: http.StatusInternalServerError,
		}
	}

	userFileMap := make(map[string]fileModels.FileInfo)
	if user.FileInfo == nil {
		user.FileInfo = userFileMap
	}
	//if _, ok := user.FileInfo[fileInfo.FileName]; ok {
	user.FileInfo[fileInfo.FileName] = fileInfo
	//} else {
	//	user.FileInfo[fileInfo.FileName] = fileInfo
	//}
	return user, nil
}

// DownloadFile returns the presigned URL for downloading the attachment
func (fm *FileManager) DownloadFile(ctx context.Context, userID string, fileName string) (fileModels.DownloadFileInfo, *commonModels.ErrorResponse) {
	downloadAttachmentInfo := fileModels.DownloadFileInfo{}

	user, err := fm.UserSvc.GetAndValidateUser(ctx, userID)
	if err != nil {
		return downloadAttachmentInfo, err
	}

	userFileInfoMap := make(map[string]fileModels.FileInfo)
	if user.FileInfo == nil {
		user.FileInfo = userFileInfoMap
	}
	if _, ok := user.FileInfo[fileName]; !ok {
		return downloadAttachmentInfo, &commonModels.ErrorResponse{
			Message:              fmt.Sprintf("File %s not for user %s", fileName, userID),
			RecommendationAction: []string{"Ensure that the file name is correct"},
			ErrorStatusCode:      http.StatusBadRequest,
		}
	}
	downloadAttachmentInfo, signErr := fm.AWSS3Svc.GenerateS3PresignedURL(ctx, userID, fileName)
	if signErr != nil {
		return downloadAttachmentInfo, &commonModels.ErrorResponse{
			Message:         fmt.Sprintf("Error while getting presigned URL for the attachment %s. Error: %s", fileName, signErr.Error()),
			ErrorStatusCode: http.StatusInternalServerError,
		}
	}
	return downloadAttachmentInfo, nil
}

// CheckValidDescription checks for valid description
func (fm *FileManager) CheckValidDescription(ctx context.Context, description string) *commonModels.ErrorResponse {
	if description == "" {
		errRes := &commonModels.ErrorResponse{
			Message:              "User Attachment description is empty. Cannot update empty attachment description ",
			RecommendationAction: []string{"Expected non empty user attachment description"},
			ErrorStatusCode:      http.StatusBadRequest,
		}
		return errRes
	}
	if len(description) > limitFileDescriptionChars {
		recommendedActions := fmt.Sprintf("File description should be less than %v characters", limitFileDescriptionChars)
		errRes := &commonModels.ErrorResponse{
			Message:              "File description is too lengthy. Cannot update file description ",
			RecommendationAction: []string{recommendedActions},
			ErrorStatusCode:      http.StatusBadRequest,
		}
		return errRes
	}
	return nil
}

// UpdateFileAttachmentDescription updates the file attachments description
func (fm *FileManager) UpdateFileAttachmentDescription(ctx context.Context, userID string, updateFiles []string, updateDescription fileModels.UpdateFileInfo) (usrModels.UserDynamo, *commonModels.ErrorResponse) {
	user, err := fm.UserSvc.GetAndValidateUser(ctx, userID)
	if err != nil {
		return user, err
	}
	userFileAttchmentMap := make(map[string]fileModels.FileInfo)
	if user.FileInfo == nil {
		user.FileInfo = userFileAttchmentMap
	}
	var updateFileName string
	for _, updateFileName = range updateFiles {
		if _, ok := user.FileInfo[updateFileName]; !ok {
			return user, &commonModels.ErrorResponse{
				Message:              fmt.Sprintf("File %s not found for user %s", updateFileName, userID),
				RecommendationAction: []string{"Check for file name"},
				ErrorStatusCode:      http.StatusBadRequest,
			}
		}
		fileInfo := user.FileInfo[updateFileName]
		if len(updateFiles) == 1 && fileInfo.Description == updateDescription.Description {
			return user, &commonModels.ErrorResponse{
				Message:              fmt.Sprintf("Error updating file %s", updateFileName),
				RecommendationAction: []string{"Check file description"},
				ErrorStatusCode:      http.StatusBadRequest,
			}
		}
		fileInfo.Description = updateDescription.Description
		user.FileInfo[updateFileName] = fileInfo
	}
	return user, nil
}

// UpdateUserFileDescription updates the user file description
func (fm *FileManager) UpdateUserFileDescription(ctx context.Context, userID string, updateFiles []string, updateDescription fileModels.UpdateFileInfo) (usrModels.UserDynamo, *commonModels.ErrorResponse) {
	user, err := fm.UserSvc.GetAndValidateUser(ctx, userID)
	if err != nil {
		return user, err
	}
	userFileMap := make(map[string]fileModels.FileInfo)
	if user.FileInfo == nil {
		user.FileInfo = userFileMap
	}
	var updateFileName string
	for _, updateFileName = range updateFiles {
		if _, ok := user.FileInfo[updateFileName]; !ok {
			return user, &commonModels.ErrorResponse{
				Message:              fmt.Sprintf("Error file %s of user %s not found", updateFileName, userID),
				RecommendationAction: []string{"Check for attachment name"},
				ErrorStatusCode:      http.StatusBadRequest,
			}
		}
		fileInfo := user.FileInfo[updateFileName]
		if len(updateFiles) == 1 && fileInfo.Description == updateDescription.Description {
			return user, &commonModels.ErrorResponse{
				Message:              fmt.Sprintf("Error in file name %s", updateFileName),
				RecommendationAction: []string{"Check file name"},
				ErrorStatusCode:      http.StatusBadRequest,
			}
		}
		fileInfo.Description = updateDescription.Description
		user.FileInfo[updateFileName] = fileInfo
	}
	return user, nil
}

// DeleteFile deletes the attachment in s3
func (fm *FileManager) DeleteFile(ctx context.Context, userID string, fileName string) (usrModels.UserDynamo, *commonModels.ErrorResponse) {
	userDB, err := fm.UserSvc.GetAndValidateUser(ctx, userID)
	if err != nil {
		return userDB, err
	}
	userAttchmentMap := make(map[string]fileModels.FileInfo)
	if userDB.FileInfo == nil {
		userDB.FileInfo = userAttchmentMap
	}

	if _, ok := userDB.FileInfo[fileName]; !ok {
		return userDB, &commonModels.ErrorResponse{
			Message:              fmt.Sprintf("File not found %s for user %s", fileName, userID),
			RecommendationAction: []string{"Check the file name"},
			ErrorStatusCode:      http.StatusBadRequest,
		}
	}
	delErr := fm.AWSS3Svc.DeleteFileInS3(ctx, userID, fileName)
	if delErr != nil {
		return userDB, &commonModels.ErrorResponse{
			Message:         fmt.Sprintf("Error while deleting file %s in S3 . Error: %s", fileName, delErr.Error()),
			ErrorStatusCode: http.StatusInternalServerError,
		}
	}
	userDB, err = fm.UserSvc.GetAndValidateUser(ctx, userID)
	if err != nil {
		return userDB, err
	}
	delete(userDB.FileInfo, fileName)

	return userDB, nil
}
