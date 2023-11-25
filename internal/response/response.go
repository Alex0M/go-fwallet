package response

import (
	"go-fwallet/internal/models"
)

func NotFound() models.Response {
	return models.Response{
		Status:  false,
		Message: "Not Found",
	}
}

func ErrorResponse(e string) models.Response {
	return models.Response{
		Status:  false,
		Message: "Error",
		Data: models.ErrorResponse{
			Error: e,
		},
	}
}

func SuccessResponse(m string, data any) models.Response {
	return models.Response{
		Status:  true,
		Message: m,
		Data:    data,
	}
}
