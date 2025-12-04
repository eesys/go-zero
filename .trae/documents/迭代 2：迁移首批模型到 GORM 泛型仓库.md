## 目标

* 在不改动现有 `sqlx/sqlc` 模型的前提下，新增并行的 GORM 泛型仓库实现，覆盖两类简单模型接口，用于验证行为与编译兼容。

* 严格遵守约束：本迭代新增/修改 Go 文件≤10 个，单文件改动≤100 行。

## 范围与选择

* 迁移对象（测试模型，接口清晰，领域字段友好）：

  * `UserModel`（tools/goctl/model/sql/test/model/usermodel.go:24）

  * `StudentModel`（tools/goctl/model/sql/test/model/studentmodel.go:27）

* 保持原模型与测试不变；新增 GORM 版本的并行实现与最小编译测试。

## 变更清单（预计 4 个 Go 文件）

1. 新增 `tools/goctl/model/sql/test/model/userrepo_gorm.go`

   * 在 `model` 包内实现 `UserRepo`，方法与 `UserModel` 一致：`Insert/FindOne/FindOneByUser/FindOneByMobile/FindOneByName/Update/Delete`。

   * 内部持有 `gormx.Repository[User]` 与 `*gormx.Router`，读路径使用 `router.Read(ctx)`，写路径使用 `router.Write(ctx)`。

   * 错误映射：`gormx.ErrNotFound` 转换为本包 `ErrNotFound`（vars.go:5）。
2. 新增 `tools/goctl/model/sql/test/model/studentrepo_gorm.go`

   * 在 `model` 包内实现 `StudentRepo`，方法与 `StudentModel` 一致（本迭代暂不接入缓存）。

   * 内部持有 `gormx.Repository[Student]` 与路由，读写分离与错误映射同上。

   * `FindOneByClassName(class,name)` 使用结构条件或链式 `Where`。
3. 新增 `tools/goctl/model/sql/test/model/userrepo_compile_test.go`

   * 最小编译测试，验证仓库接口方法签名与基本调用编译通过（不进行真实数据库操作）。
4. 新增 `tools/goctl/model/sql/test/model/studentrepo_compile_test.go`

   * 最小编译测试，验证学生仓库方法签名与基本调用编译通过。

## 技术实现要点

* 复用现有结构体：`User` 与 `Student`（usermodel.go:41、studentmodel.go:42），不新增 `gorm` 标签，依赖 `gormx` 默认命名策略（`schema.NamingStrategy{SingularTable: true}`）。

* 读写分离：使用 `gormx.WithReadReplica/WithReadPrimary/WithWrite` 上下文标记与 `Router.Read/Write` 实现（core/stores/gormx/mode.go:19、23；core/stores/gormx/rw\.go:19、27）。

* 事务：预留与服务层对接接口，仓库支持在 `gormx.WithTx` 中被调用（core/stores/gormx/tx.go:12）。

* 错误语义：将 `gormx.ErrNotFound` 映射为 `model.ErrNotFound`（tools/goctl/model/sql/test/model/vars.go:5）。

* 指标：如需，在仓库实例化时可通过 `gormx.WithObserver[T]` 包装（core/stores/gormx/metrics.go:15）。

## 验证与度量

* 编译级验证：新增的 `*_compile_test.go` 仅进行编译与最小调用验证，避免引入真实 DB 依赖。

* 行为对比：下一迭代在保留原测试的前提下，补充对比测试（同样的输入与预期输出）。

## 兼容策略

* 原 `sqlx/sqlc` 模型保持不变，服务调用可选择性切换到新仓库实现。

* 错误别名保持一致：`model.ErrNotFound` 继续作为统一对外错误。

## 风险与回滚

* 风险：字段映射差异、读写路由误用、错误语义不一致。

* 回滚：并行实现，不替换原路径；删除新仓库文件即可回滚。

## 交付物

* 新增 `UserRepo` 与 `StudentRepo` 的 GORM 泛型仓库文件各 1 个。

* 新增对应的最小编译测试各 1 个。

* 不超过 4 个 Go 文件，单文件均控制在 ≤100 行。

