package validation

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Message string `json:"message"`
}

var fieldTranslations = map[string]string{
	"FullName":    "نام کامل",
	"AccountType": "نوع حساب",
	"SaleCount":   "تعداد فروش",
	"Mobile":      "شماره موبایل",
	"Password":    "رمز عبور",
	"Name":        "نام",
}

func translateField(fieldName string) string {
	if translation, ok := fieldTranslations[fieldName]; ok {
		return translation
	}
	return fieldName
}

func TranslateValidationError(tag string) string {
	switch tag {
	case "required":
		return "اجباری است"
	case "min":
		return "کمتر از حد مجاز است"
	case "max":
		return "بیشتر از حد مجاز است"
	case "len":
		return "باید دقیقاً طول مشخص شده داشته باشد"
	case "mobile":
		return "شماره موبایل معتبر نیست"
	default:
		return "خطا در اعتبارسنجی"
	}
}

func GetValidationErrors(err error) []*string {
	var validationErrors []*string
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := translateField(err.Field()) + " " + TranslateValidationError(err.Tag())
			validationErrors = append(validationErrors, &errorMessage)
		}
		return validationErrors
	}
	return nil
}
