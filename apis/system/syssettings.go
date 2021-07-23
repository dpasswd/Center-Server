package system

import (
	"dh-passwd/models"
	"dh-passwd/tools"
	"dh-passwd/tools/app"

	"github.com/gin-gonic/gin"
	//"github.com/google/uuid"
)

// @Summary 列表基础设置信息
// @Description 获取JSON
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/sysbasic/info [get]
// @Security Bearer
func GetSettingBasic(c *gin.Context) {
	var data models.SysBasic
	var err error
	// fmt.Println(&c)
	//data.DataScope = tools.GetUserIdStr(c)
	result, err := data.GetSettingInfo()
	tools.HasError(err, "", -1)

	app.OK(c, result, "Ok！")
}

// @Summary 修改基础设置信息
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body models.SysBasic true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/sysbasic/{id} [put]
func UpdateSettingBasic(c *gin.Context) {
	var data models.SysBasic
	err := c.Bind(&data)
	tools.HasError(err, "数据解析失败", -1)
	//data.UpdateBy = tools.GetUserIdStr(c)
	result, err := data.Update(1)
	tools.HasError(err, "修改失败", 500)
	app.OK(c, result, "修改成功")
}

// @Summary 列表设置信息
// @Description 获取JSON
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/sysbasic/info [get]
// @Security Bearer
func GetSysSettingsInfo(c *gin.Context) {
	var data models.SysSettings
	var err error
	// fmt.Println(&c)
	//data.DataScope = tools.GetUserIdStr(c)
	result, err := data.GetSettingsInfo()
	tools.HasError(err, "", -1)
	app.OK(c, result, "Ok！")
}

// @Summary 列表设置信息
// @Description 获取JSON
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/sysbasic/info [get]
// @Security Bearer
func UpdateSysSettingsInfo(c *gin.Context) {
	var data models.SysSettings
	err := c.Bind(&data)
	tools.HasError(err, "数据解析失败", -1)
	result, err := data.Update(data.SName)
	tools.HasError(err, "修改失败", 500)
	app.OK(c, result, "修改成功")
}
