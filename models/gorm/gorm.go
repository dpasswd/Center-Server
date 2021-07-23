package gorm

import (
	"dh-passwd/models"

	"github.com/jinzhu/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	db.SingularTable(true)
	err := db.AutoMigrate(new(models.LoginLog)).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysUser)).Error
	if err != nil {
		return err
	} else {
		models.AddDefaultUser()
	}
	err = db.AutoMigrate(new(models.SysBasic)).Error
	if err != nil {
		return err
	} else {
		models.AddDefaultSettingData()
	}
	err = db.AutoMigrate(new(models.SysSettings)).Error
	if err != nil {
		return err
	} else {
		models.AddDefaultSettingsData()
	}
	err = db.AutoMigrate(new(models.CasbinRule)).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysDept)).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.Menu)).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysRoleDept)).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.RoleMenu)).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.SysRole)).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.Test)).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.Backup)).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(new(models.Share)).Error
	if err != nil {
		return err
	}

	return err
}
