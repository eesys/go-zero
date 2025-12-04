## 目标
- 在不改动现有 `sqlx/sqlc` 模型与测试的前提下，为 `User` 和 `Student` 增加 GORM 泛型仓库并行实现，用于验证行为与编译兼容。
- 严格遵守约束：本迭代新增/修改 Go 文件≤10 个，单文件改动≤100 行。

## 范围
- 目标模型（测试目录下，接口清晰，字段适配 ORM）：
  - `UserModel`（tools/goctl/model/sql/test/model/usermodel.go）
  - `StudentModel`（tools/goctl/model/sql/test/model/studentmodel.go）
- 保留原路径不动；新增 GORM 版本的仓库文件与最小编译测试文件。

## 变更清单（预计 4 个文件，均≤100 行）
1. 新增 `tools/goctl/model/sql/test/model/userrepo_gorm.go`
   - 类型 `UserRepo`，方法与 `UserModel` 对齐：`Insert/FindOne/FindOneByUser/FindOneByMobile/FindOneByName/Update/Delete`。
   - 依赖：`github.com/eesys/go-zero/core/stores/gormx`，持有 `gormx.Repository[User]` 与 `*gormx.Router`。
   - 读路径：`router.Read(ctx)`；写路径：`router.Write(ctx)`；错误映射：`gormx.ErrNotFound → model.ErrNotFound`。
2. 新增 `tools/goctl/model/sql/test/model/studentrepo_gorm.go`
   - 类型 `StudentRepo`，方法与 `StudentModel` 对齐：`Insert/FindOne/FindOneByClassName/Update/Delete`（本迭代不接入缓存）。
   - 依赖与读写/错误映射同上；`FindOneByClassName` 使用 `Where(map[string]any{"class": class, "name": name})`。
3. 新增 `tools/goctl/model/sql/test/model/userrepo_compile_test.go`
   - 最小编译测试：构造 `UserRepo`，调用各方法以验证编译与签名一致。
4. 新增 `tools/goctl/model/sql/test/model/studentrepo_compile_test.go`
   - 最小编译测试：构造 `StudentRepo`，调用各方法以验证编译与签名一致。

## 方法映射（关键实现）
- User
  - `Insert(u)` → `repo.Create(ctx, &u)`
  - `FindOne(id)` → `repo.First(ctx, User{ID: id})`
  - `FindOneByUser(user)` → `repo.First(ctx, User{User: user})`
  - `FindOneByMobile(mobile)` → `repo.First(ctx, User{Mobile: mobile})`
  - `FindOneByName(name)` → `repo.First(ctx, User{Name: name})`
  - `Update(u)` → `repo.Save(ctx, &u)`
  - `Delete(id)` → `repo.Delete(ctx, User{ID: id})`
- Student（无缓存版）
  - `Insert(s)` → `repo.Create(ctx, &s)`
  - `FindOne(id)` → `repo.First(ctx, Student{Id: id})`
  - `FindOneByClassName(class,name)` → `repo.First(ctx, map[string]any{"class": class, "name": name})`
  - `Update(s)` → `repo.Save(ctx, &s)`
  - `Delete(id, class, name)` → `repo.Delete(ctx, Student{Id: id})`
- 错误语义：将 `gormx.ErrNotFound` 转换为 `model.ErrNotFound`，其余错误透传。

## 读写分离与事务
- 读：对 `Find*` 方法在 `ctx` 标记 `gormx.WithReadReplica(ctx)` 后调用 `router.Read(ctx)`。
- 写：`Insert/Update/Delete` 在 `ctx` 标记 `gormx.WithWrite(ctx)` 后调用 `router.Write(ctx)`。
- 事务：仓库可在未来由服务层通过 `gormx.WithTx(ctx, db, fn)` 包裹，当前迭代保持无事务示例。

## 验证
- 编译级：新增 `*_compile_test.go` 进行最小编译验证，避免真实 DB 依赖。
- 约束核验：文件数 4 个；每文件控制在 ≤100 行。

## 风险与回滚
- 风险：字段映射不一致、读写路由误用、错误语义不符。
- 回滚：并行实现，不替换原模型；删除新增仓库文件即可撤回。

## 交付物
- 两个并行的 GORM 泛型仓库文件（User/Student）。
- 两个最小编译测试文件。