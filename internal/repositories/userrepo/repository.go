package userrepo

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/tiemfah/lenswear-service/internal/core/domain"
	"github.com/tiemfah/lenswear-service/pkg/hash"
	"google.golang.org/api/iterator"
)

type Repository struct {
	DS   *datastore.Client
	hash hash.Hash
}

func NewUserRepository(ds *datastore.Client, hash hash.Hash) *Repository {
	return &Repository{
		DS:   ds,
		hash: hash,
	}
}

func (r *Repository) CreateUser(ctx context.Context, userID string, in *domain.CreateUserReq) (*domain.EmptyRes, error) {
	tx, err := r.DS.NewTransaction(ctx)
	if err != nil {
		return nil, err
	}
	userKey := datastore.NameKey(domain.UserDataStoreKey, userID, nil)
	user := &domain.UserWithPassword{
		UserID:     userID,
		UserRoleID: domain.RoleUser,
		Username:   in.Username,
		Password:   in.Password,
	}
	if err := tx.Get(userKey, user); err != datastore.ErrNoSuchEntity {
		tx.Rollback()
		return nil, err
	}
	if _, err := tx.Put(userKey, user); err != nil {
		tx.Rollback()
		return nil, err
	}
	if _, err = tx.Commit(); err != nil {
		return nil, err
	}
	return &domain.EmptyRes{}, nil
}

func (r *Repository) GetUsersAsAdmin(ctx context.Context, in *domain.GetUsersAsAdminReq) (*domain.Users, error) {
	return &domain.Users{Users: []*domain.User{}}, nil
}

func (r *Repository) GetUserByUserID(ctx context.Context, in *domain.GetUserByUserIDReq) (*domain.User, error) {
	return &domain.User{}, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return &domain.User{}, nil
}

func (r *Repository) ModifyUser(ctx context.Context, in *domain.ModifyUserReq) (*domain.EmptyRes, error) {
	return &domain.EmptyRes{}, nil
}

func (r *Repository) ResetUserPassword(ctx context.Context, in *domain.ResetUserPasswordReq) (*domain.EmptyRes, error) {
	return &domain.EmptyRes{}, nil
}

func (r *Repository) CheckUserIsValidPassword(ctx context.Context, in *domain.CheckUserIsValidPasswordReq) (bool, error) {
	query := datastore.NewQuery(domain.UserDataStoreKey).Filter("Username =", in.Username).Limit(1)
	it := r.DS.Run(ctx, query)
	var user domain.UserWithPassword
	for {
		_, err := it.Next(&user)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, err
		}
	}
	isValid, err := r.hash.ComparePassword(in.Password, user.Password)
	if err != nil {
		return false, err
	}
	return isValid, nil
}
