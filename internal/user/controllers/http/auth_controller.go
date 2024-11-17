package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/mrizkisaputra/expenses-api/internal/user"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/mrizkisaputra/expenses-api/internal/user/model/dto"
	"github.com/mrizkisaputra/expenses-api/pkg/httpErrors"
	"github.com/mrizkisaputra/expenses-api/pkg/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

// authController acts as a struct for injecting an implementation of AuthController interface
// for use in controller methods
type authController struct {
	cfg     *config.Config
	logger  *logrus.Logger
	service user.AuthService
}

// NewAuthController is a factory function
// initializing a authController with its service layer dependencies
func NewAuthController(config *ControllerConfig) user.AuthController {
	return &authController{
		cfg:     config.Config,
		logger:  config.Logger,
		service: config.AuthService,
	}
}

func (ac *authController) RegisterNewUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// read and validate request
		request := new(dto.UserRegisterRequest)
		if err := utils.ReadRequest(ctx, request, binding.JSON); err != nil {
			utils.LogErrorResponse(ctx, ac.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		// register new user
		entityUser := &model.User{
			Email:    request.Email,
			Password: request.Password,
			Information: model.Information{
				FirstName: request.FirstName,
				LastName:  request.LastName,
			},
		}
		userResponse, err := ac.service.Register(context.Background(), entityUser)
		if err != nil {
			utils.LogErrorResponse(ctx, ac.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		// return success response
		ctx.JSON(http.StatusCreated, dto.ApiUserResponse{
			Status:  http.StatusCreated,
			Message: "Created",
			Data:    userResponse,
		})
	}
}

func (ac *authController) LoginNewUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// read and validate request
		request := new(dto.UserLoginRequest)
		if err := utils.ReadRequest(ctx, request, binding.JSON); err != nil {
			utils.LogErrorResponse(ctx, ac.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		// login
		jwtToken, err := ac.service.Login(context.Background(), &model.User{
			Email:    request.Email,
			Password: request.Password,
		})
		if err != nil {
			utils.LogErrorResponse(ctx, ac.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			return
		}

		// return response token
		ctx.JSON(http.StatusOK, dto.UserTokenResponse{
			Status:  http.StatusOK,
			Message: "OK",
			Jwt:     *jwtToken,
		})
	}
}
