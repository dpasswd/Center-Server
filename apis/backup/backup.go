package backup

import (
	"dh-passwd/models"
	"dh-passwd/tools"
	"dh-passwd/tools/app"
	"errors"
	"fmt"
	"os"
	"time"

	"dh-passwd/tools/config"

	"github.com/gin-gonic/gin"
	//"github.com/google/uuid"
)

// @Summary 根据用户id获取备份列表
// @Description 获取JSON
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/backup/list [get]
// @Security Bearer
func GetBackupList(c *gin.Context) {
	var data models.Backup
	var pageSize = 10
	var pageNo = 1
	var userId = tools.GetUserId(c)
	var err error

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = tools.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageNo"); index != "" {
		pageNo = tools.StrToInt(err, index)
	}

	result, totalCount, err := data.GetBackupListByUserId(userId, pageSize, pageNo)
	tools.HasError(err, "", -1)
	app.PageOK(c, result, totalCount, pageNo, pageSize, "")
}

// @Summary 根据用户id获取配额
// @Description 获取JSON
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/backup/size [get]
// @Security Bearer
func GetBackupSize(c *gin.Context) {
	var data models.Backup
	var userId = tools.GetUserId(c)
	var user models.SysUser
	user.UserId = userId
	var err error

	result, err := data.GetBackupSizeByUserId(userId)
	userInfo, _ := user.GetUserList()
	result.Free = userInfo[0].Quota - result.Used
	result.Total = userInfo[0].Quota
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

// @Summary 上传备份
// @Description 获取JSON
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/backup/item [post]
// @Security Bearer
func UploadBackup(c *gin.Context) {
	var userId = tools.GetUserId(c)
	var data models.Backup
	var user models.SysUser
	result, _ := data.GetBackupSizeByUserId(userId)
	userInfo, _ := user.GetUserList()
	result.Free = userInfo[0].Quota - result.Used
	if result.Free == 0 {
		err := errors.New("配额已满")
		tools.HasError(err, "upload error", 500)
		return
	} else {
		var basePath = config.UploadConfig.Path + tools.IntToString(userId) + string(os.PathSeparator)
		if _, err := os.Stat(basePath); err != nil {
			os.MkdirAll(basePath, os.ModePerm)
		}
		file, err := c.FormFile("file")
		if err != nil {
			tools.HasError(err, "参数错误", 500)
			return
		}
		bufName := "backup" + "_" + time.Now().Format("20060102150405")
		filename := basePath + bufName
		md5 := tools.GetFileMD5(filename)
		data.BackupName = bufName
		data.FileSize = int(file.Size) / 1024 /// 1024
		data.FileMD5 = md5
		data.BackupTime = time.Now().Format("2006-01-02 15:04:05")
		data.UserId = userId

		err = data.InsertBackupItem()
		if err != nil {
			tools.HasError(err, "记录失败", 500)
			return
		}
		if err := c.SaveUploadedFile(file, filename); err != nil {
			tools.HasError(err, "save", 500)
			return
		}
	}
	app.OK(c, data, "上传成功")
}

// @Summary 下载备份
// @Description 获取 file
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/backup/download [get]
// @Security Bearer
func DownloadBackup(c *gin.Context) {
	var userId = tools.GetUserId(c)
	var basePath = config.UploadConfig.Path + tools.IntToString(userId) + string(os.PathSeparator)
	fileName := c.Request.FormValue("file_name")
	var data models.Backup
	data.LastTime = time.Now().Format("2006-01-02 15:04:05")
	if res := data.GetBackupItemByBackupName(userId, fileName); !res {
		tools.HasError(nil, "下载失败1", 500)
	} else {
		data.SetBackupItemLastTime(userId, fileName)
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		c.File(basePath + fileName)
	}
}

// @Summary 删除备份
// @Description 获取 file
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/backup/item [delete]]
// @Security Bearer
func DelBackup(c *gin.Context) {
	var userId = tools.GetUserId(c)
	var basePath = config.UploadConfig.Path + tools.IntToString(userId) + string(os.PathSeparator)

	// var basePath = config.UploadConfig.Path + tools.IntToString(userId) + string(os.PathSeparator)
	var data models.Backup
	//data.UpdateBy = tools.GetUserIdStr(c)
	IDS := tools.IdsStrToIdsIntGroup2("backup_id", c)
	nameList, result, err := data.BatchDeleteBackupItem(userId, IDS)
	for i := 0; i < len(nameList); i++ {
		filename := basePath + nameList[i]
		tools.FileRemove(filename)
	}
	tools.HasError(err, "", 500)
	app.OK(c, result, "删除成功")
}
