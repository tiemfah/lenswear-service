package authrepo

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tiemfah/lenswear-service/internal/core/domain"
	ttjwt "github.com/tiemfah/lenswear-service/pkg/token/jwt"
)

const (
	issuer                         = "transformer-tracker-authentication-service"
	accessToken                    = "accessToken"
	refreshToken                   = "refreshToken"
	accessTokenExpireTimeInMinute  = 60
	refreshTokenExpireTimeInMinute = 60
)

type Repository struct {
	tokenManager ttjwt.JWTToken
}

func NewAuthenticationRepository(tokenManager ttjwt.JWTToken) *Repository {
	return &Repository{
		tokenManager: tokenManager,
	}
}

func (r *Repository) GrantTokens(ctx context.Context, user *domain.User) (*domain.Tokens, error) {
	accessToken, err := r.tokenManager.SignToken(&ttjwt.TokenClaim{
		UserID:   user.UserID,
		UserRole: user.UserRoleID,
		StandardClaims: jwt.StandardClaims{
			Id:        user.UserID,
			Issuer:    issuer,
			Subject:   accessToken,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * accessTokenExpireTimeInMinute).Unix(),
		},
	})
	if err != nil {
		return nil, err
	}
	return &domain.Tokens{
		AccessToken:  *accessToken,
		RefreshToken: "",
	}, nil
}
