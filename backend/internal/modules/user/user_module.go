package user

import (
	"github.com/D3rise/dchat/internal/modules/user/repositories"
	"github.com/D3rise/dchat/internal/modules/user/services"
	"go.uber.org/fx"
)

var Module = fx.Module("user",
	// Repositories
	fx.Provide(
		repositories.NewUserRepository,
	),

	// Services
	fx.Provide(
		services.NewUserService,
	),
)
