package models

import (
	"crypto/md5"
	orm "dh-passwd/global"
	"dh-passwd/tools"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User
type User struct {
	// key
	IdentityKey string
	// 用户名
	UserName  string
	FirstName string
	LastName  string
	// 角色
	Role string
}

type UserName struct {
	Username string `gorm:"size:64" json:"username"`
}

type PassWord struct {
	// 密码
	Password string `gorm:"size:128" json:"password"`
}

type LoginM struct {
	UserName
	PassWord
}

type SysUserId struct {
	UserId int `gorm:"primary_key;AUTO_INCREMENT"  json:"userId"` // 编码
}

type SysUserB struct {
	NickName      string    `gorm:"size:128" json:"nickName"`          // 昵称
	Phone         string    `gorm:"size:11" json:"telephone"`          // 手机号
	RoleId        int       `gorm:"" json:"roleId"`                    // 角色编码
	Salt          string    `gorm:"size:255" json:"salt"`              //盐
	Avatar        string    `gorm:"size:255" json:"avatar"`            //头像
	Sex           string    `gorm:"size:255" json:"sex"`               //性别
	Email         string    `gorm:"size:128" json:"email"`             //邮箱
	LastLoginIP   string    `gorm:"size:128" json:"lastLoginIp"`       //最近登录ip
	LastLoginTime time.Time `gorm:"type:timestamp" json:"lastLogTime"` //最近登录时间
	FirstLogin    int       `gorm:"" json:"firstLogin"`                //是否为第一次登录
	GroupId       int       `gorm:"" json:"groupId"`                   //组编码
	PostId        int       `gorm:"" json:"postId"`                    //职位编码
	//DeptId        int       `gorm:"" json:"deptId"`                    //部门编码
	AmsKey    string `gorm:"size:128" json:"amsKey"`     //Ams的保存KEY
	IsAdmin   int    `gorm:"" json:"isAdmin"`            //是否管理员，1启用，2禁用
	Quota     int    `gorm:"" json:"Quota"`              //配额
	PublicKey string `gorm:"size:2000" json:"publicKey"` //公钥
	Lang      string `gorm:"size:255" json:"lang"`       //语言
	Status    int    `gorm:"" json:"status"`             //用户状态，1启用，2禁用
	BaseModel
}

type SimpUser struct {
	UserId    int    `gorm:"" json:"userId"`           // 编码
	NickName  string `gorm:"size:128" json:"nickName"` // 昵称
	Username  string `gorm:"size:64" json:"username"`
	PublicKey string `gorm:"size:255" json:"publicKey"` //公钥
}

type SysUser struct {
	SysUserId
	SysUserB
	LoginM
}

func (SysUser) TableName() string {
	return "dhp_user"
}

type SysUserPwd struct {
	FirstLogin  int    `json:"firstLogin"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type SysUserPage struct {
	SysUserId
	SysUserB
	LoginM
	//DeptName string `gorm:"-" json:"dept_name"`
	DeptName string `gorm:"column:dept_name"  json:"dept_name"`
	RoleName string `gorm:"column:role_name"  json:"role_name"`
}

type SysUserView struct {
	SysUserId
	SysUserB
	LoginM
	RoleName string `gorm:"column:role_name"  json:"role_name"`
}

// 获取用户数据
func (e *SysUser) Get() (SysUserView SysUserView, err error) {

	table := orm.Eloquent.Table(e.TableName()).Select([]string{"dhp_user.*", "dhp_role.role_name"})
	table = table.Joins("left join dhp_role on dhp_user.role_id=dhp_role.role_id")
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}

	if e.GroupId != 0 {
		table = table.Where("group_id = ?", e.GroupId)
	}

	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}

	if err = table.First(&SysUserView).Error; err != nil {
		return
	}

	SysUserView.Password = ""
	return
}

// GetUserInfo 根据用户名查询用户信息
func (e *SysUser) GetUserInfo() (SysUserView SysUserView, err error) {
	table := orm.Eloquent.Table(e.TableName()).Select(`user_id, nick_name, phone,
	 role_id, avatar, sex, email, group_id, post_id, ams_key, public_key, 
	 lang, status, create_time, update_time, username`)

	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if err = table.First(&SysUserView).Error; err != nil {
		return
	}
	return
}

// GetUserInfo 根据用户名查询用户信息
func (e *SysUser) GetUserAllInfo() (SysUserView SysUserView, err error) {
	table := orm.Eloquent.Table(e.TableName()).Select(`*`)

	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if err = table.First(&SysUserView).Error; err != nil {
		return
	}
	return
}

// GetUserInfoShare 获取用户精简用户列表，用于分享
func (e *SysUser) GetUserInfoShare(username string) ([]SimpUser, error) {
	var set []SimpUser
	table := orm.Eloquent.Table(e.TableName()).Select(`user_id, username, nick_name, public_key`)

	if username != "" {
		table = table.Where("username = ?", username)
	}

	if err := table.Order("user_id").Error; err != nil {
		return nil, err
	}
	// var count int
	if err := table.Find(&set).Error; err != nil {
		return nil, err
	}
	// count = len(set)
	return set, nil
}

func (e *SysUser) GetUserInfoBak() (SysUserView SysUserView, err error) {

	table := orm.Eloquent.Table(e.TableName()).Select([]string{"dhp_user.*", "dhp_role.role_name"})
	table = table.Joins("left join dhp_role on dhp_user.role_id=dhp_role.role_id")
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}

	if e.GroupId != 0 {
		table = table.Where("group_id = ?", e.GroupId)
	}

	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}

	if err = table.First(&SysUserView).Error; err != nil {
		return
	}
	return
}

func (e *SysUser) GetList() (SysUserView []SysUserView, err error) {

	table := orm.Eloquent.Table(e.TableName()).Select([]string{"dhp_user.*", "dhp_role.role_name"})
	table = table.Joins("left join dhp_role on dhp_user.role_id=dhp_role.role_id")
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}

	if e.GroupId != 0 {
		table = table.Where("group_id = ?", e.GroupId)
	}

	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}

	if err = table.Find(&SysUserView).Error; err != nil {
		return
	}
	return
}

func (e *SysUser) GetUserName(users string) (SysUserList []SysUser, err error) {
	table := orm.Eloquent.Table(e.TableName()).Select("*")
	table = table.Where("username in (?)", strings.Split(users, ","))

	if err = table.Find(&SysUserList).Error; err != nil {
		return SysUserList, err
	}
	return SysUserList, nil
}

func (e *SysUser) GetUserList() (SysUserView []SysUserView, err error) {
	table := orm.Eloquent.Table(e.TableName()).Select("user_id, username, nick_name, quota")
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}
	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}
	if err = table.Find(&SysUserView).Error; err != nil {
		return
	}
	return
}

// 返回用户详细信息列表
func (e *SysUser) GetUserInfoList(pageSize int, pageNo int) ([]SysUser, int, error) {
	var set []SysUser
	table := orm.Eloquent.Table(e.TableName()).Select(`user_id, phone, username, nick_name, avatar, sex
			, email, last_login_ip, last_login_time, is_admin, quota, public_key, role_id, lang, status`)

	if err := table.Order("user_id").Offset((pageNo - 1) * pageSize).Limit(pageSize).Error; err != nil {
		return nil, 0, err
	}
	var count int
	if err := table.Find(&set).Error; err != nil {
		return nil, 0, err
	}
	count = len(set)
	return set, count, nil
}

func (e *SysUser) GetPage(pageSize int, pageNo int) ([]SysUserPage, int, error) {
	var doc []SysUserPage
	table := orm.Eloquent.Select("dhp_user.*,sys_dept.dept_name").Table(e.TableName())
	table = table.Joins("left join sys_dept on sys_dept.dept_id = dhp_user.dept_id")

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}
	if e.Status != 0 {
		table = table.Where("dhp_user.status = ?", e.Status)
	}

	if e.Phone != "" {
		table = table.Where("dhp_user.phone = ?", e.Phone)
	}

	if e.GroupId != 0 {
		table = table.Where("dhp_user.dept_id in (select dept_id from sys_dept where dept_path like ? )", "%"+tools.IntToString(e.GroupId)+"%")
	}

	// 数据权限控制
	//dataPermission := new(DataPermission)
	//dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	//table, err := dataPermission.GetDataScope("sys_user", table)
	//if err != nil {
	//	return nil, 0, err
	//}
	var count int

	if err := table.Offset((pageNo - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("dhp_user.deleted_at IS NULL").Count(&count)
	return doc, count, nil
}

func (e *SysUser) GetFindPage(pageSize int, pageNo int, find string) ([]SysUserPage, int, error) {
	var doc []SysUserPage
	table := orm.Eloquent.Select("dhp_user.*,sys_dept.dept_name").Table(e.TableName())
	table = table.Joins("left join sys_dept on sys_dept.dept_id = dhp_user.dept_id")

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}
	if e.Status != 0 {
		table = table.Where("dhp_user.status = ?", e.Status)
	}

	if e.Phone != "" {
		table = table.Where("dhp_user.phone = ?", e.Phone)
	}

	if e.GroupId != 0 {
		table = table.Where("dhp_user.dept_id in (select dept_id from sys_dept where dept_path like ? )", "%"+tools.IntToString(e.GroupId)+"%")
	}

	if find != "" {
		table = table.Where("dhp_user.phone like ? or dhp_user.username like ?", "%"+find+"%", "%"+find+"%")
	}

	// 数据权限控制
	//dataPermission := new(DataPermission)
	//dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	//table, err := dataPermission.GetDataScope("sys_user", table)
	//if err != nil {
	//	return nil, 0, err
	//}
	var count int

	if err := table.Offset((pageNo - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("dhp_user.deleted_at IS NULL").Count(&count)
	return doc, count, nil
}

//加密
func (e *SysUser) Encrypt() (err error) {
	if e.Password == "" {
		return
	}

	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		e.Password = string(hash)
		return
	}
}

//添加
func (e SysUser) Insert() (id int, code string, err error) {
	if err = e.Encrypt(); err != nil {
		return
	}

	// check 用户名
	var count int
	orm.Eloquent.Table(e.TableName()).Where("username = ?", e.Email).Count(&count)
	if count > 0 {
		err = errors.New("账户已存在！")
		return
	}
	// 默认配置
	e.SysUserB.NickName = e.Email
	//e.SysUserB.RoleId = "admin"
	e.SysUserB.RoleId = 2
	e.SysUserB.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	e.SysUserB.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	e.SysUserB.LastLoginIP = "127.0.0.1"
	e.SysUserB.LastLoginTime = time.Now()
	e.SysUserB.Status = 2
	e.SysUserB.Avatar = "static/uploadfile/jZUIxmJycoymBprLOUbT.png"
	e.SysUserB.Lang = "zh-CN"
	e.LoginM.UserName.Username = e.Email
	e.SysUserB.IsAdmin = 2
	e.SysUserB.Quota = 10240

	// e.Encrypt()
	//添加数据
	if err = orm.Eloquent.Table(e.TableName()).Create(&e).Error; err != nil {
		return
	}
	id = e.UserId
	code = base64.StdEncoding.EncodeToString([]byte(e.CreateTime))
	return
}

func LdapCreateOrUpdateUser(data tools.UserInfo, pwd string) (user string, passwd string) {
	var err error
	var e SysUser
	var count int
	var passWd string

	//添加基础数据
	e.SysUserB.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	e.SysUserB.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	e.SysUserB.LastLoginTime = time.Now()
	e.SysUserB.Lang = "zh-CN"
	e.SysUserB.Quota = 10240
	e.SysUserB.NickName = data.Nickname
	e.SysUserB.Phone = data.Phone
	e.SysUserB.Email = data.Email
	e.Username = data.Username

	w := md5.New()
	io.WriteString(w, pwd)
	passWd = fmt.Sprintf("%x", w.Sum(nil))

	e.Password = passWd

	//加密数据
	if err = e.Encrypt(); err != nil {
		return
	}
	orm.Eloquent.Table(e.TableName()).Where("username = ?", e.Username).Count(&count)
	if count > 0 {
		if err = orm.Eloquent.Table(e.TableName()).Where("username = ?", e.Username).Update(&e).Error; err != nil {
			return
		}
	} else {
		// 创建的默认数据，更新不触发
		e.SysUserB.RoleId = 2
		e.SysUserB.IsAdmin = 2
		e.SysUserB.Status = 1
		e.SysUserB.Avatar = "static/uploadfile/jZUIxmJycoymBprLOUbT.png"

		if err = orm.Eloquent.Table(e.TableName()).Create(&e).Error; err != nil {
			return
		}
	}
	return e.Username, passWd
}

// 激活用户
func (e *SysUser) ActivateUser() (bool, error) {
	if err := orm.Eloquent.Table(e.TableName()).Where("user_id = ? and create_time = ?", e.UserId, e.CreateTime).Updates(&e).Error; err != nil {
		return false, errors.New("激活失败")
	}
	return true, nil
}

// 添加默认管理员
func AddDefaultUser() (err error) {
	e := new(SysUser)
	e.SysUserB.NickName = "admin"
	//e.SysUserB.RoleId = "admin"
	e.SysUserB.RoleId = 1
	e.SysUserB.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	e.SysUserB.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	e.SysUserB.LastLoginIP = "127.0.0.1"
	e.SysUserB.LastLoginTime = time.Now()
	e.SysUserB.Status = 1
	e.SysUserB.Avatar = "static/uploadfile/jZUIxmJycoymBprLOUbT.png"
	e.SysUserB.Lang = "zh-CN"
	e.LoginM.UserName.Username = "admin"
	e.LoginM.PassWord.Password = "21232f297a57a5a743894a0e4a801fc3"
	e.SysUserB.IsAdmin = 1
	e.SysUserB.Quota = 10240
	// su := new(SysUser)
	// su.Password = "admin"
	// su.Username = "admin"
	e.Encrypt()

	err = orm.Eloquent.Table(e.TableName()).Create(&e).Error
	return err
}

//修改
func (e *SysUser) Update(id int) (update SysUser, err error) {
	if e.Password != "" {
		if err = e.Encrypt(); err != nil {
			return
		}
	}
	if err = orm.Eloquent.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}
	if e.RoleId == 0 {
		e.RoleId = update.RoleId
	}
	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

//根据user_id 修改用户数据
func (e *SysUser) UpdateByUserID(userID int) (bool, error) {
	if err := orm.Eloquent.Table(e.TableName()).Where("user_id = ?", userID).Updates(&e).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (e *SysUser) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("user_id in (?)", id).Delete(&SysUser{}).Error; err != nil {
		return
	}
	Result = true
	return
}

// 验证老密码后修改密码
func (e *SysUser) SetPwd(pwd SysUserPwd) (Result bool, err error) {
	user, err := e.GetUserAllInfo()
	if err != nil {
		tools.HasError(err, "获取用户数据失败(代码202)", 500)
	}
	_, err = tools.CompareHashAndPassword(user.Password, pwd.OldPassword)
	if err != nil {
		if strings.Contains(err.Error(), "hashedPassword is not the hash of the given password") {
			tools.HasError(err, "密码错误(代码202)", 500)
		}
		log.Print(err)
		return
	}
	e.Password = pwd.NewPassword
	_, err = e.Update(e.UserId)
	tools.HasError(err, "更新密码失败(代码202)", 500)
	return
}

// 修改密码
func (e *SysUser) SetFirstPwd(pwd SysUserPwd) (Result bool, err error) {
	if err != nil {
		tools.HasError(err, "获取用户数据失败(代码202)", 500)
	}
	e.FirstLogin = 2
	e.Password = pwd.NewPassword
	_, err = e.Update(e.UserId)
	tools.HasError(err, "更新密码失败(代码202)", 500)
	return
}
