package apparelsrv

import (
	"context"

	"github.com/tiemfah/lenswear-service/internal/core/domain"
	"github.com/tiemfah/lenswear-service/internal/core/ports"
	"github.com/tiemfah/lenswear-service/pkg/uidgen"
)

type Service struct {
	apparelRepository ports.ApparelRepository
	uidGen            uidgen.UIDGen
}

func NewApparelService(apparelRepository ports.ApparelRepository, uidGen uidgen.UIDGen) *Service {
	return &Service{
		apparelRepository: apparelRepository,
		uidGen:            uidGen,
	}
}

func (s *Service) CreateApparel(ctx context.Context, in *domain.CreateApparelReq) (*domain.EmptyRes, error) {
	if in.Requester.UserRole != domain.RoleAdmin {
		return nil, domain.ErrInsufficientPermissions
	}
	apparelID := s.uidGen.New(domain.ApparelPrefix + in.ApparelTypeID + "_")
	return s.apparelRepository.CreateApparel(ctx, apparelID, in)
}

func (s *Service) GetApparels(ctx context.Context, in *domain.GetApparelsReq) (*domain.Apparels, error) {
	return s.apparelRepository.GetApparels(ctx, in)
}

func (s *Service) GetApparelByApparelID(ctx context.Context, in *domain.GetApparelByApparelIDReq) (*domain.Apparel, error) {
	return s.apparelRepository.GetApparelByApparelID(ctx, in)
}

func (s *Service) DeleteApparelByApparelID(ctx context.Context, in *domain.DeleteApparelByApparelIDReq) (*domain.EmptyRes, error) {
	if in.Requester.UserRole != domain.RoleAdmin {
		return nil, domain.ErrInsufficientPermissions
	}
	return s.apparelRepository.DeleteApparelByApparelID(ctx, in)
}
