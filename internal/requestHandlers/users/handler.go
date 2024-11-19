package users

import (
	"github.com/rzmn/Verni.App.Backend/internal/schema"
)

type RequestsHandler interface {
	GetUsers(
		subject schema.UserId,
		request schema.GetUsersRequest,
		success func(schema.StatusCode, schema.Response[[]schema.User]),
		failure func(schema.StatusCode, schema.Response[schema.Error]),
	)
}
