# 📘 DDD (领域驱动设计) 实战指南

> "软件的核心是其为用户解决领域问题的能力。" —— Eric Evans

本指南旨在用通俗易懂的语言解释 DDD 的核心概念，并结合本项目代码展示如何正确落地。

---

## 1. 核心概念速览

### 🧅 分层架构 (The Layers)

想象一个洋葱，核心是最重要的，外层依赖内层。

1.  **Domain (领域层)** - **"心脏"**
    *   **职责**: 包含所有的业务逻辑、规则、状态。
    *   **特点**: **纯净**。不依赖数据库、不依赖 HTTP、不依赖任何第三方库。
    *   **包含**: 实体 (Entity)、值对象 (Value Object)、领域服务 (Domain Service)、仓储接口 (Repository Interface)。

2.  **Application (应用层)** - **"大脑/指挥官"**
    *   **职责**: 协调工作。接收外部请求，指挥领域层干活，然后保存结果。
    *   **特点**: **薄**。不包含核心业务逻辑，只负责流程编排。
    *   **包含**: 应用服务 (Application Service)、DTO (Data Transfer Object)。

3.  **Interface (接口层)** - **"嘴巴/耳朵"**
    *   **职责**: 与外部世界交互。处理 HTTP 请求、RPC 调用等。
    *   **包含**: Handler, Router, Middleware。

4.  **Infrastructure (基础设施层)** - **"手脚/工具"**
    *   **职责**: 提供具体的技术实现。
    *   **包含**: 数据库实现、Redis、微信 SDK、文件存储等。

---

## 2. 关键组件详解

### 💎 值对象 (Value Object)

**"没有身份的属性集合"**

*   **定义**: 描述事物的特征，没有唯一标识 (ID)。
*   **特点**:
    *   **不可变**: 一旦创建就不能修改。
    *   **相等性**: 属性值相同即相等 (5元人民币 == 5元人民币)。
    *   **自验证**: 创建时就保证数据是合法的。
*   **示例**: `Email`, `Phone`, `Address`, `Money`。
*   **代码**: `internal/domain/valueobject/`

```go
// ✅ 正确用法
email, err := valueobject.NewEmail("test@example.com") // 自动验证格式
user.UpdateEmail(email)
```

### 👤 实体 (Entity)

**"有身份的生命体"**

*   **定义**: 具有唯一标识 (ID)，经历生命周期变化。
*   **特点**:
    *   **有状态**: 状态会随时间改变。
    *   **有行为**: 包含修改自身状态的方法 (充血模型)。
*   **示例**: `User`, `Order`, `Product`。
*   **代码**: `internal/domain/entity/`

```go
// ✅ 充血模型示例 (建议)
func (u *User) ChangePassword(newPwd string) error {
    // 密码规则校验逻辑...
    u.Password = newPwd
    return nil
}
```

### 🧠 领域服务 (Domain Service)

**"跨实体的业务逻辑"**

*   **定义**: 当某个逻辑不属于单个实体，或者涉及多个实体交互时。
*   **职责**: 封装复杂的业务规则。
*   **示例**: `UserDomainService` (检查用户是否活跃, 权限判断), `PaymentDomainService` (转账: A扣钱, B加钱)。
*   **代码**: `internal/domain/service/`

### 📦 仓储 (Repository)

**"领域对象的家"**

*   **定义**: 模拟一个集合，用于保存和获取聚合根 (Aggregate Root)。
*   **关键点**:
    *   **接口在 Domain 层**: 定义"我们要存取什么"。
    *   **实现在 Infrastructure 层**: 定义"怎么存取" (SQL, Redis, File)。
*   **代码**: `internal/domain/repository/` (接口) vs `internal/infrastructure/persistence/` (实现)。

### 🎬 应用服务 (Application Service)

**"流程编排者"**

*   **职责**:
    1.  接收 DTO。
    2.  调用 Repository 获取领域对象。
    3.  调用领域对象或领域服务执行业务逻辑。
    4.  调用 Repository 保存状态。
    5.  返回 DTO。
*   **代码**: `internal/application/service/`

---

## 3. 实战：如何写代码？

### 场景：用户修改邮箱

#### ❌ 错误写法 (贫血模型/脚本式)

```go
// Controller/Service 混在一起
func UpdateEmail(userID int64, email string) {
    // 1. 校验格式 (业务逻辑泄露到应用层)
    if !strings.Contains(email, "@") {
        return Error
    }
    
    // 2. 直接操作数据库 (依赖实现细节)
    db.Exec("UPDATE users SET email = ? WHERE id = ?", email, userID)
}
```

#### ✅ DDD 正确写法

**Step 1: 定义值对象 (Domain)**
`internal/domain/valueobject/email.go`
```go
func NewEmail(s string) (Email, error) {
    // 校验逻辑封装在这里
    if !isValid(s) return Error
    return Email{value: s}, nil
}
```

**Step 2: 定义实体行为 (Domain)**
`internal/domain/entity/user.go`
```go
func (u *User) UpdateEmail(e valueobject.Email) {
    u.Email = e.Value()
    u.UpdatedAt = time.Now()
}
```

**Step 3: 定义仓储接口 (Domain)**
`internal/domain/repository/user_repository.go`
```go
type UserRepository interface {
    FindByID(ctx context.Context, id int64) (*entity.User, error)
    Save(ctx context.Context, user *entity.User) error
}
```

**Step 4: 编排流程 (Application)**
`internal/application/service/user_service.go`
```go
func (s *UserService) UpdateEmail(ctx context.Context, cmd UpdateEmailCommand) error {
    // 1. 转换/校验输入
    email, err := valueobject.NewEmail(cmd.Email)
    if err != nil { return err }

    // 2. 获取实体
    user, _ := s.repo.FindByID(ctx, cmd.UserID)

    // 3. 执行业务逻辑
    user.UpdateEmail(email)

    // 4. 保存状态
    return s.repo.Save(ctx, user)
}
```

---

## 4. 常见问题 (FAQ)

**Q: 为什么要有 Domain Service 和 Application Service 两个 Service?**
*   **Domain Service**: 处理**业务规则** (Business Rules)。例如："只有VIP用户才能购买此商品"。
*   **Application Service**: 处理**应用流程** (Application Flow)。例如："开启事务 -> 获取用户 -> 检查VIP(调用Domain Service) -> 扣款 -> 提交事务 -> 发送邮件"。

**Q: DTO 应该在哪里定义?**
*   DTO (Data Transfer Object) 属于 **Application 层**。
*   Domain 层**绝对不能**依赖 DTO。
*   Interface 层可以使用 DTO，或者定义自己的 View Model (VO)。

**Q: 什么时候不需要 DDD?**
*   简单的 CRUD 系统。
*   纯数据展示服务。
*   如果你的业务逻辑只是 "读数据库 ->以此 JSON 返回"，那么 DDD 可能会显得繁琐。但对于复杂的、业务规则多变的系统，DDD 能极大降低维护成本。

---

## 5. 最佳实践清单

- [ ] **依赖倒置**: Domain 层不依赖任何其他层。
- [ ] **充血模型**: 实体应该包含行为，而不仅仅是字段。
- [ ] **值对象**: 尽可能使用值对象代替基本类型 (string, int)。
- [ ] **统一语言**: 代码中的命名应该与业务人员使用的语言一致。
- [ ] **显式架构**: 目录结构应该清晰地反映分层。
