package gormx

import (
    "context"
    "errors"
    "gorm.io/gorm"
)

var errCantNestTx = errors.New("cannot nest transaction")

type TxKey struct{}

func WithTx(ctx context.Context, db *gorm.DB, fn func(ctx context.Context, tx *gorm.DB) error) error {
    if ctx.Value(TxKey{}) != nil {
        return errCantNestTx
    }
    return mapError(db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        ctx2 := context.WithValue(ctx, TxKey{}, true)
        return fn(ctx2, tx)
    }))
}

