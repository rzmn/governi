package auth

import (
	"accounty/internal/auth/confirmation"
	"accounty/internal/auth/jwt"
	httpserver "accounty/internal/http-server"
	"accounty/internal/http-server/middleware"
	"accounty/internal/http-server/responses"
	authController "accounty/internal/http-server/router/auth"
	"accounty/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db storage.Storage, jwtService jwt.Service) {
	ensureLoggedIn := middleware.EnsureLoggedIn(db, jwtService)
	hostFromToken := func(c *gin.Context) storage.UserId {
		return storage.UserId(c.Request.Header.Get(middleware.LoggedInSubjectKey))
	}
	controller := authController.DefaultController(db, jwtService, confirmation.EmailConfirmation{})
	router.PUT("/auth/signup", func(c *gin.Context) {
		type SignupRequest struct {
			Credentials storage.UserCredentials `json:"credentials"`
		}
		var request SignupRequest
		if err := c.BindJSON(&request); err != nil {
			httpserver.AnswerWithBadRequest(c, err)
			return
		}
		session, err := controller.Signup(request.Credentials.Email, request.Credentials.Password)
		if err != nil {
			switch err.Code {
			case authController.SignupErrorAlreadyTaken:
				c.JSON(
					http.StatusConflict,
					responses.Failure(
						responses.Error{
							Code: responses.CodeAlreadyTaken,
						},
					),
				)
			case authController.SignupErrorWrongFormat:
				c.JSON(
					http.StatusUnprocessableEntity,
					responses.Failure(
						responses.Error{
							Code: responses.CodeWrongFormat,
						},
					),
				)
			default:
				httpserver.AnswerWithUnknownError(c, err)
			}
			return
		}
		c.JSON(http.StatusOK, responses.Success(session))
	})
	router.PUT("/auth/login", func(c *gin.Context) {
		type LoginRequest struct {
			Credentials storage.UserCredentials `json:"credentials"`
		}
		var request LoginRequest
		if err := c.BindJSON(&request); err != nil {
			httpserver.AnswerWithBadRequest(c, err)
			return
		}
		session, err := controller.Login(request.Credentials.Email, request.Credentials.Password)
		if err != nil {
			switch err.Code {
			default:
				httpserver.AnswerWithUnknownError(c, err)
			}
		}
		c.JSON(http.StatusOK, responses.Success(session))
	})
	router.PUT("/auth/refresh", func(c *gin.Context) {
		type RefreshRequest struct {
			RefreshToken string `json:"refreshToken"`
		}
		var request RefreshRequest
		if err := c.BindJSON(&request); err != nil {
			httpserver.AnswerWithBadRequest(c, err)
			return
		}
		session, err := controller.Refresh(request.RefreshToken)
		if err != nil {
			switch err.Code {
			default:
				httpserver.AnswerWithUnknownError(c, err)
			}
		}
		c.JSON(http.StatusOK, responses.Success(session))
	})
	router.PUT("/auth/updateEmail", ensureLoggedIn, func(c *gin.Context) {
		type UpdateEmailRequest struct {
			Email string `json:"email"`
		}
		var request UpdateEmailRequest
		if err := c.BindJSON(&request); err != nil {
			httpserver.AnswerWithBadRequest(c, err)
			return
		}
		session, err := controller.UpdateEmail(request.Email, authController.UserId(hostFromToken(c)))
		if err != nil {
			switch err.Code {
			default:
				httpserver.AnswerWithUnknownError(c, err)
			}
		}
		c.JSON(http.StatusOK, responses.Success(session))
	})
	router.PUT("/auth/updatePassword", ensureLoggedIn, func(c *gin.Context) {
		type UpdatePasswordRequest struct {
			OldPassword string `json:"old"`
			NewPassword string `json:"new"`
		}
		var request UpdatePasswordRequest
		if err := c.BindJSON(&request); err != nil {
			httpserver.AnswerWithBadRequest(c, err)
			return
		}
		session, err := controller.UpdatePassword(request.OldPassword, request.NewPassword, authController.UserId(hostFromToken(c)))
		if err != nil {
			switch err.Code {
			default:
				httpserver.AnswerWithUnknownError(c, err)
			}
		}
		c.JSON(http.StatusOK, responses.Success(session))
	})
	router.DELETE("/auth/logout", ensureLoggedIn, func(c *gin.Context) {
		if err := controller.Logout(authController.UserId(hostFromToken(c))); err != nil {
			switch err.Code {
			default:
				httpserver.AnswerWithUnknownError(c, err)
			}
		}
		c.JSON(http.StatusOK, responses.OK())
	})
	router.PUT("/auth/confirmEmail", ensureLoggedIn, func(c *gin.Context) {
		type ConfirmEmailRequest struct {
			Code string `json:"code"`
		}
		var request ConfirmEmailRequest
		if err := c.BindJSON(&request); err != nil {
			httpserver.AnswerWithBadRequest(c, err)
			return
		}
		if err := controller.ConfirmEmail(request.Code, authController.UserId(hostFromToken(c))); err != nil {
			switch err.Code {
			default:
				httpserver.AnswerWithUnknownError(c, err)
			}
		}
		c.JSON(http.StatusOK, responses.OK())
	})
	router.PUT("/auth/sendEmailConfirmationCode", ensureLoggedIn, func(c *gin.Context) {
		if err := controller.SendEmailConfirmationCode(authController.UserId(hostFromToken(c))); err != nil {
			switch err.Code {
			default:
				httpserver.AnswerWithUnknownError(c, err)
			}
		}
		c.JSON(http.StatusOK, responses.OK())
	})
	router.PUT("/auth/registerForPushNotifications", ensureLoggedIn, func(c *gin.Context) {
		type RegisterForPushNotificationsRequest struct {
			Token string `json:"token"`
		}
		var request RegisterForPushNotificationsRequest
		if err := c.BindJSON(&request); err != nil {
			httpserver.AnswerWithBadRequest(c, err)
			return
		}
		if err := controller.RegisterForPushNotifications(request.Token, authController.UserId(hostFromToken(c))); err != nil {
			switch err.Code {
			default:
				httpserver.AnswerWithUnknownError(c, err)
			}
		}
		c.JSON(http.StatusOK, responses.OK())
	})
}
