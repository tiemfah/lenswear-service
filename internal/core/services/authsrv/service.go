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
	isValid, err := s.userRepository.CheckUserIsValidPassword(ctx, &domain.CheckUserIsValidPasswordReq{
		Username: in.Username,
		Password: in.Password,
	})
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, domain.ErrInvalidCredentials
	}
	user, err := s.userRepository.GetUserByUsername(ctx, in.Username)
	if err != nil {
		return nil, err
	}
	return s.authRepository.GrantTokens(ctx, user)
}
