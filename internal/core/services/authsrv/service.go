package authsrv

import (
	"context"

	"github.com/tiemfah/lenswear-service/internal/core/domain"
	"github.com/tiemfah/lenswear-service/internal/core/ports"
)

type Service struct {
	authRepository ports.AuthenticationRepository
	userRepository ports.UserRepository
}

func NewAuthenticationService(authRepository ports.AuthenticationRepository, userRepository ports.UserRepository) *Service {
	return &Service{
		authRepository: authRepository,
		userRepository: userRepository,
	}
}

func (s *Service) Login(ctx context.Context, in *domain.LoginReq) (*domain.Tokens, error) {
	user, err := s.userRepository.CheckUserIsValidPassword(ctx, &domain.CheckUserIsValidPasswordReq{
		Username: in.Username,
		Password: in.Password,
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrInvalidCredentials
	}
	return s.authRepository.GrantTokens(ctx, user)
}
