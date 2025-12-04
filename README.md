# Polaris

基于 Gin + GORM + Wire + DDD 的 Go 后端项目模板

## ✨ 特性

- 🏗️ **DDD 四层架构** - 清晰的领域驱动设计分层
- 💉 **依赖注入** - 使用 Google Wire 实现编译时依赖注入
- 🔒 **类型安全** - 完整的类型定义和错误处理
- 📦 **值对象支持** - 内置常用值对象（Email、Phone、Money、Address）
- 🎯 **领域服务示例** - 展示如何正确使用领域服务
- 🔌 **基础设施解耦** - 通过接口实现依赖倒置

## 🛠 技术栈

| 分类 | 技术 |
|------|------|
| **Web 框架** | Gin |
| **数据库** | PostgreSQL |
| **ORM** | GORM |
| **缓存** | Redis |
| **日志** | Zap + Lumberjack |
| **依赖注入** | Google Wire |
| **JWT** | golang-jwt/jwt |
| **配置管理** | Viper |
| **API 文档** | Swagger |
| **架构模式** | DDD + Clean Architecture |

## 📁 项目结构

```
backend-template/
├── cmd/                          # 应用程序入口
│   └── server/
│       └── main.go
│
├── internal/                     # 内部应用代码（DDD 分层）
│   │
│   ├── domain/                   # 🔵 领域层（核心业务逻辑）
│   │   ├── entity/              # 领域实体
│   │   │   ├── user.go
│   │   │   └── app_version.go
│   │   ├── valueobject/         # 值对象 ⭐新增
│   │   │   ├── email.go
│   │   │   ├── phone.go
│   │   │   ├── money.go
│   │   │   └── address.go
│   │   ├── service/             # 领域服务 ⭐新增
│   │   │   ├── user_domain_service.go
│   │   │   ├── notification_domain_service.go
│   │   │   └── payment_domain_service.go
│   │   ├── repository/          # 仓储接口
│   │   │   ├── user_repository.go
│   │   │   └── app_version_repository.go
│   │   └── errors/              # 领域错误
│   │
│   ├── application/             # 🟢 应用层（用例编排）
│   │   ├── dto/                 # 数据传输对象
│   │   │   ├── auth_dto.go
│   │   │   └── app_version_dto.go
│   │   └── service/             # 应用服务
│   │       ├── auth_service.go
│   │       ├── wechat_service.go
│   │       ├── upload_service.go
│   │       └── app_version_service.go
│   │
│   ├── infrastructure/          # 🟡 基础设施层（技术实现）
│   │   ├── persistence/         # 持久化实现
│   │   │   ├── database.go
│   │   │   ├── user_repository_impl.go
│   │   │   ├── app_version_repository_impl.go
│   │   │   └── redis.go
│   │   ├── cache/               # 缓存实现
│   │   ├── logger/              # 日志配置
│   │   ├── config/              # 配置加载
│   │   └── wechat/              # 微信 SDK 集成
│   │
│   └── interface/               # 🟠 接口层（外部交互）
│       ├── http/
│       │   ├── handler/         # HTTP 处理器
│       │   │   ├── auth_handler.go
│       │   │   └── app_version_handler.go
│       │   └── router/          # 路由配置
│       │       └── router.go
│       └── middleware/          # 中间件
│           ├── auth.go
│           ├── cors.go
│           └── logger.go
│
├── pkg/                         # 公共包（可跨项目复用）
│   ├── errors/                  # 错误定义
│   ├── response/                # 响应封装
│   ├── snowflake/               # 雪花 ID 生成器
│   └── utils/                   # 工具函数
│
├── wire/                        # Wire 依赖注入配置
│   ├── wire.go                  # Wire 定义文件
│   ├── wire_gen.go              # Wire 生成代码
│   └── app.go                   # 应用组装
│
├── config/                      # 配置文件
│   └── config.yaml
│
├── migrations/                  # 数据库迁移脚本
│   └── sql/
│
├── docs/                        # 文档
│   ├── DDD-GUIDE.md            # DDD 教程（必读）⭐
│   └── swagger/                 # Swagger API 文档
│
├── go.mod
├── go.sum
├── Makefile
└── README.md
```


## 🏗️ DDD 四层架构说明

| 层 | 职责 | 依赖方向 | 示例 |
|----|------|---------|------|
| **Interface 层** | 处理外部请求，转换为 DTO | → Application | HTTP Handler、gRPC Server |
| **Application 层** | 编排用例流程，协调领域对象 | → Domain | AuthService、OrderService |
| **Domain 层** | 核心业务逻辑和规则 | 不依赖外层 | Entity、ValueObject、DomainService |
| **Infrastructure 层** | 技术实现细节 | 实现 Domain 接口 | Database、Redis、第三方API |

### 依赖规则

```
┌─────────────┐
│  Interface  │ ──┐
└─────────────┘   │
┌─────────────┐   │
│ Application │ ──┤  都依赖
└─────────────┘   │    ↓
┌─────────────┐   │
│   Domain    │ ←─┘  核心层（不依赖任何外层）
└─────────────┘
       ↑
       │ 实现接口
┌─────────────┐
│Infrastructure│
└─────────────┘
```

## 📚 学习资源

### 必读文档

1. **[DDD 实战教程](docs/DDD-GUIDE.md)** ⭐ - 从零开始学习 DDD
2. **[API 文档](docs/swagger/)** - Swagger UI

### 建议阅读顺序

```
1️⃣ 阅读 DDD-GUIDE.md 理解核心概念
    ↓
2️⃣ 查看 internal/domain/valueobject/ 学习值对象
    ↓
3️⃣ 查看 internal/domain/service/ 学习领域服务
    ↓
4️⃣ 查看 internal/application/service/ 对比应用服务
    ↓
5️⃣ 完整走查一个用例：登录流程
   Handler → AppService → DomainService → Repository
```

## 🚀 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 安装工具

```bash
make install-tools
```

### 3. 生成依赖注入代码

```bash
make wire
```

### 4. 配置数据库

编辑 `config/config.yaml`：

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: backend_template
```

### 5. 运行数据库迁移

```bash
make migrate-up
```

### 6. 启动服务

```bash
make run
```

服务默认运行在 `http://localhost:8080`

## 🛠️ 开发命令

| 命令 | 说明 |
|------|------|
| `make wire` | 生成 Wire 依赖注入代码 |
| `make swag` | 生成 Swagger API 文档 |
| `make run` | 启动开发服务器 |
| `make test` | 运行测试 |
| `make migrate-up` | 执行数据库迁移 |
| `make migrate-down` | 回滚数据库迁移 |
| `make fmt` | 格式化代码 |
| `make lint` | 代码检查 |
| `make clean` | 清理构建产物 |
| `make help` | 查看所有命令 |

## ⚠️ 开发注意事项

### DDD 最佳实践

1. **值对象优先**
   ```go
   // ❌ 不好的做法
   type User struct {
       Email string
   }
   
   // ✅ 好的做法
   type User struct {
       Email valueobject.Email
   }
   ```

2. **领域服务 vs 应用服务**
   ```go
   // 领域服务 - 纯业务规则
   func (s *UserDomainService) IsUserActive(userID int64) bool
   
   // 应用服务 - 流程编排
   func (s *AuthService) WechatLogin(req *dto.LoginRequest) (*dto.LoginResponse, error)
   ```

3. **仓储接口定义在 Domain 层**
   ```go
   // ✅ 定义在 domain/repository/
   type UserRepository interface {
       FindByID(ctx context.Context, id int64) (*entity.User, error)
   }
   
   // ✅ 实现在 infrastructure/persistence/
   type UserRepositoryImpl struct { ... }
   ```

4. **DTO 只在 Application 层使用**
   ```go
   // ❌ 不要在 Domain 层使用 DTO
   func (s *UserDomainService) Create(dto *dto.UserDTO)
   
   // ✅ Domain 层只操作领域对象
   func (s *UserDomainService) Validate(user *entity.User)
   ```

### 常见陷阱

- ❌ **贫血模型** - 实体只有字段没有行为
- ❌ **层级混乱** - Application 层调用 Infrastructure 层
- ❌ **过度设计** - 简单的 CRUD 不需要领域服务
- ✅ **合理使用** - 复杂业务逻辑才用完整 DDD

## 📖 延伸阅读

- [Domain-Driven Design (Eric Evans)](https://www.domainlanguage.com/ddd/)
- [Implementing Domain-Driven Design (Vaughn Vernon)](https://vaughnvernon.com/)
- [Clean Architecture (Robert C. Martin)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 📄 License

MIT License

## 🤝 Contributing

欢迎提交 Issue 和 Pull Request！
