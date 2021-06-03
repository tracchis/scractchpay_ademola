package validatorutil

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var uni *ut.UniversalTranslator
var trans ut.Translator

func GetValidator() *validator.Validate {
	en := en.New()
	uni = ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ = uni.GetTranslator("en")

	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	validate.SetTagName("validate")
	// Print JSON name on validator.FieldError.Field()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		return name
	})

	return validate
}

func GetTranslatedErrors(err error) map[string]string {
	errs := err.(validator.ValidationErrors)
	errorMessages := errs.Translate(trans)
	return errorMessages
}

func GetAttributeErrorMessages() map[string]string {
	attrErrorMessages := map[string]string{}
	return attrErrorMessages
}
