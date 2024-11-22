package validator

import (
	"context"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

func ValidateStruct(ctx context.Context, s interface{}) error {
	return Validate.StructCtx(ctx, s)
}
