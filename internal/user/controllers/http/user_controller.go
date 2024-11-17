package http

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/mrizkisaputra/expenses-api/internal/middleware"
	"github.com/mrizkisaputra/expenses-api/internal/user"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/mrizkisaputra/expenses-api/internal/user/model/dto"
	"github.com/mrizkisaputra/expenses-api/pkg/httpErrors"
	"github.com/mrizkisaputra/expenses-api/pkg/utils"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// userController acts as a struct for injecting an implementation of UserController interface
// for use in controller methods
type userController struct {
	cfg     *config.Config
	logger  *logrus.Logger
	service user.UserService
}

// NewUserController is a factory function
// initializing a authController with its service layer dependencies
func NewUserController(config *ControllerConfig) user.UserController {
	return &userController{
		cfg:     config.Config,
		logger:  config.Logger,
		service: config.UserService,
	}
}

func (u *userController) GetCurrentUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := middleware.GetAuth(ctx)

		userResponse, err := u.service.GetCurrentUser(context.Background(), auth.Id.String())
		if err != nil {
			utils.LogErrorResponse(ctx, u.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		ctx.JSON(http.StatusOK, &dto.ApiUserResponse{
			Status:  http.StatusOK,
			Message: "OK",
			Data:    userResponse,
		})
	}
}

func (u *userController) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := middleware.GetAuth(ctx)

		// read and validate update payload body request
		request := new(dto.UserUpdateRequest)
		if err := utils.ReadRequest(ctx, request, binding.JSON); err != nil {
			utils.LogErrorResponse(ctx, u.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		userResponse, err := u.service.Update(context.Background(), &model.User{
			Id:    auth.Id,
			Email: request.Email,
			Information: model.Information{
				FirstName:   request.FirstName,
				LastName:    request.LastName,
				City:        request.City,
				PhoneNumber: request.PhoneNumber,
			},
		})

		if err != nil {
			utils.LogErrorResponse(ctx, u.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		ctx.JSON(http.StatusAccepted, &dto.ApiUserResponse{
			Status:  http.StatusAccepted,
			Message: "Accept",
			Data:    userResponse,
		})
	}
}

func (u *userController) PostAvatar() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := middleware.GetAuth(ctx)

		// read request query parameter
		bucketName, exist := ctx.GetQuery("bucket")
		if !exist {
			err := httpErrors.NewBadRequestError("query bucket param is required")
			utils.LogErrorResponse(ctx, u.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		// read and validate request post file
		image, err := utils.ReadImageRequest(ctx, "avatar")
		if err != nil {
			utils.LogErrorResponse(ctx, u.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		file, err := image.Open()
		defer file.Close()
		if err != nil {
			utils.LogErrorResponse(ctx, u.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		//
		binaryImage := bytes.NewBuffer(nil)
		if _, err := io.Copy(binaryImage, file); err != nil {
			utils.LogErrorResponse(ctx, u.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		// post avatar
		reader := bytes.NewReader(binaryImage.Bytes())
		userResponse, err := u.service.UploadAvatar(context.Background(), auth.Id, &model.UserUploadInput{
			Object:      reader,
			ObjectName:  image.Filename,
			ObjectSize:  image.Size,
			BucketName:  bucketName,
			ContentType: image.Header.Get("Content-Type"),
		})

		ctx.JSON(http.StatusCreated, &dto.ApiUserResponse{
			Status:  http.StatusCreated,
			Message: "Created",
			Data:    userResponse.Avatar,
		})
	}
}
