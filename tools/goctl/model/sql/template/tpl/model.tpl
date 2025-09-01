package {{.pkg}}
{{if .withCache}}
import (
	"sync"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)
{{else}}

import (
	"sync"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)
{{end}}
var (
	{{.lowerStartCamelObject}}ModelInstance {{.upperStartCamelObject}}Model
	onceFor{{.upperStartCamelObject}}Model sync.Once
)

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
	onceFor{{.upperStartCamelObject}}Model.Do(func() {
		{{.lowerStartCamelObject}}ModelInstance = &custom{{.upperStartCamelObject}}Model{
			default{{.upperStartCamelObject}}Model: new{{.upperStartCamelObject}}Model(conn, c, opts...),
		}
	})
	return {{.lowerStartCamelObject}}ModelInstance
}

func Get{{.upperStartCamelObject}}Model() {{.upperStartCamelObject}}Model {
	return {{.lowerStartCamelObject}}ModelInstance
}

{{if not .withCache}}
func (m *custom{{.upperStartCamelObject}}Model) withSession(session sqlx.Session) {{.upperStartCamelObject}}Model {
    return New{{.upperStartCamelObject}}Model(sqlx.NewSqlConnFromSession(session))
}
{{end}}

