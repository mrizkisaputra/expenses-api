package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mrizkisaputra/expenses-api/pkg/contextutils"
	"github.com/mrizkisaputra/expenses-api/pkg/httpErrors"
	. "github.com/mrizkisaputra/expenses-api/pkg/validator"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

// ReadRequest is a function for read and validate request
func ReadRequest(ctx *gin.Context, request any, binding binding.Binding) error {
	// binding request to dto struct
	if err := ctx.ShouldBindWith(request, binding); err != nil {
		return httpErrors.NewBadRequestError(errors.Wrap(err, "ReadRequest.ShouldBindWith"))
	}

	// validate request
	if err := Validate.Struct(request); err != nil {
		return errors.Wrap(err, "ReadRequest.Validate.Struct")
	}

	return nil
}

// ReadImageRequest is a function for read and validate post request file image
func ReadImageRequest(ctx *gin.Context, fieldName string) (*multipart.FileHeader, error) {
	// parse and read file request
	fh, err := ctx.FormFile(fieldName)
	if err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "UserController.PostAvatar.FormFile"))
	}

	// specify maximum file size (in bytes)
	maxFileSize := int64(1 << 20) // 1048576 bytes (1MB)
	if fh.Size > maxFileSize {    // bytes > bytes
		return nil, httpErrors.NewError(http.StatusBadRequest, httpErrors.MaxFileSizeMsg, nil)
	}

	// specify header content-type
	allowedImagesContentType := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
	}
	contentType := fh.Header.Get("Content-Type")
	if !allowedImagesContentType[contentType] {
		return nil, httpErrors.NewError(
			http.StatusBadRequest,
			httpErrors.NotAllowedImageHeaderMsg,
			fmt.Sprintf("content-type '%s' not allowed", contentType),
		)
	}

	// specify file extension
	allowedFileExtensions := map[string]bool{
		".jpeg": true,
		".jpg":  true,
		".png":  true,
	}
	extension := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowedFileExtensions[extension] {
		return nil, httpErrors.NewError(
			http.StatusBadRequest,
			httpErrors.NotAllowedFileExtensionMsg,
			fmt.Sprintf("file extension '%s' not allowed", extension),
		)
	}

	// specify real content type (magic bytes)
	imageFile, err := fh.Open()
	defer imageFile.Close()
	if err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "UserController.ReadImageRequest.Open"))
	}

	bytes := make([]byte, 512)
	if _, err := imageFile.Read(bytes); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "UserController.ReadImageRequest.Read"))
	}

	detectContentType := http.DetectContentType(bytes)
	if !allowedImagesContentType[detectContentType] {
		return nil, httpErrors.NewError(
			http.StatusBadRequest,
			httpErrors.NotAllowedImageHeaderMsg,
			fmt.Sprintf("content-type '%s' not allowed", contentType),
		)
	}

	// validate finish
	return fh, nil
}

// LogErrorResponse is a function for write log response error
func LogErrorResponse(ctx *gin.Context, logger *logrus.Logger, err error) {
	logger.WithError(err).WithFields(logrus.Fields{
		"requestId": contextutils.GetRequestId(ctx),
		"IPAddress": contextutils.GetIPAddress(ctx),
	}).Error("Logger error response")
}
