package ports

import (
	"context"

	"github.com/tiemfah/lenswear-service/internal/core/domain"
)

type AuthenticationRepository interface {
	GrantTokens(ctx context.Context, user *domain.User) (*domain.Tokens, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, userID string, in *domain.CreateUserReq) (*domain.EmptyRes, error)
	GetUsersAsAdmin(ctx context.Context, in *domain.GetUsersAsAdminReq) (*domain.Users, error)
	GetUserByUserID(ctx context.Context, in *domain.GetUserByUserIDReq) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	ModifyUser(ctx context.Context, in *domain.ModifyUserReq) (*domain.EmptyRes, error)
	ResetUserPassword(ctx context.Context, in *domain.ResetUserPasswordReq) (*domain.EmptyRes, error)
	// helpers
	CheckUserIsValidPassword(ctx context.Context, in *domain.CheckUserIsValidPasswordReq) (*domain.User, error)
}

type ApparelRepository interface {
	CreateApparel(ctx context.Context, ApparelID string, in *domain.CreateApparelReq) (*domain.EmptyRes, error)
	GetApparels(ctx context.Context, in *domain.GetApparelsReq) (*domain.Apparels, error)
	GetApparelByApparelID(ctx context.Context, in *domain.GetApparelByApparelIDReq) (*domain.Apparel, error)
	DeleteApparelByApparelID(ctx context.Context, in *domain.DeleteApparelByApparelIDReq) (*domain.EmptyRes, error)
}
