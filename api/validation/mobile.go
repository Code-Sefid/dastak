package validation

import (
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
)

func IranianMobileNumberValidator(fld validator.FieldLevel) bool {

	value, ok := fld.Field().Interface().(string)
	if !ok {
		return false
	}

	return IranianMobileNumberValidate(value)
}

const iranianMobileNumberPattern string = `^(09|09[0-3,9])[0-9]{8}$`

func IranianMobileNumberValidate(mobileNumber string) bool {
	res, err := regexp.MatchString(iranianMobileNumberPattern, mobileNumber)
	if err != nil {
		log.Print(err.Error())
	}
	return res
}
