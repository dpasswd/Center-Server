package models

import (
	orm "dh-passwd/global"
)

type BackupId struct {
	BackupId int `gorm:"primary_key;AUTO_INCREMENT"  json:"backupId"` // 编码
}

type BackupT struct {
	UserId     int    `gorm:"" json:"userId"`             // 用户id
	BackupName string `gorm:"size:25" json:"backupName"`  // 备份文件名
	FileSize   int    `gorm:"" json:"fileSize"`           // 文件大小
	BackupTime string `gorm:"size:128" json:"backupTime"` //备份时间
	FileMD5    string `gorm:"size:128" json:"fileMd5"`    //文件MD5
	LastTime   string `gorm:"size:128" json:"lastTime"`   //最后一次下载时间
}

type BackupTotal struct {
	Total int
	Used  int
	Free  int
}

type Backup struct {
	BackupId
	BackupT
}

func (Backup) TableName() string {
	return "dhp_backup_list"
}

// GetBackupListByUserid 根据用户列表获得备份列表
func (e *Backup) GetBackupListByUserId(userID int, pageSize int, pageNo int) ([]Backup, int, error) {
	var set []Backup
	table := orm.Eloquent.Table(e.TableName()).Select(`backup_id,user_id,backup_name,
	file_size,file_md5,backup_time,last_time`)

	table = table.Where("user_id = ?", userID)

	if err := table.Order("backup_id").Offset((pageNo - 1) * pageSize).Limit(pageSize).Error; err != nil {
		return nil, 0, err
	}
	var count int
	if err := table.Find(&set).Error; err != nil {
		return nil, 0, err
	}
	table.Where("user_id = ?", userID).Count(&count)
	return set, count, nil
}

// GetBackupListByUserid 根据用户列表获得备份大小
func (e *Backup) GetBackupSizeByUserId(userID int) (bt BackupTotal, err error) {
	table := orm.Eloquent.Table(e.TableName()).Select(`sum(file_size) as used`)

	table = table.Where("user_id = ?", userID)

	// table = table.Group("file_size")

	if err = table.First(&bt).Error; err != nil {
		return
	}
	return
}

// SetBackupItemByBackupId 修改备份信息
func (e *Backup) SetBackupItemByBackupId(userID int, backupID int) (bool, error) {
	if err := orm.Eloquent.Table(e.TableName()).Where("user_id = ? and backup_id = ?", userID, backupID).Updates(&e).Error; err != nil {
		return false, err
	}
	return true, nil
}

// SetBackupItemByBackupId 根据备份明 获取备份信息
func (e *Backup) GetBackupItemByBackupName(userID int, backupName string) bool {
	var set []Backup
	table := orm.Eloquent.Table(e.TableName()).Select(`1`)
	table = table.Where("user_id = ? and backup_name = ?", userID, backupName)
	// var count int
	if err := table.Find(&set).Error; err != nil {
		return false
	}
	if len(set) == 0 {
		return false
	}
	return true
}

// SetBackupItemByBackupId 修改最后访问时间
func (e *Backup) SetBackupItemLastTime(userID int, backupName string) bool {
	if err := orm.Eloquent.Table(e.TableName()).Where("user_id = ? and backup_name = ?", userID, backupName).Updates(&e).Error; err != nil {
		return false
	}
	return true
}

// InsertBackupItem 插入备份记录
func (e *Backup) InsertBackupItem() (err error) {
	err = orm.Eloquent.Table(e.TableName()).Create(&e).Error
	return err
}

// UpdateBackupItemByBackupId 根据备份ip删除记录
func (e *Backup) DeleteBackupItemByBackupID(userID int, backupID int) (bool, error) {
	if err := orm.Eloquent.Table(e.TableName()).Where("user_id = ? and backup_id = ?", userID, backupID).Delete(&e).Error; err != nil {
		return false, err
	}
	return true, nil
}

// BatchDeleteBackupItem 删除备份记录
func (e *Backup) BatchDeleteBackupItem(userID int, id []int) (nameList []string, Result bool, err error) {
	var set []Backup
	table := orm.Eloquent.Table(e.TableName()).Select("*")
	table = table.Where("user_id = ? and backup_id in (?)", userID, id)
	if err = table.Find(&set).Error; err != nil {
		return
	}
	// nameList := make([]string, len(set))
	for i := 0; i < len(set); i++ {
		nameList = append(nameList, set[i].BackupName)
	}
	if err = orm.Eloquent.Table(e.TableName()).Where("user_id = ? and backup_id in (?)", userID, id).Delete(&Backup{}).Error; err != nil {
		return
	}
	Result = true
	return
}
