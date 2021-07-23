package database

import (
	"dh-passwd/global"
	"dh-passwd/tools/config"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PgSql struct {
}

func (e *PgSql) Setup() {
	var err error
	var db Database

	db = new(PgSql)
	global.Source = db.GetConnect()
	log.Println(global.Source)
	global.Eloquent, err = db.Open(db.GetDriver(), db.GetConnect())
	if err != nil {
		log.Fatalf("%s connect error %v", db.GetDriver(), err)
	} else {
		log.Printf("%s connect success!", db.GetDriver())
	}

	if global.Eloquent.Error != nil {
		log.Fatalf("database error %v", global.Eloquent.Error)
	}
}

// 打开数据库连接
func (*PgSql) Open(dbType string, conn string) (db *gorm.DB, err error) {
	eloquent, err := gorm.Open(dbType, conn)
	return eloquent, err
}

func (e *PgSql) GetConnect() string {
	return config.DatabaseConfig.Source
}

func (e *PgSql) GetDriver() string {
	return config.DatabaseConfig.Driver
}
