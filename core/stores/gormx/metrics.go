package gormx

import (
    "context"
    "time"
)

type Observer interface {
    Observe(method string, duration time.Duration, err error)
}

type observedRepo[T any] struct {
    inner Repository[T]
    obs   Observer
}

func WithObserver[T any](inner Repository[T], obs Observer) Repository[T] {
    return &observedRepo[T]{inner: inner, obs: obs}
}

func (o *observedRepo[T]) wrap(ctx context.Context, method string, fn func() error) error {
    start := time.Now()
    err := fn()
    if o.obs != nil {
        o.obs.Observe(method, time.Since(start), err)
    }
    return err
}

func (o *observedRepo[T]) First(ctx context.Context, where any) (T, error) {
    var out T
    err := o.wrap(ctx, "First", func() error {
        var e error
        out, e = o.inner.First(ctx, where)
        return e
    })
    return out, err
}

func (o *observedRepo[T]) List(ctx context.Context, where any, limit int, offset int) ([]T, error) {
    var list []T
    err := o.wrap(ctx, "List", func() error {
        var e error
        list, e = o.inner.List(ctx, where, limit, offset)
        return e
    })
    return list, err
}

func (o *observedRepo[T]) Create(ctx context.Context, v *T) error {
    return o.wrap(ctx, "Create", func() error { return o.inner.Create(ctx, v) })
}

func (o *observedRepo[T]) Save(ctx context.Context, v *T) error {
    return o.wrap(ctx, "Save", func() error { return o.inner.Save(ctx, v) })
}

func (o *observedRepo[T]) Delete(ctx context.Context, where any) error {
    return o.wrap(ctx, "Delete", func() error { return o.inner.Delete(ctx, where) })
}

