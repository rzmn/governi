package avatars

import (
	"github.com/rzmn/governi/internal/schema"
)

type RequestsHandler interface {
	GetAvatars(
		request schema.GetAvatarsRequest,
		success func(schema.StatusCode, schema.Response[map[schema.ImageId]schema.Image]),
		failure func(schema.StatusCode, schema.Response[schema.Error]),
	)
}
