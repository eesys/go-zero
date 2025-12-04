package gormx

import "context"

const (
    readPrimaryMode readWriteMode = "read-primary"
    readReplicaMode readWriteMode = "read-replica"
    writeMode       readWriteMode = "write"
    notSpecifiedMode readWriteMode = ""
)

type readWriteModeKey struct{}
type readWriteMode string

func (m readWriteMode) isValid() bool {
    return m == readPrimaryMode || m == readReplicaMode || m == writeMode
}

func WithReadPrimary(ctx context.Context) context.Context {
    return context.WithValue(ctx, readWriteModeKey{}, readPrimaryMode)
}

func WithReadReplica(ctx context.Context) context.Context {
    return context.WithValue(ctx, readWriteModeKey{}, readReplicaMode)
}

func WithWrite(ctx context.Context) context.Context {
    return context.WithValue(ctx, readWriteModeKey{}, writeMode)
}

func getReadWriteMode(ctx context.Context) readWriteMode {
    if mode := ctx.Value(readWriteModeKey{}); mode != nil {
        if v, ok := mode.(readWriteMode); ok && v.isValid() {
            return v
        }
    }
    return notSpecifiedMode
}

func usePrimary(ctx context.Context) bool {
    return getReadWriteMode(ctx) != readReplicaMode
}

