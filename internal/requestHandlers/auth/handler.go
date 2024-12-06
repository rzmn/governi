package auth

import (
	"github.com/rzmn/governi/internal/schema"
)

type RequestsHandler interface {
	Signup(
		request schema.SignupRequest,
		success func(schema.StatusCode, schema.Response[schema.Session]),
		failure func(schema.StatusCode, schema.Response[schema.Error]),
	)
	Login(
		request schema.LoginRequest,
		success func(schema.StatusCode, schema.Response[schema.Session]),
		failure func(schema.StatusCode, schema.Response[schema.Error]),
	)
	Refresh(
		request schema.RefreshRequest,
		success func(schema.StatusCode, schema.Response[schema.Session]),
		failure func(schema.StatusCode, schema.Response[schema.Error]),
	)
	UpdateEmail(
		subject schema.UserId,
		request schema.UpdateEmailRequest,
		success func(schema.StatusCode, schema.Response[schema.Session]),
		failure func(schema.StatusCode, schema.Response[schema.Error]),
	)
	UpdatePassword(
		subject schema.UserId,
		request schema.UpdatePasswordRequest,
		success func(schema.StatusCode, schema.Response[schema.Session]),
		failure func(schema.StatusCode, schema.Response[schema.Error]),
	)
	RegisterForPushNotifications(
		subject schema.UserId,
		request schema.RegisterForPushNotificationsRequest,
		success func(schema.StatusCode, schema.VoidResponse),
		failure func(schema.StatusCode, schema.Response[schema.Error]),
	)
	Logout(
		subject schema.UserId,
		success func(schema.StatusCode, schema.VoidResponse),
		failure func(schema.StatusCode, schema.Response[schema.Error]),
	)
}
