package models

import (
	orm "dh-passwd/global"
)

type SettingsId struct {
	SettingId int `gorm:"primary_key;AUTO_INCREMENT"  json:"settingId"` // 编码
}

type SettingsT struct {
	SStatus int    `gorm:"" json:"sStatus"`       // 设置状态 1 开始  0 禁用
	SName   string `gorm:"size:128" json:"sName"` // 设置名 LDAP,EMail,dingtask
	SData   string `gorm:"size:512" json:"sData"` // 设置内容
	BaseModel
}

type SysSettings struct {
	SettingsId
	SettingsT
}

func (SysSettings) TableName() string {
	return "dhp_settings"
}

// GetSettingsInfo 获取系统基础配置信息
func (e *SysSettings) GetSettingsInfo() (sysSettings []SysSettings, err error) {
	var setlist []SysSettings
	table := orm.Eloquent.Table(e.TableName()).Select("*")

	// table = table.Where("sName = ?", e.SName)

	if err = table.Find(&setlist).Error; err != nil {
		return
	}
	return setlist, nil
}

// 添加默认数据
func AddDefaultSettingsData() (err error) {
	e := new(SysSettings)
	e.SStatus = 0
	e.SName = "dingtask"
	e.SData = "{\"status\":0,\"AppKey\":\"11111\",\"AppSecret\":\"ssssss\"}"
	err = orm.Eloquent.Table(e.TableName()).Create(&e).Error
	return err
}

//修改
func (e *SysSettings) Update(sName string) (update SysSettings, err error) {
	//参数1:是要修改的数据
	//参数2:是修改的数据
	var count int
	orm.Eloquent.Table(e.TableName()).Where("s_name = ?", sName).Count(&count)
	if count > 0 {
		if err = orm.Eloquent.Table(e.TableName()).Where("s_name = ?", sName).Updates(&e).Error; err != nil {
			return
		}
	} else {
		if err = orm.Eloquent.Table(e.TableName()).Create(&e).Error; err != nil {
			return
		}
	}

	return
}
