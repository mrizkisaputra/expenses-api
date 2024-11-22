package server

import (
	"github.com/gin-gonic/gin"
	expenseController "github.com/mrizkisaputra/expenses-api/internal/expense/controllers/http"
	expenseRoute "github.com/mrizkisaputra/expenses-api/internal/expense/controllers/http"
	expenseRepository "github.com/mrizkisaputra/expenses-api/internal/expense/repository"
	expenseService "github.com/mrizkisaputra/expenses-api/internal/expense/service"
	"github.com/mrizkisaputra/expenses-api/internal/middleware"
	userController "github.com/mrizkisaputra/expenses-api/internal/user/controllers/http"
	userRoute "github.com/mrizkisaputra/expenses-api/internal/user/controllers/http"
	userRepository "github.com/mrizkisaputra/expenses-api/internal/user/repository"
	userService "github.com/mrizkisaputra/expenses-api/internal/user/service"
	"net/http"
)

func (s *Server) Bootstrap() error {
	// -----------------------------------------------------------------------------------------------------------
	// create a new instance repositories
	userPostgresRepo := userRepository.NewUserPostgresRepository(s.db)
	userRedisRepo := userRepository.NewUserRedisRepository(s.redisClient)
	userAwsRepo := userRepository.NewAWSUserRepository(s.awsClient)

	expenseRepo := expenseRepository.NewExpensePgRepository(s.db)

	// -----------------------------------------------------------------------------------------------------------
	// create a new instance services
	authSV := userService.NewAuthService(&userService.ServiceConfig{
		Config:                 s.cfg,
		UserPostgresRepository: userPostgresRepo,
	})
	userSV := userService.NewUserService(&userService.ServiceConfig{
		Config:                 s.cfg,
		Logger:                 s.logger,
		UserPostgresRepository: userPostgresRepo,
		UserRedisRepository:    userRedisRepo,
		AwsUserRepository:      userAwsRepo,
	})

	expenseSV := expenseService.NewExpenseService(&expenseService.ServiceConfig{
		PgRepo: expenseRepo,
		Logger: s.logger,
		Config: s.cfg,
	})

	// -----------------------------------------------------------------------------------------------------------
	// create a new instance controllers
	authController := userController.NewAuthController(&userController.ControllerConfig{
		Config:      s.cfg,
		Logger:      s.logger,
		AuthService: authSV,
	})
	usrController := userController.NewUserController(&userController.ControllerConfig{
		Config:      s.cfg,
		Logger:      s.logger,
		UserService: userSV,
	})

	expenseCntrl := expenseController.NewExpenseController(&expenseController.ControllerConfig{
		ExpenseService: expenseSV,
		Logger:         s.logger,
	})

	// -----------------------------------------------------------------------------------------------------------
	// create a new instance middleware
	middlewareManager := middleware.NewMiddlewareManager(&middleware.MiddlewareConfig{
		Logger: s.logger,
		Config: s.cfg,
	})

	s.app.Use(middlewareManager.RequestIdMiddleware())
	s.app.Use(middlewareManager.RequestLoggerMiddleware())

	// -----------------------------------------------------------------------------------------------------------
	// setup routes
	apiV1 := s.app.Group("/api/v1")
	{
		// group user routes
		userGroup := apiV1.Group("/user")
		{
			userRoute.MapAuthRoutes(userGroup, authController)
			userRoute.MapUserRoutes(userGroup, usrController, middlewareManager)
		}

		// group expense routes
		expenseGroup := apiV1.Group("/expenses")
		{
			expenseRoute.MapExpenseRoutes(expenseGroup, expenseCntrl, middlewareManager)
		}
	}

	apiV1.GET("/ping", middlewareManager.AuthJwtMiddleware(), func(ctx *gin.Context) {
		auth := middleware.GetAuth(ctx)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Pong!!",
			"user_id": auth.Id,
		})
	})

	return nil
}
