package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mrizkisaputra/expenses-api/pkg/httpErrors"
	"github.com/mrizkisaputra/expenses-api/pkg/utils"
	"github.com/sirupsen/logrus"
	"strings"
)

type auth struct {
	Id    uuid.UUID
	Email string
}

// AuthJwtMiddleware is a middleware authentication request.
// cara autentikasi JWT menggunakan cookie atau header otorisasi
func (mw *MiddlewareManager) AuthJwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("Authorization")

		mw.logger.WithFields(logrus.Fields{"Authorization": authorizationHeader}).
			Debug("auth middleware authorization header")

		var tokenString string

		// cara authentikasi menggunakan header authorization
		// cek apakah authorization header ada?
		if authorizationHeader != "" {
			headerParts := strings.Split(authorizationHeader, " ") // ['Bearer','token']
			if len(headerParts) < 2 {
				errResponse := httpErrors.NewUnauthorizedError("MiddlewareManager.AuthJwtMiddleware.Split")
				utils.LogErrorResponse(ctx, mw.logger, errResponse)
				ctx.JSON(httpErrors.ErrorResponse(ctx, errResponse))
				ctx.Abort()
				return
			}
			tokenString = headerParts[1]
		} else {
			// cara authentikasi menggunakan cookie, jika tidak ada authorization header
			cookie, err := ctx.Cookie("jwt-token") // setcookie dengan jwt belum diimplemntasikan
			if err != nil {
				errResponse := httpErrors.NewUnauthorizedError(err)
				utils.LogErrorResponse(ctx, mw.logger, errResponse)
				ctx.JSON(httpErrors.ErrorResponse(ctx, errResponse))
				ctx.Abort()
				return
			}
			tokenString = cookie
		}

		// verifikasi jwt token dan mengambil claims dari token
		claims, err := utils.ValidateJwtToken(tokenString, mw.cfg)
		if err != nil {
			utils.LogErrorResponse(ctx, mw.logger, err)
			ctx.JSON(httpErrors.ErrorResponse(ctx, err))
			ctx.Abort()
			return
		}

		// menyimpan claims kedalam context gin
		auth := &auth{
			Id:    claims.ID,
			Email: claims.Email,
		}
		ctx.Set("auth", auth)

		ctx.Next()
	}
}

func GetAuth(ctx *gin.Context) *auth {
	value, _ := ctx.Get("auth")
	return value.(*auth)
}
