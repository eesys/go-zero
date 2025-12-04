package model

import (
	"context"
	"database/sql"

	"github.com/eesys/go-zero/core/stores/gormx"
)

type UserRepo struct {
	r  gormx.Repository[User]
	rw *gormx.Router
}

func NewUserRepo(r gormx.Repository[User], rw *gormx.Router) UserRepo {
	return UserRepo{r: r, rw: rw}
}

type simpleResult struct{}

func (simpleResult) LastInsertId() (int64, error) { return 0, nil }
func (simpleResult) RowsAffected() (int64, error) { return 0, nil }

func (u UserRepo) Insert(data User) (sql.Result, error) {
	ctx := gormx.WithWrite(context.Background())
	if err := u.r.Create(ctx, &data); err != nil {
		return nil, err
	}
	return simpleResult{}, nil
}

func (u UserRepo) FindOne(id int64) (*User, error) {
	ctx := gormx.WithReadReplica(context.Background())
	out, err := u.r.First(ctx, User{ID: id})
	if err == gormx.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (u UserRepo) FindOneByUser(user string) (*User, error) {
	ctx := gormx.WithReadReplica(context.Background())
	out, err := u.r.First(ctx, User{User: user})
	if err == gormx.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (u UserRepo) FindOneByMobile(mobile string) (*User, error) {
	ctx := gormx.WithReadReplica(context.Background())
	out, err := u.r.First(ctx, User{Mobile: mobile})
	if err == gormx.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (u UserRepo) FindOneByName(name string) (*User, error) {
	ctx := gormx.WithReadReplica(context.Background())
	out, err := u.r.First(ctx, User{Name: name})
	if err == gormx.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (u UserRepo) Update(data User) error {
	ctx := gormx.WithWrite(context.Background())
	return u.r.Save(ctx, &data)
}

func (u UserRepo) Delete(id int64) error {
	ctx := gormx.WithWrite(context.Background())
	return u.r.Delete(ctx, User{ID: id})
}
