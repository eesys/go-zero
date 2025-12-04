package gormx

import (
    "context"
    "sync/atomic"
    "gorm.io/gorm"
)

type Router struct {
    primary *gorm.DB
    replicas []*gorm.DB
    rr uint64
}

func NewRouter(primary *gorm.DB, replicas ...*gorm.DB) *Router {
    return &Router{primary: primary, replicas: replicas}
}

func (r *Router) Read(ctx context.Context) *gorm.DB {
    if len(r.replicas) == 0 || usePrimary(ctx) {
        return r.primary
    }
    i := atomic.AddUint64(&r.rr, 1)
    return r.replicas[int(i)%len(r.replicas)]
}

func (r *Router) Write(context.Context) *gorm.DB {
    return r.primary
}

