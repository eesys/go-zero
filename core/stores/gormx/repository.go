package gormx

import "context"

type Repository[T any] interface {
    First(ctx context.Context, where any) (T, error)
    List(ctx context.Context, where any, limit int, offset int) ([]T, error)
    Create(ctx context.Context, v *T) error
    Save(ctx context.Context, v *T) error
    Delete(ctx context.Context, where any) error
}

type repo[T any] struct {
    r *Router
}

func NewRepository[T any](r *Router) Repository[T] {
    return &repo[T]{r: r}
}

func (p *repo[T]) First(ctx context.Context, where any) (T, error) {
    var out T
    err := p.r.Read(ctx).WithContext(ctx).Where(where).First(&out).Error
    return out, mapError(err)
}

func (p *repo[T]) List(ctx context.Context, where any, limit int, offset int) ([]T, error) {
    var list []T
    db := p.r.Read(ctx).WithContext(ctx).Where(where)
    if limit > 0 {
        db = db.Limit(limit)
    }
    if offset > 0 {
        db = db.Offset(offset)
    }
    err := db.Find(&list).Error
    return list, mapError(err)
}

func (p *repo[T]) Create(ctx context.Context, v *T) error {
    return mapError(p.r.Write(ctx).WithContext(ctx).Create(v).Error)
}

func (p *repo[T]) Save(ctx context.Context, v *T) error {
    return mapError(p.r.Write(ctx).WithContext(ctx).Save(v).Error)
}

func (p *repo[T]) Delete(ctx context.Context, where any) error {
    var t T
    return mapError(p.r.Write(ctx).WithContext(ctx).Where(where).Delete(&t).Error)
}
