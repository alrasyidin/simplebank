package api

import (
	"github.com/alrasyidin/simplebank-go/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsValidCurrency(currency)
	}

	return false
}
