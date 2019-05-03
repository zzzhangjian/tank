package rest

import (
	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/tool/builder"
	"github.com/eyebluecn/tank/code/tool/result"
	"github.com/nu7hatch/gouuid"
	"time"
)

type UserDao struct {
	BaseDao
}

//初始化方法
func (this *UserDao) Init() {
	this.BaseDao.Init()
}

//创建用户
func (this *UserDao) Create(user *User) *User {

	if user == nil {
		panic("参数不能为nil")
	}

	timeUUID, _ := uuid.NewV4()
	user.Uuid = string(timeUUID.String())
	user.CreateTime = time.Now()
	user.UpdateTime = time.Now()
	user.LastTime = time.Now()
	user.Sort = time.Now().UnixNano() / 1e6

	db := core.CONTEXT.GetDB().Create(user)
	this.PanicError(db.Error)

	return user
}

//按照Id查询用户，找不到返回nil
func (this *UserDao) FindByUuid(uuid string) *User {

	// Read
	var user *User = &User{}
	db := core.CONTEXT.GetDB().Where(&User{Base: Base{Uuid: uuid}}).First(user)
	if db.Error != nil {
		return nil
	}
	return user
}

//按照Id查询用户,找不到抛panic
func (this *UserDao) CheckByUuid(uuid string) *User {

	if uuid == "" {
		panic("uuid必须指定")
	}

	// Read
	var user = &User{}
	db := core.CONTEXT.GetDB().Where(&User{Base: Base{Uuid: uuid}}).First(user)
	this.PanicError(db.Error)
	return user
}

//查询用户。
func (this *UserDao) FindByUsername(username string) *User {

	var user = &User{}
	db := core.CONTEXT.GetDB().Where(&User{Username: username}).First(user)
	if db.Error != nil {
		if db.Error.Error() == result.DB_ERROR_NOT_FOUND {
			return nil
		} else {
			panic(db.Error)
		}
	}
	return user
}

//显示用户列表。
func (this *UserDao) Page(page int, pageSize int, username string, email string, phone string, status string, sortArray []builder.OrderPair) *Pager {

	var wp = &builder.WherePair{}

	if username != "" {
		wp = wp.And(&builder.WherePair{Query: "username LIKE ?", Args: []interface{}{"%" + username + "%"}})
	}

	if email != "" {
		wp = wp.And(&builder.WherePair{Query: "email LIKE ?", Args: []interface{}{"%" + email + "%"}})
	}

	if phone != "" {
		wp = wp.And(&builder.WherePair{Query: "phone = ?", Args: []interface{}{phone}})
	}

	if status != "" {
		wp = wp.And(&builder.WherePair{Query: "status = ?", Args: []interface{}{status}})
	}

	count := 0
	db := core.CONTEXT.GetDB().Model(&User{}).Where(wp.Query, wp.Args...).Count(&count)
	this.PanicError(db.Error)

	var users []*User
	orderStr := this.GetSortString(sortArray)
	if orderStr == "" {
		db = core.CONTEXT.GetDB().Where(wp.Query, wp.Args...).Offset(page * pageSize).Limit(pageSize).Find(&users)
	} else {
		db = core.CONTEXT.GetDB().Where(wp.Query, wp.Args...).Order(orderStr).Offset(page * pageSize).Limit(pageSize).Find(&users)
	}

	this.PanicError(db.Error)

	pager := NewPager(page, pageSize, count, users)

	return pager
}

//查询某个用户名是否已经有用户了
func (this *UserDao) CountByUsername(username string) int {
	var count int
	db := core.CONTEXT.GetDB().
		Model(&User{}).
		Where("username = ?", username).
		Count(&count)
	this.PanicError(db.Error)
	return count
}

//保存用户
func (this *UserDao) Save(user *User) *User {

	user.UpdateTime = time.Now()
	db := core.CONTEXT.GetDB().
		Save(user)
	this.PanicError(db.Error)
	return user
}

//执行清理操作
func (this *UserDao) Cleanup() {
	this.logger.Info("[UserDao]执行清理：清除数据库中所有User记录。")
	db := core.CONTEXT.GetDB().Where("uuid is not null and role != ?", USER_ROLE_ADMINISTRATOR).Delete(User{})
	this.PanicError(db.Error)
}