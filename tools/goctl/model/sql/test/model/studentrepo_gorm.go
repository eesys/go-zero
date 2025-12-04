package model

import (
	"context"
	"database/sql"

	"github.com/eesys/go-zero/core/stores/gormx"
)

type StudentRepo struct {
	r gormx.Repository[Student]
}

func NewStudentRepo(r gormx.Repository[Student]) StudentRepo {
	return StudentRepo{r: r}
}

type simpleResult2 struct{}

func (simpleResult2) LastInsertId() (int64, error) { return 0, nil }
func (simpleResult2) RowsAffected() (int64, error) { return 0, nil }

func (s StudentRepo) Insert(data Student) (sql.Result, error) {
	ctx := gormx.WithWrite(context.Background())
	if err := s.r.Create(ctx, &data); err != nil {
		return nil, err
	}
	return simpleResult2{}, nil
}

func (s StudentRepo) FindOne(id int64) (*Student, error) {
	ctx := gormx.WithReadReplica(context.Background())
	out, err := s.r.First(ctx, Student{Id: id})
	if err == gormx.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s StudentRepo) FindOneByClassName(class, name string) (*Student, error) {
	ctx := gormx.WithReadReplica(context.Background())
	out, err := s.r.First(ctx, map[string]any{"class": class, "name": name})
	if err == gormx.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s StudentRepo) Update(data Student) error {
	ctx := gormx.WithWrite(context.Background())
	return s.r.Save(ctx, &data)
}

func (s StudentRepo) Delete(id int64, className, studentName string) error {
	ctx := gormx.WithWrite(context.Background())
	return s.r.Delete(ctx, Student{Id: id})
}
