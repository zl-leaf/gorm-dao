# gorm-dao
模板生成gorm的dao

## 简介

- 根据Model生成Dao
- 生成的Dao支持简易链式操作

## 开始
### 下载gorm-dao
```shell
cd $your_project
go get github.com/zl-leaf/gorm-dao
```

### 开始生成

```go
package main

import (
	gormdao "github.com/zl-leaf/gorm-dao"
	"model"
)

func main() {
	g := gormdao.NewGenerator("./gen/dao") // 初始化，指定生成dao代码的目录
	g.Apply(model.User{}, model.Role{})    // 指定Model对象生成dao代码
	err := g.Generate()
	if err != nil {
		panic(err)
	}
}
```


## 开始使用
```go
var (
    db = *gorm.DB
    userList []*model.User
    total int64
    err error
)
userList, total, err = dao.NewDao(db).
    WhereUsernameLike("demo").
    WhereCreatedAtLt(time.Now().Unix()).
    Find()
```