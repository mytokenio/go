# MySQL ORM
* 查询单条记录
* 查询多条记录
* 添加单条记录
* 添加多条记录
* 更新记录
* 删除记录
* 支持事务操作

## 使用介绍
### 初始化MySQL操作
```
import (
	"os"

	"github.com/mytokenio/go/log"
	"github.com/mytokenio/go/mysql"
)

func init() {

	// init mysql
	dataSource := "write your mysql data source"
	if err := mysql.Init(dataSource); err != nil {
		log.Errorf("init mysql | %v", err)
		os.Exit(-1)
	}
}
```

### 模型代码示例
```
import (
	"github.com/mytokenio/go/mysql"
)

type User struct {
	mysql.Model `db:"-" json:"-"` // db:"-" 表示数据库无该字段
	ID          int64  `json:"id"` // 默认数据库字段与结构体字段一致
	// Username string `db:"username"` 映射数据库中的“username”
	Username    string `json:"username"`
	Address     string `json:"address"`
	Sex         int    `json:"sex"`
	Comment     string `json:"comment"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
}

// ---------------------------------------------------------------------------------------------------------------------

func NewUser() *User {
	return &User{
		Model: mysql.Model{
			TableName: "ddy_user",
		},
	}
}
```

### 获取单条数据

```
func (u *User) GetSingleByExp(exp map[string]interface{}) error {

	// 获取指定列数据
	builder := u.Select("ID, Username, Address, Sex, Comment, CreateTime, UpdateTime").Form(u.TableName).Limit(1)
	rows, err := u.SelectWhere(builder, exp)
	if err != nil {
		return err
	}

	// 加载数据
	if err := u.LoadStruct(rows, u); err != nil {
		return err
	}

	return nil
}
```

### 获取多条数据

```
func (u *User) GetAll() ([]*User, error) {
	builder := u.Select("ID, Username, Address, Sex, Comment, CreateTime, UpdateTime").Form(u.TableName)
	rows, err := u.SelectWhere(builder, nil)
	if err != nil {
		log.Errorf("select where: %v", err)
		return nil, err
	}

	// 加载数据
	users := make([]*User, 0, 20)
	if err := u.LoadStruct(rows, &users); err != nil {
		return nil, err
	}

	return users, nil
}
```

### 插入操作
* 指定列插入单条数据 ———— Insert()

    ```
    now := time.Now().Unix()
    params := map[string]interface{}{
        "Username":   "zhang_san",
        "Comment":    "zhang_san is good boy",
        "CreateTime": now,
        "UpdateTime": now,
    }
    id, err := model.NewUser().Insert(params)
    if err != nil {
        log.Errorf("insert user | %v", err)
        return
    }
    log.Infof("insert id: %d", id)
    ```

* 指定列插入多条数据 ———— MInsert()
> **调用MInsert时，不用担心数据量过大的问题，这是因为ORM底层进行了分批批量
> 插入，目前设置的每批插入500条。所以即使你调用MInsert要求插入1200条，实
> 际上会分3批进行批量插入。**

    ```
    now := time.Now().Unix()
    params1 := map[string]interface{}{
        "Username":   "zhang_san_1",
        "Comment":    "zhang_san_1 is good boy",
        "CreateTime": now, // CreateTime 在 UpdateTime 之前
        "UpdateTime": now,
    }
    params2 := map[string]interface{}{
        "Username":   "zhang_san_2",
        "Comment":    "zhang_san_2 is good boy",
        "UpdateTime": now,
        "CreateTime": now, // CreateTime 在 UpdateTime 之后
    }
    // MInsert 也支持插入一条数据。设计方式类似于Redis的 MSET、HMSET
    // MInsert 要求每个params中的“key”字段必须一样，但对于“key”的排列顺序无要求
    latestId, err := model.NewUser().MInsert(params1, params2)
    if err != nil {
        log.Errorf("insert user | %v", err)
        return
    }
    log.Infof("insert latestId: %d", latestId)
    ```

* 插入单条对象数据 ———— Insert()

    ```
    user := model.NewUser()
    user.Username = "hezhixiong"
    user.Address = "shang_hai"
    user.Sex = 1 // sex: 0.未知、1.男、2.女
    user.Comment = "Gopher boy"
    user.CreateTime = time.Now().Unit()
    user.UpdateTime = user.CreateTime
    id, err := user.Insert(user)
    if err != nil {
        log.Errorf("insert user | %v", err)
        return
    }
    log.Infof("insert id: %d", id)
    ```

* 插入多条对象数据 ———— Insert()
> **调用MInsert时，不用担心数据量过大的问题，这是因为ORM底层进行了分批批量
> 插入，目前设置的每批插入500条。所以即使你调用MInsert要求插入1200条，实
> 际上会分3批进行批量插入。**

    ```
    var user1, user2 model.User

    // 赋值需要插入的数据
    user1.Username = "hezhixiong1"
    user1.Comment = "Gopher boy"
    user1.CreateTime = time.Now().Unit()
    user1.UpdateTime = user.CreateTime

    user2.Username = "hezhixiong2"
    user2.Comment = "Gopher boy"
    user2.CreateTime = time.Now().Unit()
    user2.UpdateTime = user.CreateTime

    latestId, err := model.NewUser().Insert(&user1, &user2)
    if err != nil {
        log.Errorf("minsert user | %v", err)
        return
    }
    log.Infof("minsert latestId: %d", latestId)
    ```

* 批量插入数据 ———— BatchInsert()

    ```
    columns := []string{"Username", "Sex", "Comment"}
    values := make([]interface{}, 0, 20)
    for i:=0; i<20; i++ {
        values = append(values, []interface{}{"name", i%2, ""})
    }
    latestId, err := model.NewUser().BatchInsert(columns, values)
    if err != nil {
        log.Errorf("batch insert users | %v", err)
        return
    }
    log.Infof("batch insert users latestId: %d", latestId)
    ```

### 更新操作
* 更新指定列值数据

    ```
    user := model.NewUser()
    exp := map[string]interface{}{
        "ID=?": 1,
        "Sex=?": 0,
    }
    params := map[string]interface{}{
        "Sex": 1,
        "UpdateTime": time.Now().Unit(),
    }

    // 原生SQL为：UPDATE `ddy_user` SET `Sex`='1', `UpdateTime`='1545322066' WHERE ID='1' AND Sex='0'
    affected, err := user.Update(params, exp)
    if err != nil {
        log.Errorf("update user | %v", err)
        return
    }
    log.Infof("update user affected: %d", affected)
    ```
    ```
    // 如何实现：UPDATE `ddy_user` SET `Sex`='1', `UpdateTime`='1545322066'
    // WHERE ID='1' AND (Sex='0' OR Username='hezhixiong')
    user := model.NewUser()
    exp := map[string]map[string]interface{}{
        "AND": {
            "ID=?": 1,
        },
        "OR": {
            "Sex=?":      0,
            "Username=?": "hezhixiong",
        },
    }
    params := map[string]interface{}{
        "Sex": 1,
        "UpdateTime": time.Now().Unit(),
    }
    affected, err := user.Update(params, exp)
    if err != nil {
        log.Errorf("update user | %v", err)
        return
    }
    log.Infof("update user affected: %d", affected)
    ```

* 更新整个对象数据

    ```
        exp := map[string]interface{}{
            "ID=?": 1,
            "Sex=?": 0,
        }
        user := model.NewUser()
        user.Username = "new_name"
        user.CreateTime = time.Now().Unit()
        user.UpdateTime = user.CreateTime
        affected, err := user.Update(params, user)
        if err != nil {
            log.Errorf("update user | %v", err)
            return
        }
        log.Infof("update user affected: %d", affected)
    ```

### 删除操作

```
user := model.NewUser()
exp := map[string]interface{}{
    "ID=?": 1,
}
affected, err := user.Delete(exp)
if err != nil {
    log.Errorf("delete user | %v", err)
    return
}
log.Infof("delete affected: %d", affected)
```
