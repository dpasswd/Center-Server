package models

import (
	orm "dh-passwd/global"
	"log"
	"time"

	"github.com/google/uuid"
)

type Share struct {
	ShareId     string `gorm:"size:128" json:"shareId"`     // 分享ID，uuid
	SrcUserId   int    `gorm:"" json:"srcUserId"`           // 分享人
	DstUserId   int    `gorm:"" json:"dstUserId"`           // 接收人
	SaveCount   int    `gorm:"" json:"saveCount"`           //保存次数
	SecPasswd   string `gorm:"size:50000" json:"secPasswd"` //密码块
	ShareStatus int    `gorm:"" json:"shareStatus"`         // 1开启分享，2关闭分享 和EndTime 同时控制分享状态
	EndTime     string `gorm:"size:128" json:"endTime"`     //分享结束时间
	ShareTime   string `gorm:"size:128" json:"shareTime"`   //分享时间
	ShareNode   string `gorm:"size:200" json:"shareNode"`   //分享备注
	LastTime    string `gorm:"size:128" json:"lastTime"`    //最后一次下载时间
}

type ShareView struct {
	Share
	Username string `gorm:"column:username"  json:"username"`
	Email    string `gorm:"column:email"  json:"email"`
}

func (Share) TableName() string {
	return "dhp_share"
}

// InsertBackupItem 插入分享记录
func (e *Share) InsertShareItem() (err error) {
	u1, _ := uuid.NewUUID()
	e.ShareId = u1.String()
	e.ShareTime = time.Now().Format("2006-01-02 15:04:05")
	err = orm.Eloquent.Table(e.TableName()).Create(&e).Error
	return err
}

// 返回分享详细信息列表
func (e *Share) GetShareListByUserId(userId int, pageSize int, pageNo int) ([]ShareView, int, error) {
	var set []ShareView
	table := orm.Eloquent.Table(e.TableName()).Select(`dhp_share.*,dhp_user.username,dhp_user.email`)
	table = table.Joins("left join dhp_user on dhp_share.dst_user_id = dhp_user.user_id")
	table = table.Where("src_user_id = ?", userId)

	if err := table.Order("share_id").Offset((pageNo - 1) * pageSize).Limit(pageSize).Error; err != nil {
		return nil, 0, err
	}
	var count int
	if err := table.Find(&set).Error; err != nil {
		return nil, 0, err
	}
	count = len(set)
	return set, count, nil
}

// 取消共享
func (e *Share) SetShareListByShare() bool {
	if err := orm.Eloquent.Table(e.TableName()).Where("share_id = ? and src_user_id = ?", e.ShareId, e.SrcUserId).Updates(&e).Error; err != nil {
		return false
	}
	return true
}

// 根据share_id 查询分享信息
func (e *Share) GetShareItemByShareId(shareID string) (Share Share, err error) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	table := orm.Eloquent.Table(e.TableName()).Select(`*`)
	table = table.Where("share_id = ? and share_status = 1 and end_time > ?", shareID, currentTime)

	if err = table.First(&Share).Error; err != nil {
		return
	}
	Share.SaveCount = Share.SaveCount + 1
	if err := orm.Eloquent.Table(e.TableName()).Where("share_id = ?", Share.ShareId).Updates(&Share).Error; err != nil {
		log.Println(err)
	}
	return
}
