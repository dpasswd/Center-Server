package system

import (
	"dh-passwd/models"
	"dh-passwd/tools"
	"dh-passwd/tools/app"
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	//"log"

	"github.com/gin-gonic/gin"
)

// @Summary 获取用户详细信息
// @Description 获取JSON
// @Tags 业务信息
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/test/list [get]
// @Security Bearer
func GetUserInfo(c *gin.Context) {

	var roles = make([]string, 1)
	roles[0] = tools.GetRoleName(c)

	var permissions = make([]string, 1)
	permissions[0] = "*:*:*"

	var buttons = make([]string, 1)
	buttons[0] = "*:*:*"

	RoleMenu := models.RoleMenu{}
	RoleMenu.RoleId = tools.GetRoleId(c)

	var mp = make(map[string]interface{})
	mp["roles"] = roles

	//正式环境更改为以下代码
	//if tools.GetRoleName(c) == "admin" || tools.GetRoleName(c) == "系统管理员" {
	//	mp["permissions"] = permissions
	//	mp["buttons"] = buttons
	//} else {
	//	list, _ := RoleMenu.GetPermis()
	//	mp["permissions"] = list
	//	mp["buttons"] = list
	//}

	list, _ := RoleMenu.GetPermis()
	mp["permissions"] = list
	mp["buttons"] = list

	sysuser := models.SysUser{}
	sysuser.UserId = tools.GetUserId(c)
	user, err := sysuser.Get()
	tools.HasError(err, "", 500)

	mp["introduction"] = ""

	mp["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	if user.Avatar != "" {
		mp["avatar"] = user.Avatar
	}

	if user.FirstLogin == 1 {
		mp["password"] = user.Password
	}

	mp["firstLogin"] = user.FirstLogin
	mp["userName"] = user.Username
	mp["userId"] = user.UserId
	mp["nickName"] = user.NickName
	mp["email"] = user.Email
	mp["phone"] = user.Phone
	mp["publicKey"] = user.PublicKey
	mp["quota"] = user.Quota

	app.OK(c, mp, "")
}

// @Summary 修改用户信息
// @Description 获取JSON
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/sysbasic/info [get]
// @Security Bearer
func UpdateUserInfoByUserId(c *gin.Context) {
	var userId = tools.GetUserId(c)
	var data models.SysUser
	err := c.Bind(&data)
	data.Quota = 0
	tools.HasError(err, "数据解析失败", -1)
	//data.UpdateBy = tools.GetUserIdStr(c)
	result, err := data.UpdateByUserID(userId)
	tools.HasError(err, "用户名已存在", 500)
	app.OK(c, result, "修改成功")
}

// @Summary 添加用户
// @Description 获取JSON
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/sysbasic/info [get]
// @Security Bearer
func AddUser(c *gin.Context) {
	// mail.suffix
	var sysInfo models.SysBasic
	sysInfo, _ = sysInfo.GetSettingInfo()
	if sysInfo.RegStatus == 1 {
		var data models.SysUser
		err := c.Bind(&data)
		tools.HasError(err, "数据解析失败", -1)
		if sysInfo.EmailStatus == 1 {
			_buf := strings.Split(data.Email, "@")
			if len(_buf) == 2 {
				suffix := strings.Split(sysInfo.EmailSuffix, ",")
				if !tools.ArrayExists(_buf[1], suffix) {
					app.Error(c, 500, nil, "域名后缀不匹配")
				} else {
					var sysEmail models.SysSettings
					sysI, _ := sysEmail.GetSettingsInfo()
					var emailSys tools.EmailParam
					for _, item := range sysI {
						if item.SName == "EMail" {
							if err := json.Unmarshal([]byte(item.SData), &emailSys); err != nil {
								log.Println(err)
							}
							emailSys.Toers = data.Email
							// tools.SendEmail(emailSys, "test", "test")
						}
					}
					result, code, err := data.Insert()
					url := sysInfo.SiteDomain + "/user/activate?id=" + strconv.Itoa(result) + "&code=" + code
					if result != 0 {
						err := emailSys.SendRegisterMail(data.Email, url)
						if err != nil {
							tools.HasError(err, "邮件发送失败，请与管理员联系", 500)
							app.OK(c, result, "修改成功")
						}
					}
					tools.HasError(err, "注册失败", 500)
					app.OK(c, result, "修改成功")
				}
			}
		} else {
			var sysEmail models.SysSettings
			sysI, _ := sysEmail.GetSettingsInfo()
			var emailSys tools.EmailParam
			for _, item := range sysI {
				if item.SName == "EMail" {
					if err := json.Unmarshal([]byte(item.SData), &emailSys); err != nil {
						log.Println(err)
					}
					emailSys.Toers = data.Email
					// tools.SendEmail(emailSys, "test", "test")
				}
			}
			result, code, err := data.Insert()
			url := sysInfo.SiteDomain + "/user/activate?id=" + strconv.Itoa(result) + "&code=" + code
			if result != 0 {
				err := emailSys.SendRegisterMail(data.Email, url)
				if err != nil {
					tools.HasError(err, "邮件发送失败，请与管理员联系", 500)
					app.OK(c, result, "修改成功")
				}
			}
			tools.HasError(err, "注册失败", 500)
			app.OK(c, result, "修改成功")
		}
		//data.UpdateBy = tools.GetUserIdStr(c)
	} else {
		app.Error(c, 500, nil, "注册已关闭")
	}

}

// @Summary 激活用户
// @Description 获取JSON
// @Tags 业务信息
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/test/list [get]
// @Security Bearer
func SetActivateUser(c *gin.Context) {

	var data models.SysUser
	var err error

	userId := c.Request.FormValue("id")
	createTime, _ := base64.StdEncoding.DecodeString(c.Request.FormValue("code"))

	data.UserId, _ = strconv.Atoi(userId)
	data.CreateTime = string(createTime)
	data.Status = 1

	// err = c.Bind(&data)
	// if err == nil {
	//	tools.HasError(err, "", -1)
	// }
	result, err := data.ActivateUser()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

// @Summary 获取用户列表用于分享
// @Description 获取JSON
// @Tags 业务信息
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/test/list [get]
// @Security Bearer
func GetUserListShare(c *gin.Context) {

	var data models.SysUser
	var err error

	username := c.Request.FormValue("username")

	// err = c.Bind(&data)
	// if err == nil {
	//	tools.HasError(err, "", -1)
	// }
	result, err := data.GetUserInfoShare(username)
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}
