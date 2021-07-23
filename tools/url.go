package tools

import (
	"strings"

	"github.com/gin-gonic/gin"
)

//获取URL中批量id并解析
func IdsStrToIdsIntGroup(key string, c *gin.Context) []int {
	return idsStrToIdsIntGroup(c.Param(key))
}

func IdsStrToIdsIntGroup2(key string, c *gin.Context) []int {
	return idsStrToIdsIntGroup(c.Request.FormValue(key))
}

func idsStrToIdsIntGroup(keys string) []int {
	IDS := make([]int, 0)
	ids := strings.Split(keys, ",")
	for i := 0; i < len(ids); i++ {
		ID, _ := StringToInt(ids[i])
		IDS = append(IDS, ID)
	}
	return IDS
}
