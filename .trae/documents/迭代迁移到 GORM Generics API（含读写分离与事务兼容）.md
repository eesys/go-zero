## 目标
- 将现有关系数据库访问从 `sqlx/sqlc` 迁移到基于 GORM 的泛型仓库 API（Repository[T]）。
- 保留读写分离、事务封装、错误语义（`ErrNotFound`）、指标与熔断逻辑的能力。
- 以小步快跑方式迭代，每周期≤10 个 Go 文件变更，单文件≤100 行改动。

## 约束
- 每迭代周期：新建或修改的 Go 文件≤10；单文件新增/修改行数≤100。
- 避免一次性替换所有模型；先并行引入 `gormx`，再逐步替换调用方。

## 当前架构（简述）
- SQL 抽象：`core/stores/sqlx`（`SqlConn/Session/TransactCtx/QueryRow*` 等）
- 缓存型 DAO：`core/stores/sqlc/CachedConn`（读缓存、索引→主键两段式）
- 生成器：`tools/goctl/model/sql`（模板包含带缓存/不带缓存模型）
- 读写分离：`sqlx` 基于 `SqlConf.Replicas/Policy` 和上下文选择（主/副）

## 迁移策略
- 并行引入 `core/stores/gormx` 包，封装：
  - 连接与 Dialector 选择（兼容 MySQL/PostgreSQL，映射 `SqlConf`）。
  - 读写分离路由（主库、只读副本，保持上下文选择语义）。
  - 事务封装（`WithTx/Transaction` 与现有 `TransactCtx` 语义对齐）。
  - 泛型仓库 `Repository[T]`（`Get/List/Create/Update/Delete/First/Where` 等）。
  - 错误映射（`ErrNotFound` 对齐）、指标/熔断包装。
- 首先迁移少量模型与它们的调用路径，验证功能与性能；随后扩大覆盖范围。
- 缓存层分阶段接入：先直连查询，随后为热点模型引入带缓存的装饰器 `CacheRepo[T]`。
- 模型结构体逐步补充 `gorm:"column:..."` 标签（或命名策略）以减少一次性改动。

## 迭代计划
### 迭代 1：引入 gormx 基础（≤9 文件，单文件≤100 行）
- 新增：
  - `core/stores/gormx/config.go`：从 `SqlConf` 构造 `gorm.Config` 与 Dialector。
  - `core/stores/gormx/db.go`：初始化 `*gorm.DB`（含连接池设置、命名策略）。
  - `core/stores/gormx/rw.go`：读写分离路由器（主库/副本选择，遵循上下文标记）。
  - `core/stores/gormx/tx.go`：事务封装（`WithTx(ctx, fn)`）、嵌套事务防护。
  - `core/stores/gormx/repository.go`：`Repository[T]` 泛型接口与默认实现。
  - `core/stores/gormx/errors.go`：统一 `ErrNotFound`、错误转换。
  - `core/stores/gormx/metrics.go`：请求时长/错误计数与熔断包装。
- 适配：
  - 保留 `sqlx/sqlc` 不动；新增极小示例用法（测试或示例，不超过 2 文件）。
- 验证：
  - 基础连通性、事务提交/回滚、读副本选择、错误语义对齐的单元测试。

### 迭代 2：迁移首批模型（≤8 文件，单文件≤100 行）
- 选择 1–2 个简单模型（仅主键查询、少量写操作）。
- 为对应结构体补充必要的 `gorm` 标签或命名策略，新增 `Repo[T]` 实例化代码。
- 替换这些模型的服务层调用为 `Repository[T]`，保留行为一致。
- 验证：
  - 与原实现对比测试（读写一致性、事务行为、错误码）。

### 迭代 3：接入缓存装饰器（≤6 文件，单文件≤100 行）
- 新增 `core/stores/gormx/cache.go`：`CacheRepo[T]`（读缓存、索引→主键两段式）。
- 为首批模型添加缓存键生成与装饰器接入（只在读路径启用）。
- 验证：
  - 命中率、失效策略（写后删除相关键）、并发一致性。

### 迭代 4：事务与批量写路径迁移（≤8 文件，单文件≤100 行）
- 将服务层中使用 `TransactCtx` 的路径迁到 `gormx.WithTx`。
- 如有批量插入，提供 `db.CreateInBatches` 或自定义批量器。
- 验证：
  - 回滚场景、并发冲突、死锁重试策略（如需要）。

### 迭代 5：扩展覆盖、沉淀模板（≤10 文件，单文件≤100 行）
- 将更多模型迁移到 `Repository[T]`。
- 在生成器中增加 GORM 模板选项（保留 goctl 原模板，新增 gorm 模板）。
- 标注 `sqlx/sqlc` 为可用但不推荐路径，文档指导新项目使用 `gormx`。

## 技术细节（关键点）
- Dialector 选择：`driver=mysql/postgres` 映射到 `gorm.io/driver/mysql|postgres`。
- 命名策略：使用 `schema.NamingStrategy` 降低对 `gorm` 标签的依赖，逐步补标签。
- 读写分离：维护两个 `*gorm.DB`（primary/replica），路由依据上下文标记（与现有 `WithReadPrimary/WithReadReplica/WithWrite` 一致）。
- 事务：`db.Transaction(func(tx *gorm.DB){...})`；在 `gormx` 暴露 `WithTx(ctx, func(r Repo[T]) error)`。
- 错误：将 `gorm.ErrRecordNotFound` 转换为统一 `ErrNotFound`；其余错误透传并计数。
- 指标与熔断：在仓库方法入口/出口打点与熔断，沿用现有 `metrics/breaker` 规范。

## 验证与度量
- 单元测试：连接、读写分离、事务、错误语义、缓存一致性。
- 基准：对比 `sqlx` 与 `gormx` 在典型查询/写入上的性能。
- 观测：请求时长与错误率仪表；缓存命中率与失效率。

## 风险与回滚
- 风险：ORM 行为差异、隐式事务、命名映射不一致、读写路由错误。
- 回滚：保留 `sqlx/sqlc` 路径；按模型回滚切换（配置开关或构造注入）。

## 交付物
- 新增 `core/stores/gormx/*` 基础库与示例。
- 首批模型的仓库迁移与服务层替换。
- 缓存装饰器与验证用例。
- 文档：使用指南、迁移注意事项、对比与最佳实践。