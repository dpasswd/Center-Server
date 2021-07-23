package share

import (
	"dh-passwd/models"
	"dh-passwd/tools"
	"dh-passwd/tools/app"

	"github.com/gin-gonic/gin"
	//"github.com/google/uuid"
)

// @Summary 根据用户id获取分享列表
// @Description 获取JSON
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/share/list [get]
// @Security Bearer
func GetShareList(c *gin.Context) {
	var data models.Share
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

	result, totalCount, err := data.GetShareListByUserId(userId, pageSize, pageNo)
	tools.HasError(err, "", -1)
	app.PageOK(c, result, totalCount, pageNo, pageSize, "")
}

// @Summary 创建分享
// @Description 获取JSON
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/share/info [put]
// @Security Bearer
func CreateShare(c *gin.Context) {
	var userId = tools.GetUserId(c)

	var data models.Share
	err := c.Bind(&data)
	if err != nil {
		tools.HasError(err, "bing fail", 500)
	}
	data.SrcUserId = userId

	err = data.InsertShareItem()
	if err != nil {
		tools.HasError(err, "记录失败", 500)
		return
	}
	var sys models.SysBasic
	result, _ := sys.GetSettingInfo()
	shareUrl := result.SiteDomain + "/publicshare/key?id=" + data.ShareId
	data.ShareNode = shareUrl
	app.OK(c, data, "分享成功")
}

// @Summary 开启或关闭分享
// @Description 获取JSON
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/share/item [post]
// @Security Bearer
func UpdateShare(c *gin.Context) {
	var userId = tools.GetUserId(c)

	var data models.Share
	err := c.Bind(&data)
	if err != nil {
		tools.HasError(err, "bing fail", 500)
	}
	data.SrcUserId = userId

	if !data.SetShareListByShare() {
		tools.HasError(err, "失败", 500)
	}

	app.OK(c, data, "成功")
}

// @Summary 根据share_id获取详细信息
// @Description 获取JSON
// @Tags 用户
// @Param username query string false "username"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/share/item [get]
// @Security Bearer
func GetShareItem(c *gin.Context) {
	var data models.Share
	var err error
	id := c.Request.FormValue("id")
	//data.DataScope = tools.GetUserIdStr(c)
	result, err := data.GetShareItemByShareId(id)
	if result.ShareId == "" {
		app.Error(c, 500, err, "select fail!")
	} else {
		app.OK(c, result, "Ok!")
	}
	// tools.HasError(err, "", -1)
}
