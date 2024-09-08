package updateEmail

import (
	"accounty/internal/http-server/responses"
	"accounty/internal/storage"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestHandler interface {
	Validate(c *gin.Context, request Request) *Error
	Handle(c *gin.Context, request Request) (storage.AuthenticatedSession, *Error)
}

func handleError(c *gin.Context, err Error) {
	switch err.Code {
	case responses.CodeWrongFormat:
		c.JSON(http.StatusUnprocessableEntity, Failure(err))
	case responses.CodeBadRequest:
		c.JSON(http.StatusBadRequest, Failure(err))
	default:
		c.JSON(http.StatusInternalServerError, Failure(err))
	}
}

func New(requestHandler RequestHandler) func(c *gin.Context) {
	return func(c *gin.Context) {
		const op = "handlers.auth.updateEmail"
		var request Request
		if err := c.BindJSON(&request); err != nil {
			handleError(c, ErrBadRequest(fmt.Sprintf("%s: request failed %v", op, err)))
			return
		}
		if err := requestHandler.Validate(c, request); err != nil {
			handleError(c, *err)
			return
		}
		token, err := requestHandler.Handle(c, request)
		if err != nil {
			handleError(c, *err)
			return
		}
		c.JSON(http.StatusCreated, Success(token))
	}
}
