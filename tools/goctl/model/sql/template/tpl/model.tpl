package {{.pkg}}
{{if .withCache}}
import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)
{{else}}

import "github.com/zeromicro/go-zero/core/stores/sqlx"
{{end}}
var {{.lowerStartCamelObject}}ModelSingleton {{.upperStartCamelObject}}Model

type (
	// {{.upperStartCamelObject}}Model is an interface to be customized, add more methods here,
	// and implement the added methods in custom{{.upperStartCamelObject}}Model.
	{{.upperStartCamelObject}}Model interface {
		{{.lowerStartCamelObject}}Model
		{{if not .withCache}}withSession(session sqlx.Session) {{.upperStartCamelObject}}Model{{end}}
	}

	custom{{.upperStartCamelObject}}Model struct {
		*default{{.upperStartCamelObject}}Model
	}
)

// New{{.upperStartCamelObject}}Model returns a model for the database table.
func New{{.upperStartCamelObject}}Model(conn sqlx.SqlConn{{if .withCache}}, c cache.CacheConf, opts ...cache.Option{{end}}) {{.upperStartCamelObject}}Model {
	if {{.lowerStartCamelObject}}ModelSingleton != nil {
		return {{.lowerStartCamelObject}}ModelSingleton
	}
	{{.lowerStartCamelObject}}ModelSingleton = &custom{{.upperStartCamelObject}}Model{
		default{{.upperStartCamelObject}}Model: new{{.upperStartCamelObject}}Model(conn, c, opts...),
	}
	return {{.lowerStartCamelObject}}ModelSingleton
}

{{if not .withCache}}
func (m *custom{{.upperStartCamelObject}}Model) withSession(session sqlx.Session) {{.upperStartCamelObject}}Model {
    return New{{.upperStartCamelObject}}Model(sqlx.NewSqlConnFromSession(session))
}
{{end}}

