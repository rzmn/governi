package setDisplayName

import (
	"accounty/internal/http-server/responses"
)

type Request struct {
	DisplayName string `json:"displayName"`
}

type Error struct {
	responses.Error
}

func Success() responses.VoidResponse {
	return responses.OK()
}

func Failure(err Error) responses.Response[responses.Error] {
	return responses.Failure(err.Error)
}

func ErrInternal() Error {
	return Error{responses.Error{Code: responses.CodeInternal}}
}

func ErrBadRequest(description string) Error {
	return Error{responses.Error{Code: responses.CodeBadRequest, Description: &description}}
}

func ErrWrongFormat() Error {
	return Error{responses.Error{Code: responses.CodeWrongFormat}}
}