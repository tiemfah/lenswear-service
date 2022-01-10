package apparelrepo

import (
	"context"
	"database/sql"
	"log"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/tiemfah/lenswear-service/internal/core/domain"
	"github.com/tiemfah/lenswear-service/internal/rdbms/postgresql/gen/lenswear/public/model"
	"github.com/tiemfah/lenswear-service/internal/rdbms/postgresql/gen/lenswear/public/table"
	"github.com/tiemfah/lenswear-service/pkg/bucket"
)

type Repository struct {
	DB        *sql.DB
	GCPBucket *bucket.GCPBucket
}

func createTable(db *sql.DB) error {
	var err error

	const apparelTypeQry = `CREATE TABLE IF NOT EXISTS apparel_type (
		apparel_type_id  varchar(200) PRIMARY KEY NOT NULL
	)`

	if _, err = db.Exec(apparelTypeQry); err != nil {
		log.Fatal("cannot create apparel_type table: ", err)
		return err
	}

	const apparelTypeMasterDataQry = `
	INSERT INTO apparel_type (apparel_type_id)
	VALUES
	('shirt'),
	('tshirt'),
	('pants'),
	('shorts'),
	('skirt'),
	('dress')
	ON CONFLICT DO NOTHING;
	`
	if _, err = db.Exec(apparelTypeMasterDataQry); err != nil {
		log.Fatal("cannot insert into apparel_type table: ", err)
		return err
	}

	const apparelQry = `CREATE TABLE IF NOT EXISTS apparel (
		apparel_id						varchar(200) PRIMARY KEY NOT NULL,
		apparel_type_id					varchar(200) NOT NULL,
		name     						varchar(200) NOT NULL,
		brand							varchar(200) NOT NULL,
		price							varchar(200) NOT NULL,
		store_url						varchar(200) NOT NULL,
		create_date 					timestamp NOT NULL DEFAULT NOW(),
		update_date 					timestamp,
		create_by   					varchar(200),
		update_by   					varchar(200),
		CONSTRAINT fk_apparel_type
			FOREIGN KEY(apparel_type_id)
				REFERENCES apparel_type(apparel_type_id)
	)`

	if _, err = db.Exec(apparelQry); err != nil {
		log.Fatal("cannot create apparel table: ", err)
		return err
	}

	const apparelImgQry = `CREATE TABLE IF NOT EXISTS apparel_image (
		apparel_id						varchar(200) NOT NULL,
		apparel_image_url				varchar(200) NOT NULL,
		create_date 					timestamp NOT NULL DEFAULT NOW(),
		update_date 					timestamp,
		create_by   					varchar(200),
		update_by   					varchar(200),
		CONSTRAINT fk_apparel
			FOREIGN KEY(apparel_id)
				REFERENCES apparel(apparel_id)
	)`

	if _, err = db.Exec(apparelImgQry); err != nil {
		log.Fatal("cannot create apparel_image table: ", err)
		return err
	}

	return nil
}

func NewApprelRepository(db *sql.DB, gcpBucket *bucket.GCPBucket) *Repository {
	createTable(db)
	return &Repository{
		DB:        db,
		GCPBucket: gcpBucket,
	}
}

func (r *Repository) CreateApparel(ctx context.Context, ApparelID string, in *domain.CreateApparelReq) (*domain.EmptyRes, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	filePaths, err := r.GCPBucket.UploadFile(in.ApparelTypeID, ApparelID, in.Files)
	if err != nil {
		return nil, err
	}
	tcStmt := table.Apparel.
		INSERT(
			table.Apparel.AllColumns,
		).MODEL(model.Apparel{
		ApparelID:     ApparelID,
		ApparelTypeID: in.ApparelTypeID,
		Name:          in.Name,
		Brand:         in.Brand,
		Price:         in.Price,
		StoreURL:      in.StoreURL,
		CreateBy:      &in.Requester.UserID,
	})
	if _, err := tcStmt.Exec(tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, p := range filePaths {
		tStmt := table.ApparelImage.
			INSERT(
				table.ApparelImage.AllColumns,
			).MODEL(model.ApparelImage{
			ApparelID:       ApparelID,
			ApparelImageURL: p,
			CreateBy:        &in.Requester.UserID,
		})
		if _, err := tStmt.Exec(tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &domain.EmptyRes{}, nil
}

func (r *Repository) GetApparels(ctx context.Context, in *domain.GetApparelsReq) (*domain.Apparels, error) {
	res := &domain.Apparels{Apparels: []*domain.Apparel{}}
	stmt := postgres.SELECT(
		table.Apparel.AllColumns,
		table.ApparelImage.AllColumns,
	).FROM(
		table.Apparel.LEFT_JOIN(table.ApparelImage, table.ApparelImage.ApparelID.EQ(table.Apparel.ApparelID)),
	).ORDER_BY(
		table.Apparel.CreateDate.ASC(),
	).LIMIT(
		in.Limit,
	).OFFSET(
		in.Offset,
	)
	var apparels []struct {
		model.Apparel
		ApparelImages []model.ApparelImage
	}
	if err := stmt.Query(r.DB, &apparels); err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return res, nil
		}
		return nil, err
	}
	for _, a := range apparels {
		imgURLs := []string{}
		for _, ai := range a.ApparelImages {
			imgURLs = append(imgURLs, ai.ApparelImageURL)
		}
		res.Apparels = append(res.Apparels, &domain.Apparel{
			ApparelID: a.ApparelID,
			Name:      a.Name,
			Brand:     a.Brand,
			Price:     a.Price,
			StoreURL:  a.StoreURL,
			ImgURLs:   imgURLs,
		})
	}
	return res, nil
}

func (r *Repository) GetApparelByApparelID(ctx context.Context, in *domain.GetApparelByApparelIDReq) (*domain.Apparel, error) {
	stmt := postgres.SELECT(
		table.Apparel.AllColumns,
		table.ApparelImage.AllColumns,
	).FROM(
		table.Apparel.LEFT_JOIN(table.ApparelImage, table.ApparelImage.ApparelID.EQ(table.Apparel.ApparelID)),
	).WHERE(
		table.Apparel.ApparelID.EQ(postgres.String(in.ApparelID)),
	)
	var apparel struct {
		model.Apparel
		ApparelImages []model.ApparelImage
	}
	if err := stmt.Query(r.DB, &apparel); err != nil {
		return nil, err
	}
	imgURLs := []string{}
	for _, ai := range apparel.ApparelImages {
		imgURLs = append(imgURLs, ai.ApparelImageURL)
	}
	return &domain.Apparel{
		ApparelID: apparel.ApparelID,
		Name:      apparel.Name,
		Brand:     apparel.Brand,
		Price:     apparel.Price,
		StoreURL:  apparel.StoreURL,
		ImgURLs:   imgURLs,
	}, nil
}

func (r *Repository) DeleteApparelByApparelID(ctx context.Context, in *domain.DeleteApparelByApparelIDReq) (*domain.EmptyRes, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	stmt := table.Apparel.
		DELETE().
		WHERE(table.Apparel.ApparelID.EQ(postgres.String(in.ApparelID)))
	if _, err := stmt.Exec(tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &domain.EmptyRes{}, nil
}
