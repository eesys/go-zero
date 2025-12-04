package model

import (
	"context"
	"testing"

	"github.com/eesys/go-zero/core/stores/gormx"
)

func TestUserRepoCompile(t *testing.T) {
	var r gormx.Repository[User]
	var router *gormx.Router
	repo := NewUserRepo(r, router)

	_, _ = repo.Insert(User{User: "u"})
	_, _ = repo.FindOne(1)
	_, _ = repo.FindOneByUser("u")
	_, _ = repo.FindOneByMobile("m")
	_, _ = repo.FindOneByName("n")
	_ = repo.Update(User{ID: 1})
	_ = repo.Delete(1)

	_ = context.Background()
}
