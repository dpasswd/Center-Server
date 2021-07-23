package system

import (
	"dh-passwd/global"
	"dh-passwd/models"
	"dh-passwd/tools"
	"dh-passwd/tools/app"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	//"github.com/google/uuid"
)

// @Summary 列表用户信息数据
// @Description 获取JSON
// @Tags 用户
// @Param username query string false "username"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/user/list [get]
// @Security Bearer
func GetUserList(c *gin.Context) {
	var data models.SysUser
	var err error
	//data.DataScope = tools.GetUserIdStr(c)
	result, err := data.GetUserList()
	tools.HasError(err, "", -1)

	app.OK(c, result, "Ok！")
}

// @Summary 列表用户信息数据
// @Description 获取JSON
// @Tags 用户
// @Param username query string false "username"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/sysUserList [get]
// @Security Bearer
func GetSysUserList(c *gin.Context) {
	var data models.SysUser
	var err error
	var pageSize = 10
	var pageNo = 1

	size := c.Request.FormValue("pageSize")
	if size != "" {
		pageSize = tools.StrToInt(err, size)
	}

	index := c.Request.FormValue("pageNo")
	if index != "" {
		pageNo = tools.StrToInt(err, index)
	}
	data.Username = c.Request.FormValue("username")
	data.Phone = c.Request.FormValue("phone")
	find := c.Request.FormValue("find")

	postId := c.Request.FormValue("postId")
	data.PostId, _ = tools.StringToInt(postId)

	//deptId := c.Request.FormValue("deptId")
	//data.DeptId, _ = tools.StringToInt(deptId)
	//var Dept models.SysDept
	//dept, err := Dept.Get()
	//fmt.Println(dept)
	//data.DataScope = tools.GetUserIdStr(c)
	//result, count, err := data.GetPage(pageSize, pageNo)
	result, count, err := data.GetFindPage(pageSize, pageNo, find)
	for i := range result {
		var role models.SysRole
		role.RoleId = result[i].RoleId
		roleData, _ := role.Get()
		result[i].RoleName = roleData.RoleName
	}
	tools.HasError(err, "", -1)

	app.PageOK(c, result, count, pageNo, pageSize, "")
}

// @Summary 获取用户
// @Description 获取JSON
// @Tags 用户
// @Param userId path int true "用户编码"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sysUser/{userId} [get]
// @Security Bearer
func GetSysUser(c *gin.Context) {
	var SysUser models.SysUser
	SysUser.UserId, _ = tools.StringToInt(c.Param("userId"))
	result, err := SysUser.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	var SysRole models.SysRole
	var Post models.Post
	roles, err := SysRole.GetList()
	posts, err := Post.GetList()

	postIds := make([]int, 0)
	postIds = append(postIds, result.PostId)

	roleIds := make([]int, 0)
	roleIds = append(roleIds, result.RoleId)
	app.Custum(c, gin.H{
		"code":    200,
		"data":    result,
		"postIds": postIds,
		"roleIds": roleIds,
		"roles":   roles,
		"posts":   posts,
	})
}

// @Summary 获取个人中心用户
// @Description 获取JSON
// @Tags 个人中心
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/user/profile [get]
// @Security Bearer
func GetSysUserProfile(c *gin.Context) {
	var SysUser models.SysUser
	userId := tools.GetUserIdStr(c)
	SysUser.UserId, _ = tools.StringToInt(userId)
	result, err := SysUser.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	var SysRole models.SysRole
	var Post models.Post
	//var Dept models.SysDept
	//获取角色列表
	roles, err := SysRole.GetList()
	//获取职位列表
	posts, err := Post.GetList()
	//获取部门列表
	//Dept.DeptId = result.DeptId
	//dept, err := Dept.Get()

	postIds := make([]int, 0)
	postIds = append(postIds, result.PostId)

	roleIds := make([]int, 0)
	roleIds = append(roleIds, result.RoleId)

	app.Custum(c, gin.H{
		"code":    200,
		"data":    result,
		"postIds": postIds,
		"roleIds": roleIds,
		"roles":   roles,
		"posts":   posts,
		//"dept":    dept,
	})
}

// @Summary 获取用户角色和职位
// @Description 获取JSON
// @Tags 用户
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sysUser [get]
// @Security Bearer
func GetSysUserInit(c *gin.Context) {
	var SysRole models.SysRole
	var Post models.Post
	roles, err := SysRole.GetList()
	posts, err := Post.GetList()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	mp := make(map[string]interface{}, 2)
	mp["roles"] = roles
	mp["posts"] = posts
	app.OK(c, mp, "")
}

// @Summary 创建用户
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body models.SysUser true "用户数据"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/sysUser [post]
func InsertSysUser(c *gin.Context) {
	var sysuser models.SysUser
	err := c.BindWith(&sysuser, binding.JSON)
	tools.HasError(err, "非法数据格式", 500)
	//sysuser.CreateBy = tools.GetUserIdStr(c)
	id, _, err := sysuser.Insert()
	tools.HasError(err, "添加失败", 500)
	app.OK(c, id, "添加成功")
}

// @Summary 修改用户数据
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body models.SysUser true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/sysuser/{userId} [put]
func UpdateSysUser(c *gin.Context) {
	var data models.SysUser
	err := c.Bind(&data)
	tools.HasError(err, "数据解析失败", -1)
	//data.UpdateBy = tools.GetUserIdStr(c)
	result, err := data.Update(data.UserId)
	tools.HasError(err, "修改失败", 500)
	app.OK(c, result, "修改成功")
}

// @Summary 管理员修改用户信息
// @Description 获取JSON
// @Tags 用户
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/sysbasic/info [get]
// @Security Bearer
func AdminUpdateUserInfoByUserId(c *gin.Context) {
	var data models.SysUser
	err := c.Bind(&data)
	tools.HasError(err, "数据解析失败", -1)
	//data.UpdateBy = tools.GetUserIdStr(c)
	result, err := data.UpdateByUserID(data.UserId)
	tools.HasError(err, "修改失败", 500)
	app.OK(c, result, "修改成功")
}

// @Summary 删除用户数据
// @Description 删除数据
// @Tags 用户
// @Param userId path int true "userId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/sysuser/{userId} [delete]
func DeleteSysUser(c *gin.Context) {
	var data models.SysUser
	//data.UpdateBy = tools.GetUserIdStr(c)
	IDS := tools.IdsStrToIdsIntGroup2("userId", c)
	result, err := data.BatchDelete(IDS)
	tools.HasError(err, "", 500)
	app.OK(c, result, "删除成功")
}

func ShowImage(c *gin.Context) {
	imageName := c.Query("imageName")
	// log.Println(imageName)
	c.File(imageName)
}

// @Summary 修改头像
// @Description 获取JSON
// @Tags 用户
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/user/profileAvatar [post]
func InsetSysUserAvatar(c *gin.Context) {
	file, _ := c.FormFile("file")

	id0, _ := uuid.NewDCEPerson()
	userId := tools.GetUserId(c)
	guid := uuid.NewMD5(id0, []byte(strconv.Itoa(userId))).String()
	//guid := ""
	filPath := "static/uploadfile/" + guid + ".jpg"
	contentType, err := tools.GetFileContentType(file)
	if contentType != "image/jpeg" {
		tools.HasError(err, "文件类型错误", 500)
		return
	}
	global.Logger.Debug(file.Filename)
	_ = c.SaveUploadedFile(file, filPath)
	sysuser := models.SysUser{}
	sysuser.UserId = userId
	sysuser.Avatar = filPath
	//sysuser.UpdateBy = tools.GetUserIdStr(c)
	sysuser.Update(sysuser.UserId)
	app.OK(c, filPath, "修改成功")
}

// @Summary 多文件上传方法
// @Description 获取JSON
// @Tags 用户
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/user/profileAvatar [post]
func UploadMultFile(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	userId := tools.GetUserId(c)
	guid := ""
	//guid := ""
	filPath := "static/uploadfile/" + guid + ".jpg"
	for _, file := range files {
		// contentType, err := tools.GetFileContentType(file)
		global.Logger.Debug(file.Filename)
		// 上传文件至指定目录
		_ = c.SaveUploadedFile(file, filPath)
	}
	sysuser := models.SysUser{}
	sysuser.UserId = userId
	sysuser.Avatar = "/" + filPath
	//sysuser.UpdateBy = tools.GetUserIdStr(c)
	sysuser.Update(sysuser.UserId)
	app.OK(c, filPath, "修改成功")
}

// @Summary 修改密码
// @Description 获取JSON
// @Tags 用户
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/user/profileAvatar [post]
func SysUserUpdatePwd(c *gin.Context) {
	var pwd models.SysUserPwd
	err := c.Bind(&pwd)
	tools.HasError(err, "数据解析失败", 500)
	sysuser := models.SysUser{}
	sysuser.UserId = tools.GetUserId(c)
	if pwd.FirstLogin != 1 {
		sysuser.SetPwd(pwd)
	} else {
		sysuser.SetFirstPwd(pwd)
	}
	app.OK(c, "", "密码修改成功")
}

// @Summary 获取用户列表
// @Description 获取JSON
// @Tags 用户
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/user/profileAvatar [post]
func GetAllUserList(c *gin.Context) {
	var data models.SysUser
	var pageSize = 10
	var pageNo = 1
	var err error

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = tools.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageNo"); index != "" {
		pageNo = tools.StrToInt(err, index)
	}

	// err = c.Bind(&data)
	// if err == nil {
	//	tools.HasError(err, "", -1)
	// }
	result, totalCount, err := data.GetUserInfoList(pageSize, pageNo)
	tools.HasError(err, "", -1)
	app.PageOK(c, result, totalCount, pageNo, pageSize, "")
}
