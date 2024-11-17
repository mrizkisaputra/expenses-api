package http

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mrizkisaputra/expenses-api/config"
	mockObject "github.com/mrizkisaputra/expenses-api/internal/user/mock"
	"github.com/mrizkisaputra/expenses-api/internal/user/model/dto"
	"github.com/mrizkisaputra/expenses-api/pkg/httpErrors"
	"github.com/mrizkisaputra/expenses-api/pkg/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var cfg = &config.Config{
	Logger: config.LoggerConfig{Level: "panic"},
}

func TestAuthController_RegisterNewUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// scenario test case #1
	t.Run("[Test Case #1] Success register", func(t *testing.T) {
		mockAuthService := new(mockObject.AuthServiceMock)
		controller := NewAuthController(&ControllerConfig{
			Config:      cfg,
			Logger:      logger.NewLogrusLogger(cfg),
			AuthService: mockAuthService,
		})

		mockUserRegister := &dto.UserRegisterRequest{
			FirstName: "muhammat",
			LastName:  "Saputra",
			Email:     "muhammatsaputra@gmail.com",
			Password:  "secret123456",
		}

		body, err := json.Marshal(mockUserRegister)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		rec := httptest.NewRecorder()

		// create test gin context for running tests
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		handleFunc := controller.RegisterNewUser()

		// mock behavior from service layer
		mockAuthService.On("Register", mock.Anything, mock.Anything).Return(&dto.UserResponse{
			Id:        uuid.New(),
			Email:     "test@test.com",
			Password:  "secret12345",
			FirstName: "test",
			LastName:  "test",
		}, nil)

		handleFunc(ctx)

		responseContentType := rec.Header().Get("Content-Type")
		responseStatusCode := rec.Code
		var responseBody = new(dto.ApiUserResponse)
		er := json.Unmarshal(rec.Body.Bytes(), responseBody)
		require.NoError(t, er)

		require.Equal(t, http.StatusCreated, responseStatusCode)
		require.Equal(t, "application/json; charset=utf-8", responseContentType)
		require.Equal(t, "Created", responseBody.Message)
		require.NotNil(t, responseBody)
		require.NotNil(t, responseBody.Data)
	})

	// scenario test case #2
	t.Run("[Test Case #2] Error when binding request to struct", func(t *testing.T) {
		mockService := new(mockObject.AuthServiceMock)
		controller := NewAuthController(&ControllerConfig{
			Config:      cfg,
			Logger:      logger.NewLogrusLogger(cfg),
			AuthService: mockService,
		})

		body := strings.NewReader("firstName=jonson&lastName=smith&email=jsonson@gmail.com&password=secret123456") // expect 'json'
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // expect 'application/json'
		req.Header.Set("Accept", "application/json")
		rec := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req

		handleFunc := controller.RegisterNewUser()
		handleFunc(ctx)

		responseCode := rec.Code
		var responseBody = new(httpErrors.ApiErrorResponse)
		er := json.Unmarshal(rec.Body.Bytes(), responseBody)
		require.NoError(t, er)

		require.Equal(t, http.StatusBadRequest, responseCode)
		require.NotNil(t, responseBody)
		require.Equal(t, httpErrors.BadRequestErrorMsg, responseBody.ErrorInfo.Message)

	})

	// scenario test case #3
	t.Run("[Test Case #3] Error when validate request", func(t *testing.T) {
		mockAuthService := new(mockObject.AuthServiceMock)
		controller := NewAuthController(&ControllerConfig{
			Config:      cfg,
			Logger:      logger.NewLogrusLogger(cfg),
			AuthService: mockAuthService,
		})

		mockUserRegister := &dto.UserRegisterRequest{
			FirstName: "muhammat23", // invalid
			LastName:  "Saputra",
			Email:     "muhammatsaputra", // invalid
			Password:  "secret",          // invalid
		}

		body, err := json.Marshal(mockUserRegister)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		rec := httptest.NewRecorder()

		// create test gin context for running tests
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		handleFunc := controller.RegisterNewUser()

		handleFunc(ctx)
		responseCode := rec.Code
		var responseBody = new(httpErrors.ApiErrorResponse)
		er := json.Unmarshal(rec.Body.Bytes(), responseBody)
		require.NoError(t, er)

		require.Equal(t, http.StatusBadRequest, responseCode)
		require.Equal(t, httpErrors.BadRequestErrorMsg, responseBody.ErrorInfo.Message)
		require.NotNil(t, responseBody)
		require.NotNil(t, responseBody.ErrorInfo.SubError)
		require.Equal(t, 3, len(*responseBody.ErrorInfo.SubError))
	})
}

func TestAuthController_LoginNewUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// scenario test case #1
	t.Run("[Test Case #1] Success login", func(t *testing.T) {
		mockService := new(mockObject.AuthServiceMock)
		controller := NewAuthController(&ControllerConfig{
			Config:      cfg,
			Logger:      logger.NewLogrusLogger(cfg),
			AuthService: mockService,
		})

		mockUserLogin := &dto.UserLoginRequest{
			Email:    "muhammat@gamil.com",
			Password: "secret12345",
		}

		body, err := json.Marshal(mockUserLogin)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		rec := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = req
		handleFunc := controller.LoginNewUser()

		mockService.On("Login", mock.Anything, mock.Anything).
			Return(&dto.JwtToken{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
			}, nil)

		handleFunc(ctx)

		responseCode := rec.Code
		var responseBody = new(dto.UserTokenResponse)
		er := json.Unmarshal(rec.Body.Bytes(), responseBody)
		require.NoError(t, er)

		require.Equal(t, http.StatusOK, responseCode)
		require.Equal(t, "OK", responseBody.Message)
		require.NotNil(t, responseBody)
		require.NotNil(t, responseBody.Jwt)

	})
}
