package database

import (
	"dh-passwd/global"
	"dh-passwd/tools"
	"dh-passwd/tools/config"
	"time"

	_ "github.com/go-sql-driver/mysql" //加载mysql
	"github.com/jinzhu/gorm"
)

type Mysql struct {
}

func (e *Mysql) Setup() {
	var err error
	var db Database

	db = new(Mysql)
	global.Source = db.GetConnect()
	// global.Logger.Info(tools.Green(global.Source))
	global.Eloquent, err = db.Open(db.GetDriver(), db.GetConnect())
	global.Eloquent.DB().SetConnMaxLifetime(100)
	global.Eloquent.DB().SetMaxOpenConns(100)
	global.Eloquent.DB().SetConnMaxLifetime(time.Minute)
	if err != nil {
		global.Logger.Fatal(tools.Red(db.GetDriver()+" connect error :"), err)
	} else {
		global.Logger.Info(tools.Green(db.GetDriver() + " connect success !"))
	}

	if global.Eloquent.Error != nil {
		global.Logger.Fatal(tools.Red(" database error :"), global.Eloquent.Error)
	}
	global.Eloquent.LogMode(config.LoggerConfig.EnabledDB)

}

// 打开数据库连接
func (e *Mysql) Open(dbType string, conn string) (db *gorm.DB, err error) {
	return gorm.Open(dbType, conn)
}

// 获取数据库连接
func (e *Mysql) GetConnect() string {
	return config.DatabaseConfig.Source
}

func (e *Mysql) GetDriver() string {
	return config.DatabaseConfig.Driver
}
