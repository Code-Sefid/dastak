package helper

import (
	"github.com/soheilkhaledabdi/dastak/api/validation"
)

type BaseHttpResponse struct {
	Result           any       `json:"result"`
	Status           bool      `json:"status"`
	Alert            *string   `json:"alert"`
	ValidationErrors []*string `json:"validationErrors"`
}

func GenerateAlert(message string) *string {
	if message == "" {
		return nil
	}
	return &message
}

func GenerateBaseResponse(result any, Status bool, message string) *BaseHttpResponse {
	return &BaseHttpResponse{
		Result: result,
		Status: Status,
		Alert:  GenerateAlert(message),
	}
}

func GenerateBaseResponseWithError(Status bool, err any, message string) *BaseHttpResponse {
	return &BaseHttpResponse{
		Status: Status,
		Alert:  GenerateAlert(message),
	}
}

func GenerateBaseResponseWithAnyError(Status bool, err any, message string) *BaseHttpResponse {
	return &BaseHttpResponse{
		Alert:  GenerateAlert(message),
		Status: Status,
	}
}

func GenerateBaseResponseWithValidationError(Status bool, err error, message string) *BaseHttpResponse {
	validationErrs := validation.GetValidationErrors(err)
	if validationErrs != nil {
		return &BaseHttpResponse{
			Alert:            GenerateAlert(message),
			Status:           Status,
			ValidationErrors: validationErrs,
		}
	}

	return &BaseHttpResponse{
		Alert:  GenerateAlert(message),
		Status: Status,
	}
}
