package verification

import (
	"net/http"
	verificationController "verni/internal/controllers/verification"
	httpserver "verni/internal/http-server"
	"verni/internal/http-server/middleware"
	"verni/internal/http-server/responses"

	"github.com/gin-gonic/gin"
)

type VerificationController verificationController.Controller

func RegisterRoutes(
	router *gin.Engine,
	tokenChecker middleware.AccessTokenChecker,
	verification VerificationController,
) {
	methodGroup := router.Group("/verification", tokenChecker.Handler)
	methodGroup.PUT("/confirmEmail", func(c *gin.Context) {
		type ConfirmEmailRequest struct {
			Code string `json:"code"`
		}
		var request ConfirmEmailRequest
		if err := c.BindJSON(&request); err != nil {
			httpserver.AnswerWithBadRequest(c, err)
			return
		}
		if err := verification.ConfirmEmail(verificationController.UserId(tokenChecker.AccessToken(c)), request.Code); err != nil {
			switch err.Code {
			case verificationController.ConfirmEmailErrorWrongConfirmationCode:
				httpserver.Answer(c, err, http.StatusConflict, responses.CodeIncorrectCredentials)
			default:
				httpserver.AnswerWithUnknownError(c, err)
			}
		}
		c.JSON(http.StatusOK, responses.OK())
	})
	methodGroup.PUT("/sendEmailConfirmationCode", tokenChecker.Handler, func(c *gin.Context) {
		if err := verification.SendConfirmationCode(verificationController.UserId(tokenChecker.AccessToken(c))); err != nil {
			switch err.Code {
			default:
				httpserver.AnswerWithUnknownError(c, err)
			}
		}
		c.JSON(http.StatusOK, responses.OK())
	})
}