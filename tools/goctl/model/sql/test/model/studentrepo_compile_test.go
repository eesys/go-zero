package model

import (
	"testing"

	"github.com/eesys/go-zero/core/stores/gormx"
)

func TestStudentRepoCompile(t *testing.T) {
	var r gormx.Repository[Student]
	repo := NewStudentRepo(r)

	_, _ = repo.Insert(Student{Class: "c", Name: "n"})
	_, _ = repo.FindOne(1)
	_, _ = repo.FindOneByClassName("c", "n")
	_ = repo.Update(Student{Id: 1})
	_ = repo.Delete(1, "c", "n")
}
