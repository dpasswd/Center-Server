package models

import (
	orm "dh-passwd/global"
)

type SettingId struct {
	SettingId int `gorm:"primary_key;AUTO_INCREMENT"  json:"settingId"` // 编码
}

type SettingT struct {
	SiteDomain   string `gorm:"size:128" json:"siteDomain"`  // api域名
	ApiKey       string `gorm:"size:11" json:"apiKey"`       // api密钥
	RegStatus    int    `gorm:"" json:"regStatus"`           // 注册状态 1 开启 2 关闭
	EmailStatus  int    `gorm:"" json:"emailStatus"`         //邮箱验证状态 1 开启，2 关闭
	EmailSuffix  string `gorm:"size:255" json:"emailSuffix"` //邮箱后缀
	DefaultQuota int    `gorm:"" json:"defaultQuota"`        //默认配额，单位MB
	BaseModel
}

type SysBasic struct {
	SettingId
	SettingT
}

func (SysBasic) TableName() string {
	return "dhp_setting_basic"
}

// GetSettingInfo 获取系统基础配置信息
func (e *SysBasic) GetSettingInfo() (sysBasic SysBasic, err error) {
	table := orm.Eloquent.Table(e.TableName()).Select(`site_domain,api_key,
	reg_status,email_status,email_suffix,default_quota`)

	if err = table.First(&sysBasic).Error; err != nil {
		return
	}
	return
}

// 添加默认数据
func AddDefaultSettingData() (err error) {
	e := new(SysBasic)
	e.SiteDomain = "http://localhost"
	e.ApiKey = "defaultPwd"
	e.RegStatus = 2
	e.EmailStatus = 1
	e.EmailSuffix = "qq.com,163.com"
	e.DefaultQuota = 100

	err = orm.Eloquent.Table(e.TableName()).Create(&e).Error
	return err
}

//修改
func (e *SysBasic) Update(id int) (update SysBasic, err error) {
	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Updates(&e).Error; err != nil {
		return
	}
	return
}
