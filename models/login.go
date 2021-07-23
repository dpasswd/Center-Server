package models

import (
	orm "dh-passwd/global"
	"dh-passwd/tools"
)

type Login struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
	UUID     string `form:"UUID" json:"uuid" binding:"required"`
}

func (u *Login) GetUser() (user SysUser, role SysRole, e error) {

	e = orm.Eloquent.Table("dhp_user").Where("username = ? and status = 1", u.Username).Find(&user).Error
	if e != nil {
		return
	}
	_, e = tools.CompareHashAndPassword(user.Password, u.Password)
	if e != nil {
		return
	}
	e = orm.Eloquent.Table("dhp_role").Where("role_id = ? ", user.RoleId).First(&role).Error
	if e != nil {
		return
	}

	return
}
