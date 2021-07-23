package models

import (
	"dh-passwd/tools"
	"dh-passwd/tools/config"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

type DataPermission struct {
	DataScope string
	UserId    int
	DeptId    int
	RoleId    int
}

func (e *DataPermission) GetDataScope(tbname string, table *gorm.DB) (*gorm.DB, error) {

	if !config.ApplicationConfig.EnableDP {
		usageStr := `数据权限已经为您` + tools.Green(`关闭`) + `，如需开启请参考配置文件字段说明`
		fmt.Printf("%s\n", usageStr)
		return table, nil
	}
	SysUser := new(SysUser)
	SysRole := new(SysRole)
	SysUser.UserId = e.UserId
	user, err := SysUser.Get()
	if err != nil {
		return nil, errors.New("获取用户数据出错 msg:" + err.Error())
	}
	SysRole.RoleId = user.RoleId
	role, err := SysRole.Get()
	if err != nil {
		return nil, errors.New("获取用户数据出错 msg:" + err.Error())
	}
	if role.DataScope == "2" {
	}
	if role.DataScope == "3" {
	}
	if role.DataScope == "4" {
		if tbname == "job" {
			table = table.Where(tbname+".game_id in (SELECT game_id from game where manager like ?)", "%"+user.Username+"%")
		} else if tbname == "game" {
			table = table.Where(tbname+".manager like ?", "%"+user.Username+"%")
		} else if tbname == "db_manage" {
			table = table.Where(tbname+".pdt_id in (SELECT cmdb_id from game where manager like ?)", "%"+user.Username+"%")
		}
	}
	if role.DataScope == "5" || role.DataScope == "" {
	}

	return table, nil
}
