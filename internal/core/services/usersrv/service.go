package usersrv

import (
	"context"

	"github.com/tiemfah/lenswear-service/internal/core/domain"
	"github.com/tiemfah/lenswear-service/internal/core/ports"
	"github.com/tiemfah/lenswear-service/pkg/hash"
	"github.com/tiemfah/lenswear-service/pkg/uidgen"
)

type Service struct {
	userRepository ports.UserRepository
	uidGen         uidgen.UIDGen
	hash           hash.Hash
}

func NewUserService(userRepository ports.UserRepository, uidGen uidgen.UIDGen, hash hash.Hash) *Service {
	return &Service{
		userRepository: userRepository,
		uidGen:         uidGen,
		hash:           hash,
	}
}

func (s *Service) CreateUser(ctx context.Context, in *domain.CreateUserReq) (*domain.EmptyRes, error) {
	if in.Requester.UserRole != domain.RoleAdmin {
		return nil, domain.ErrInsufficientPermissions
	}
	hashedPassword, err := s.hash.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	in.Password = hashedPassword
	return s.userRepository.CreateUser(ctx, s.uidGen.New(domain.UserPrefix), in)
}

func (s *Service) GetUsersAsAdmin(ctx context.Context, in *domain.GetUsersAsAdminReq) (*domain.Users, error) {
	if in.Requester.UserRole != domain.RoleAdmin {
		return nil, domain.ErrInsufficientPermissions
	}
	return s.userRepository.GetUsersAsAdmin(ctx, in)
}

func (s *Service) GetUserByUserID(ctx context.Context, in *domain.GetUserByUserIDReq) (*domain.User, error) {
	return s.userRepository.GetUserByUserID(ctx, in)
}

func (s *Service) ModifyUser(ctx context.Context, in *domain.ModifyUserReq) (*domain.EmptyRes, error) {
	if in.Requester.UserRole != domain.RoleAdmin {
		return nil, domain.ErrInsufficientPermissions
	}
	return s.userRepository.ModifyUser(ctx, in)
}

func (s *Service) ResetUserPassword(ctx context.Context, in *domain.ResetUserPasswordReq) (*domain.EmptyRes, error) {
	if in.Requester.UserRole != domain.RoleAdmin {
		return nil, domain.ErrInsufficientPermissions
	}
	hashedPassword, err := s.hash.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	in.Password = hashedPassword
	return s.userRepository.ResetUserPassword(ctx, in)
}
