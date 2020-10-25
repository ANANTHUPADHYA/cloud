package main

import (
	awss3 "github.com/ANANTHUPADHYA/cloud/internal/pkg/aws-s3"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	fileMgHndlr "github.com/ANANTHUPADHYA/cloud/internal/app/files-manager/handlers/v1"
	fileSvc "github.com/ANANTHUPADHYA/cloud/internal/app/files-manager/services"
	userMgHndlr "github.com/ANANTHUPADHYA/cloud/internal/app/user-manager/handlers/v1"
	userSvc "github.com/ANANTHUPADHYA/cloud/internal/app/user-manager/services"
	"github.com/ANANTHUPADHYA/cloud/internal/pkg/database"
	"github.com/ANANTHUPADHYA/cloud/internal/pkg/http/transport"
	"github.com/ANANTHUPADHYA/cloud/internal/pkg/utils"
	 cors "github.com/rs/cors/wrapper/gin"
)

var (
	port            = utils.GetEnvOrDefault("PORT", "3000")
	healthzEndpoint = "/healthz"
)

func main() {
	//create a router and corresponding groups
	router := gin.New()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())
	router.Use(cors.AllowAll())

	//healthz endpoint
	router.GET(healthzEndpoint, healthzCheck)

	umsV1 := router.Group("/v1")

	loggerTransport := transport.NewLoggerTransport(http.DefaultTransport)
	metricTransport := transport.NewMetricTransport(loggerTransport)

	// adding timeout explicitly to limit the wait time for each client call
	// keeping it as 15 seconds to avoid the request timing out too soon
	dbClient := &http.Client{
		Timeout:   15 * time.Second,
		Transport: metricTransport,
	}
	awsSvcSession := database.NewAWSCredsImpl()
	dynamoDBsvc, err := awsSvcSession.GetDynamodbSVC(dbClient)
	if err != nil {
		panic(err)
	}

	s3creds := awss3.NewAWSCredsImpl()
	s3Impl := awss3.NewAWSS3Impl(s3creds)
	s3IfImpl, err := s3Impl.GetS3SVC()
	if err != nil {
		panic(err)
	}
	s3Svc := awss3.NewAWSS3(s3IfImpl, s3creds)

	usersDBImpl := database.NewUsersDBImpl(dynamoDBsvc)
	userService := userSvc.NewUserService(&usersDBImpl)
	usersRouter := userMgHndlr.CreateUMSRouter(userService)
	log.Print("Starting my service")
	umsV1.POST("/users",
		usersRouter.CreateUser,
	)

	umsV1.PUT("/login",
		usersRouter.Login,
	)

	umsV1.GET(
		"/users/:user_id",
		usersRouter.GetUser,
	)

	filev1 := router.Group("/v1")
	fileService := fileSvc.NewFileService(userService, s3Svc)
	filesRouter := fileMgHndlr.CreateFileRouter(fileService, userService)

	filev1.GET(
		"/users",
		filesRouter.GetAllUsers)

	filev1.GET(
		"/files",
		filesRouter.GetAllFiles)

	filev1.PUT(
		"/users/:user_id/upload",
		filesRouter.UploadFile,
	)

	filev1.PATCH(
		"/users/:user_id/file-update",
		filesRouter.UpdateFileDescription,
	)

	filev1.GET(
		"/users/:user_id/download",
		filesRouter.DownloadFile,
	)

	filev1.GET(
		"/admin/users/:user_id/download",
		filesRouter.DownloadFile,
	)

	filev1.DELETE(
		"/users/:user_id/file",
		filesRouter.DeleteFile,
	)

	filev1.DELETE(
		"/admin/users/:user_id/file",
		filesRouter.DeleteFile,
	)

	//log.Infof(context.Background(), "Listening on %v", port)

	err = router.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

//healthzCheck returns the health check status of the user service
func healthzCheck(c *gin.Context) {
	//as of now sending ok, in future this may send the
	//health status of other services as well
	c.JSON(http.StatusOK, "OK")
}
