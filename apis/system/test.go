package system

import (
	"dh-passwd/models"
	"dh-passwd/tools"
	"dh-passwd/tools/app"
	"log"

	"github.com/gin-gonic/gin"
)

// @Summary 获取任务信息
// @Description 获取JSON
// @Tags 业务信息
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/test/list [get]
// @Security Bearer
func GetTest(c *gin.Context) {
	// models为数据库模块，获取数据列表
	var data models.Test
	//data.Name = c.GetString("name")
	log.Printf(c.Query("name"))
	result, _, err := data.GetList()
	tools.HasError(err, "获取失败", 500)
	app.OK(c, result, "")
}

func PostTest(c *gin.Context) {
	// models为数据库模块，获取数据列表
	var data models.Test
	//data.Name = c.GetString("name")
	log.Printf(c.PostForm("name"))
	result, _, err := data.GetList()
	tools.HasError(err, "获取失败", 500)
	app.OK(c, result, "")
}
