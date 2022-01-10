package userrepo

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/tiemfah/lenswear-service/internal/core/domain"
	"github.com/tiemfah/lenswear-service/internal/rdbms/postgresql/gen/lenswear/public/model"
	"github.com/tiemfah/lenswear-service/internal/rdbms/postgresql/gen/lenswear/public/table"
	"github.com/tiemfah/lenswear-service/pkg/hash"
)

type Repository struct {
	DB   *sql.DB
	hash hash.Hash
}

func NewUserRepository(db *sql.DB, hash hash.Hash) *Repository {
	createTable(db)
	return &Repository{
		DB:   db,
		hash: hash,
	}
}

func createTable(db *sql.DB) error {
	var err error

	const roleQry = `CREATE TABLE IF NOT EXISTS role (
		role_id  varchar(200) PRIMARY KEY NOT NULL
	)`

	if _, err = db.Exec(roleQry); err != nil {
		log.Fatal("cannot create role table: ", err)
		return err
	}

	const areaQry = `CREATE TABLE IF NOT EXISTS area (
		area_id  varchar(200) PRIMARY KEY NOT NULL
	)`

	if _, err = db.Exec(areaQry); err != nil {
		log.Fatal("cannot create area table: ", err)
		return err
	}

	const userQry = `CREATE TABLE IF NOT EXISTS users (
		user_id     varchar(200) PRIMARY KEY NOT NULL,
		role_id     varchar(200) NOT NULL,
		username    varchar(200) NOT NULL UNIQUE,
		password    varchar(200),
		create_date timestamp NOT NULL DEFAULT NOW(),
		update_date timestamp,
		delete_date timestamp,
		create_by   varchar(200),
		update_by   varchar(200),
		delete_by   varchar(200),
		CONSTRAINT fk_role
			FOREIGN KEY(role_id)
				REFERENCES role(role_id)
	)`

	if _, err = db.Exec(userQry); err != nil {
		log.Fatal("cannot create users table: ", err)
		return err
	}

	const userRoleQry = `
	INSERT INTO role (role_id)
	VALUES
	('user'),
	('admin')
	ON CONFLICT DO NOTHING;
	`
	if _, err = db.Exec(userRoleQry); err != nil {
		log.Fatal("cannot insert into role table: ", err)
		return err
	}

	return nil
}

func (r *Repository) CreateUser(ctx context.Context, userID string, in *domain.CreateUserReq) (*domain.EmptyRes, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	stmt := table.Users.
		INSERT(
			table.Users.UserID,
			table.Users.RoleID,
			table.Users.Username,
			table.Users.Password,
		).MODEL(model.Users{
		UserID:   userID,
		RoleID:   domain.RoleUser,
		Username: in.Username,
		Password: &in.Password,
	})
	if _, err := stmt.Exec(tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &domain.EmptyRes{}, nil
}

func (r *Repository) GetUsersAsAdmin(ctx context.Context, in *domain.GetUsersAsAdminReq) (*domain.Users, error) {
	res := &domain.Users{Users: []*domain.User{}}
	stmt := postgres.SELECT(
		table.Users.AllColumns,
	).FROM(
		table.Users,
	).ORDER_BY(
		table.Users.CreateDate.ASC(),
	).LIMIT(
		in.Limit,
	).OFFSET(
		in.Offset,
	)
	var users []model.Users
	if err := stmt.Query(r.DB, &users); err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return res, nil
		}
		return nil, err
	}
	for _, u := range users {
		res.Users = append(res.Users, &domain.User{
			UserID:     u.UserID,
			UserRoleID: u.RoleID,
			Username:   u.Username,
		})
	}
	return res, nil
}

func (r *Repository) GetUserByUserID(ctx context.Context, in *domain.GetUserByUserIDReq) (*domain.User, error) {
	stmt := postgres.SELECT(
		table.Users.AllColumns,
	).FROM(
		table.Users,
	).WHERE(
		table.Users.UserID.EQ(postgres.String(in.UserID)),
	)
	var user model.Users
	if err := stmt.Query(r.DB, &user); err != nil {
		return nil, err
	}
	return &domain.User{
		UserID:     user.UserID,
		UserRoleID: user.RoleID,
		Username:   user.Username,
	}, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	stmt := postgres.SELECT(
		table.Users.AllColumns,
	).FROM(
		table.Users,
	).WHERE(
		table.Users.Username.EQ(postgres.String(username)),
	)
	var user model.Users
	if err := stmt.Query(r.DB, &user); err != nil {
		return nil, err
	}
	return &domain.User{
		UserID:     user.UserID,
		UserRoleID: user.RoleID,
		Username:   user.Username,
	}, nil
}

func (r *Repository) ModifyUser(ctx context.Context, in *domain.ModifyUserReq) (*domain.EmptyRes, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	stmt := table.Users.UPDATE().
		SET(
			table.Users.UpdateDate.SET(postgres.TimestampT(time.Now())),
			table.Users.UpdateBy.SET(postgres.String(in.Requester.UserID)),
		).
		WHERE(table.Users.Username.EQ(postgres.String(in.Username)))
	if _, err = stmt.Exec(tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &domain.EmptyRes{}, nil
}

func (r *Repository) ResetUserPassword(ctx context.Context, in *domain.ResetUserPasswordReq) (*domain.EmptyRes, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	stmt := table.Users.UPDATE().
		SET(
			table.Users.Username.SET(postgres.String(in.Username)),
			table.Users.Password.SET(postgres.String(in.Password)),
			table.Users.UpdateDate.SET(postgres.TimestampT(time.Now())),
			table.Users.UpdateBy.SET(postgres.String(in.Requester.UserID)),
		).
		WHERE(table.Users.Username.EQ(postgres.String(in.Username)))
	if _, err = stmt.Exec(tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &domain.EmptyRes{}, nil
}

func (r *Repository) CheckUserIsValidPassword(ctx context.Context, in *domain.CheckUserIsValidPasswordReq) (bool, error) {
	stmt := postgres.SELECT(
		table.Users.AllColumns,
	).FROM(
		table.Users,
	).WHERE(
		table.Users.Username.EQ(postgres.String(in.Username)),
	)
	var user model.Users
	if err := stmt.Query(r.DB, &user); err != nil {
		return false, err
	}
	isValid, err := r.hash.ComparePassword(in.Password, *user.Password)
	if err != nil {
		return false, err
	}
	return isValid, nil
}
