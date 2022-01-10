package ports

import (
	"context"

	"github.com/tiemfah/lenswear-service/internal/core/domain"
)

type AuthenticationService interface {
	Login(ctx context.Context, in *domain.LoginReq) (*domain.Tokens, error)
}

type UserService interface {
	CreateUser(ctx context.Context, in *domain.CreateUserReq) (*domain.EmptyRes, error)
	GetUsersAsAdmin(ctx context.Context, in *domain.GetUsersAsAdminReq) (*domain.Users, error)
	GetUserByUserID(ctx context.Context, in *domain.GetUserByUserIDReq) (*domain.User, error)
	ModifyUser(ctx context.Context, in *domain.ModifyUserReq) (*domain.EmptyRes, error)
	ResetUserPassword(ctx context.Context, in *domain.ResetUserPasswordReq) (*domain.EmptyRes, error)
}

type ApparelService interface {
	CreateApparel(ctx context.Context, in *domain.CreateApparelReq) (*domain.EmptyRes, error)
	GetApparels(ctx context.Context, in *domain.GetApparelsReq) (*domain.Apparels, error)
	GetApparelByApparelID(ctx context.Context, in *domain.GetApparelByApparelIDReq) (*domain.Apparel, error)
	DeleteApparelByApparelID(ctx context.Context, in *domain.DeleteApparelByApparelIDReq) (*domain.EmptyRes, error)
}
