package defaultUsersHandler

import (
	"net/http"

	"github.com/rzmn/governi/internal/common"
	usersController "github.com/rzmn/governi/internal/controllers/users"
	"github.com/rzmn/governi/internal/requestHandlers/users"
	"github.com/rzmn/governi/internal/schema"
	"github.com/rzmn/governi/internal/services/logging"
)

func New(
	controller usersController.Controller,
	logger logging.Service,
) users.RequestsHandler {
	return &defaultRequestsHandler{
		controller: controller,
		logger:     logger,
	}
}

type defaultRequestsHandler struct {
	controller usersController.Controller
	logger     logging.Service
}

func (c *defaultRequestsHandler) GetUsers(
	subject schema.UserId,
	request schema.GetUsersRequest,
	success func(schema.StatusCode, schema.Response[[]schema.User]),
	failure func(schema.StatusCode, schema.Response[schema.Error]),
) {
	users, err := c.controller.Get(common.Map(request.Ids, func(id schema.UserId) usersController.UserId {
		return usersController.UserId(id)
	}), usersController.UserId(subject))
	if err != nil {
		switch err.Code {
		default:
			c.logger.LogError("getUsers request %v failed with unknown err: %v", request, err)
			failure(http.StatusInternalServerError, schema.Failure(err, schema.CodeInternal))
		}
		return
	}
	success(http.StatusOK, schema.Success(common.Map(users, mapUser)))
}

func (c *defaultRequestsHandler) SearchUsers(
	subject schema.UserId,
	request schema.SearchUsersRequest,
	success func(schema.StatusCode, schema.Response[[]schema.User]),
	failure func(schema.StatusCode, schema.Response[schema.Error]),
) {
	users, err := c.controller.Search(request.Query, usersController.UserId(subject))
	if err != nil {
		switch err.Code {
		default:
			c.logger.LogError("searchUsers request %v failed with unknown err: %v", request, err)
			failure(http.StatusInternalServerError, schema.Failure(err, schema.CodeInternal))
		}
		return
	}
	success(http.StatusOK, schema.Success(common.Map(users, mapUser)))
}

func mapUser(user usersController.User) schema.User {
	return schema.User{
		Id:           schema.UserId(user.Id),
		DisplayName:  user.DisplayName,
		AvatarId:     (*schema.ImageId)(user.AvatarId),
		FriendStatus: schema.FriendStatus(user.FriendStatus),
	}
}
