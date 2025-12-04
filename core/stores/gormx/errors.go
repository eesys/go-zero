package gormx

import (
    "errors"
    "gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

func mapError(err error) error {
    if err == nil {
        return nil
    }
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return ErrNotFound
    }
    return err
}

