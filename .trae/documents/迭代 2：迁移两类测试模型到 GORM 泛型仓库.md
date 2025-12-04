## 迭代目标
- 在不改动现有 `sqlx/sqlc` 模型与测试的前提下，新增并行的 GORM 泛型仓库实现，用于验证行为与编译兼容。
- 严格遵守约束：本迭代新增/修改 Go 文件≤10 个，单文件改动≤100 行。

## 迁移范围
- 目标模型（测试目录下，接口清晰、字段适配 ORM）：
  - `UserModel`（tools/goctl/model/sql/test/model/usermodel.go）
  - `StudentModel`（tools/goctl/model/sql/test/model/studentmodel.go）
- 保留原模型与测试不变；新增 GORM 版本的仓库文件与最小编译测试。

## 变更清单（预计 4 文件，均≤100 行）
1. 新增 `tools/goctl/model/sql/test/model/userrepo_gorm.go`
   - 类型 `UserRepo`，方法与 `UserModel` 对齐：`Insert/FindOne/FindOneByUser/FindOneByMobile/FindOneByName/Update/Delete`。
   - 依赖：`github.com/eesys/go-zero/core/stores/gormx`，持有 `gormx.Repository[User]` 与 `*gormx.Router`。
   - 读路径走 `router.Read(ctx)`、写路径走 `router.Write(ctx)`；错误映射：`gormx.ErrNotFound → model.ErrNotFound`。
2. 新增 `tools/goctl/model/sql/test/model/studentrepo_gorm.go`
   - 类型 `StudentRepo`，方法与 `StudentModel` 对齐：`Insert/FindOne/FindOneByClassName/Update/Delete`（本迭代暂不接入缓存）。
   - 依赖与读写/错误映射同上；`FindOneByClassName` 使用 `Where(map[string]any{"class": class, "name": name})`。
3. 新增 `tools/goctl/model/sql/test/model/userrepo_compile_test.go`
   - 最小编译测试：构造 `UserRepo`，调用各方法以验证编译与签名一致。
4. 新增 `tools/goctl/model/sql/test/model/studentrepo_compile_test.go`
   - 最小编译测试：构造 `StudentRepo`，调用各方法以验证编译与签名一致。

## 方法映射与上下文
- User 映射
  - `Insert(u)` → `repo.Create(ctxWithWrite, &u)`
  - `FindOne(id)` → `repo.First(ctxWithRead, User{ID: id})`
  - `FindOneByUser(user)` → `repo.First(ctxWithRead, User{User: user})`
  - `FindOneByMobile(mobile)` → `repo.First(ctxWithRead, User{Mobile: mobile})`
  - `FindOneByName(name)` → `repo.First(ctxWithRead, User{Name: name})`
  - `Update(u)` → `repo.Save(ctxWithWrite, &u)`
  - `Delete(id)` → `repo.Delete(ctxWithWrite, User{ID: id})`
- Student 映射（无缓存版）
  - `Insert(s)` → `repo.Create(ctxWithWrite, &s)`
  - `FindOne(id)` → `repo.First(ctxWithRead, Student{Id: id})`
  - `FindOneByClassName(class,name)` → `repo.First(ctxWithRead, map[string]any{"class": class, "name": name})`
  - `Update(s)` → `repo.Save(ctxWithWrite, &s)`
  - `Delete(id, class, name)` → `repo.Delete(ctxWithWrite, Student{Id: id})`
- 上下文标记：`ctxWithRead = gormx.WithReadReplica(ctx)`、`ctxWithWrite = gormx.WithWrite(ctx)`；必要时使用 `gormx.WithReadPrimary(ctx)` 强制读主库。

## 错误语义与事务
- 错误：统一将 `gormx.ErrNotFound` 转换为 `model.ErrNotFound`；其余错误透传。
- 事务：仓库可由服务层在未来通过 `gormx.WithTx(ctx, db, fn)` 包裹，当前迭代保留无事务示例。

## 验证
- 编译级：新增 `*_compile_test.go` 进行最小编译验证，避免真实 DB 依赖。
- 约束核验：文件数 4 个；每文件控制在 ≤100 行。

## 风险与回滚
- 风险：字段映射差异、读写路由误用、错误语义不一致。
- 回滚：并行实现，不替换原路径；删除新增仓库文件即可撤回。

## 交付物
- 两个并行的 GORM 泛型仓库文件（User/Student）。
- 两个最小编译测试文件。