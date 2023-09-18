package helper

import (
	"strings"

	"github.com/soheilkhaledabdi/dastak/api/validation"
)

type AlertMessage struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type BaseHttpResponse struct {
	Result           any                           `json:"result"`
	Status           bool                          `json:"status"`
	Alert            *AlertMessage                 `json:"alert"`
	ValidationErrors *[]validation.ValidationError `json:"validationErrors"`
	Error            any                           `json:"error"`
}

func GenerateAlert(title, message string) *AlertMessage {
	if title == "" && message == "" {
		return nil
	}
	return &AlertMessage{
		Title:   title,
		Message: message,
	}
}

func GenerateBaseResponse(result any, Status bool, Title string, Message string) *BaseHttpResponse {
	return &BaseHttpResponse{
		Result: result,
		Status: Status,
		Alert:  GenerateAlert(Title, Message),
	}
}

func GenerateBaseResponseWithError(Status bool, err any, Title string, Message string) *BaseHttpResponse {
	return &BaseHttpResponse{
		Status: Status,
		Error:  err,
		Alert:  GenerateAlert(Title, Message),
	}
}

func GenerateBaseResponseWithAnyError(Status bool, err any, Title string, Message string) *BaseHttpResponse {
	return &BaseHttpResponse{
		Alert:  GenerateAlert(Title, Message),
		Status: Status,
		Error:  err,
	}
}

func GenerateBaseResponseWithValidationError(Status bool, err error, Title string, Message string) *BaseHttpResponse {
	validationErrs := validation.GetValidationErrors(err)
	if validationErrs != nil {
		errorMessages := make([]string, 0)
		for _, validationErr := range *validationErrs {
			errorMessage := strings.ToLower(validationErr.Property)
			errorMessages = append(errorMessages, errorMessage)
		}
		titleMessage := "Error"
		alertMessage := "Please fill " + strings.Join(errorMessages, " , ") + " field correctly"
		return &BaseHttpResponse{
			Alert:            GenerateAlert(titleMessage, alertMessage),
			Status:           Status,
			ValidationErrors: validationErrs,
		}
	}

	return &BaseHttpResponse{
		Alert:  GenerateAlert(Title, Message),
		Status: Status,
	}
}
