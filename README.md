# gorm-dao

一个基于模板的 GORM DAO 代码生成工具，根据 Model 自动生成类型安全的 DAO 代码，支持链式查询操作。

## 特性

- 🚀 **自动生成**：根据 GORM Model 自动生成 DAO 代码
- 🔗 **链式操作**：支持流畅的链式查询 API
- 📝 **类型安全**：生成的代码具有完整的类型检查
- 🔍 **丰富的查询方法**：自动生成 Where、Order、Preload 等方法
- 📌 **指针类型支持**：支持 `string`、数字、`time.Time` 的指针类型（如 `*string`、`*int64`、`*time.Time`），指针字段额外生成 `IsNull` / `IsNotNull` 方法，eq 等方法的参数类型为元素类型（如 `int` 而非 `*int`）
- 💾 **CRUD 支持**：提供完整的增删改查操作
- 🔄 **事务支持**：内置事务处理功能
- 📄 **分页支持**：提供便捷的分页查询方法

## 安装

```shell
go get github.com/zl-leaf/gorm-dao
```

## 快速开始

### 1. 定义 Model

```go
package model

import "gorm.io/gorm"

type User struct {
	ID        int64  `gorm:"primaryKey"`
	Username  string `gorm:"column:username"`
	Email     string `gorm:"column:email"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

type Role struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
}
```

### 2. 生成 DAO 代码

创建一个生成脚本（例如 `cmd/gen/main.go`）：

```go
package main

import (
	gormdao "github.com/zl-leaf/gorm-dao"
	"your-project/model"
)

func main() {
	g := gormdao.NewGenerator("./gen/dao") // 指定生成目录
	g.Apply(model.User{}, model.Role{})    // 指定要生成的 Model
	err := g.Generate()
	if err != nil {
		panic(err)
	}
}
```

运行生成脚本：

```shell
go run cmd/gen/main.go
```

### 3. 使用生成的 DAO

```go
package main

import (
	"context"
	"time"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"your-project/gen/dao"
	"your-project/model"
)

func main() {
	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open("dsn"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 创建 DAO 实例
	d := dao.NewDao(db)
	ctx := context.Background()

	// 查询示例
	userList, total, err := d.UserDao(ctx).
		WhereUsernameLike("demo").
		WhereCreatedAtLt(time.Now().Unix()).
		Page(1, 10).
		Find()
	if err != nil {
		panic(err)
	}

	// 获取单条记录
	user, err := d.UserDao(ctx).
		WhereUsernameEq("admin").
		First()
	if err != nil {
		panic(err)
	}

	// 创建记录
	newUser := &model.User{
		Username:  "newuser",
		Email:     "newuser@example.com",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	err = d.UserDao(ctx).Create(newUser)
	if err != nil {
		panic(err)
	}

	// 更新记录
	user.Email = "updated@example.com"
	err = d.UserDao(ctx).Save(user)
	if err != nil {
		panic(err)
	}

	// 删除记录
	err = d.UserDao(ctx).Delete(user.ID)
	if err != nil {
		panic(err)
	}

	// 事务处理
	err = d.Transaction(ctx, func(tx *dao.Dao) error {
		// 在事务中执行多个操作
		if err := tx.UserDao(ctx).Create(newUser); err != nil {
			return err
		}
		role := &model.Role{Name: "admin"}
		if err := tx.RoleDao(ctx).Create(role); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
```

## 支持的字段类型

生成器仅对以下类型的字段生成 Where 条件方法，其它类型会跳过该字段：

| 类型 | 说明 |
|------|------|
| `string` / `*string` | 等于、不等于、Like、In、NotIn；指针时额外生成 IsNull、IsNotNull |
| `int`、`int32`、`int64`、`float32`、`float64`、`uint`、`uint32`、`uint64` 及其指针 | 等于、不等于、大于/小于、Between、In、NotIn；指针时额外生成 IsNull、IsNotNull |
| `time.Time` / `*time.Time` | 等于、大于/小于、Between；指针时额外生成 IsNull、IsNotNull |

**说明**：指针类型字段（如 `*int`）生成的 eq/gt 等方法，参数类型为**元素类型**（如 `WhereAgeEq(age int)`），便于调用时传值。

## 生成的 DAO 方法

### 查询方法

根据 Model 字段自动生成以下类型的查询方法：

- `Where{Field}Eq(value)` - 等于
- `Where{Field}Ne(value)` - 不等于
- `Where{Field}Gt(value)` - 大于
- `Where{Field}Gte(value)` - 大于等于
- `Where{Field}Lt(value)` - 小于
- `Where{Field}Lte(value)` - 小于等于
- `Where{Field}Like(value)` - LIKE 查询（字符串字段）
- `Where{Field}In(values...)` - IN 查询
- `Where{Field}NotIn(values...)` - NOT IN 查询
- `Where{Field}IsNull()` - 为空（仅指针类型字段）
- `Where{Field}IsNotNull()` - 非空（仅指针类型字段）

### 排序方法

- `OrderBy{Field}Asc()` - 升序
- `OrderBy{Field}Desc()` - 降序

### 关联查询

- `Preload{Relation}(args...)` - 预加载关联数据

### CRUD 方法

- `GetByID(id int64)` - 根据 ID 获取记录
- `First(conds ...interface{})` - 获取第一条记录
- `Find()` - 查询列表（返回列表、总数和错误）
- `Count()` - 统计数量
- `Create(record *model.{Name})` - 创建记录
- `Save(record *model.{Name})` - 保存记录
- `SaveFullAssociations(record *model.{Name})` - 保存记录及所有关联
- `Delete(conds ...interface{})` - 删除记录

### 分页方法

- `Page(page, pageSize int)` - 分页查询
- `Limit(limit int)` - 限制数量
- `Offset(offset int)` - 偏移量

## 更多示例

### 复杂查询

```go
users, total, err := d.UserDao(ctx).
	WhereUsernameLike("admin").
	WhereCreatedAtGte(startTime).
	WhereCreatedAtLt(endTime).
	OrderByCreatedAtDesc().
	Page(1, 20).
	Find()
```

### 预加载关联

```go
user, err := d.UserDao(ctx).
	PreloadRoles().
	WhereIDEq(1).
	First()
```

### 统计查询

```go
count, err := d.UserDao(ctx).
	WhereStatusEq(1).
	Count()
```

### 指针字段：IsNull / IsNotNull（可选字段）

当 Model 中字段为指针类型（如 `DeletedAt *time.Time`、`Remark *string`）时，会生成无参方法用于筛空值：

```go
// 只查未删除的（DeletedAt 为 NULL）
list, total, err := d.UserDao(ctx).
	WhereDeletedAtIsNull().
	Find()

// 只查已填备注的
list, total, err := d.UserDao(ctx).
	WhereRemarkIsNotNull().
	Page(1, 10).
	Find()
```

## 注意事项

- 生成的代码需要手动运行一次生成脚本
- Model 变更后需要重新生成 DAO 代码
- 生成的代码位于指定的输出目录中
- 确保 Model 包路径正确，以便生成的代码能正确引用

## License

MIT License - see [LICENSE](LICENSE) file for details