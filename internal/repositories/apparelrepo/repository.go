package apparelrepo

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/tiemfah/lenswear-service/internal/core/domain"
	"github.com/tiemfah/lenswear-service/pkg/bucket"
	"google.golang.org/api/iterator"
)

type Repository struct {
	DS        *datastore.Client
	GCPBucket *bucket.GCPBucket
}

func NewApprelRepository(ds *datastore.Client, gcpBucket *bucket.GCPBucket) *Repository {
	return &Repository{
		DS:        ds,
		GCPBucket: gcpBucket,
	}
}

func (r *Repository) CreateApparel(ctx context.Context, apparelID string, in *domain.CreateApparelReq) (*domain.EmptyRes, error) {
	tx, err := r.DS.NewTransaction(ctx)
	if err != nil {
		return nil, err
	}
	filePaths, err := r.GCPBucket.UploadFile(in.ApparelTypeID, apparelID, in.Files)
	if err != nil {
		return nil, err
	}
	apparelKey := datastore.NameKey(domain.ApprelDataStoreKey, apparelID, nil)
	apparel := &domain.Apparel{
		ApparelTypeID: in.ApparelTypeID,
		Name:          in.Name,
		Brand:         in.Brand,
		Price:         in.Price,
		StoreURL:      in.StoreURL,
		ImgURLs:       filePaths,
		CreateDate:    time.Now(),
		CreateBy:      in.Requester.UserID,
	}
	if err := tx.Get(apparelKey, apparel); err != datastore.ErrNoSuchEntity {
		tx.Rollback()
		return nil, err
	}
	if _, err := tx.Put(apparelKey, apparel); err != nil {
		tx.Rollback()
		return nil, err
	}
	if _, err = tx.Commit(); err != nil {
		return nil, err
	}
	return &domain.EmptyRes{}, nil
}

func (r *Repository) GetApparels(ctx context.Context, in *domain.GetApparelsReq) (*domain.Apparels, error) {
	res := &domain.Apparels{Apparels: []*domain.Apparel{}}
	var query *datastore.Query
	if in.ApparelTypeID != "" {
		query = datastore.NewQuery(domain.ApprelDataStoreKey).Filter("ApparelTypeID=", in.ApparelTypeID).Offset(int(in.Offset)).Limit(int(in.Limit))
	} else {
		query = datastore.NewQuery(domain.ApprelDataStoreKey).Offset(int(in.Offset)).Limit(int(in.Limit))
	}
	it := r.DS.Run(ctx, query)
	for {
		var apparel domain.Apparel
		key, err := it.Next(&apparel)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		apparel.ApparelID = key.Name
		res.Apparels = append(res.Apparels, &apparel)
	}
	return res, nil
}

func (r *Repository) GetApparelByApparelID(ctx context.Context, in *domain.GetApparelByApparelIDReq) (*domain.Apparel, error) {
	apparel := &domain.Apparel{}
	apparelKey := datastore.NameKey(domain.ApprelDataStoreKey, in.ApparelID, nil)
	if err := r.DS.Get(ctx, apparelKey, apparel); err != nil {
		return nil, err
	}
	return apparel, nil
}

func (r *Repository) DeleteApparelByApparelID(ctx context.Context, in *domain.DeleteApparelByApparelIDReq) (*domain.EmptyRes, error) {
	tx, err := r.DS.NewTransaction(ctx)
	if err != nil {
		return nil, err
	}
	apparelKey := datastore.NameKey(domain.ApprelDataStoreKey, in.ApparelID, nil)
	if err := tx.Delete(apparelKey); err != nil {
		return nil, err
	}
	if err := r.GCPBucket.DeleteFolder(in.ApparelTypeID, in.ApparelID); err != nil {
		return nil, err
	}
	if _, err := tx.Commit(); err != nil {
		return nil, err
	}
	return &domain.EmptyRes{}, nil
}
