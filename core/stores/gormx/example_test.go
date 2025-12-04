package gormx

import (
    "context"
    "testing"
)

type user struct {
    ID   int64 `gorm:"column:id;primaryKey"`
    Name string `gorm:"column:name"`
}

func TestRepositoryCompile(t *testing.T) {
    _ = func() {
        var r Repository[user]
        ctx := context.Background()
        _, _ = r.List(ctx, user{Name: "n"}, 10, 0)
    }
}

