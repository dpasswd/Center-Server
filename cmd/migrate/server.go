package migrate

import (
	"dh-passwd/database"
	"dh-passwd/global"
	"dh-passwd/models"
	"dh-passwd/models/gorm"
	"dh-passwd/pkg/logger"
	"dh-passwd/tools/config"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	configYml string
	mode      string
	StartCmd  = &cobra.Command{
		Use:     "init",
		Short:   "Initialize the database",
		Example: "dh-passwd init -c config/settings.yml",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

func run() {
	usage := `start init`
	fmt.Println(usage)
	//1. 读取配置
	config.Setup(configYml)
	//2. 设置日志
	logger.Setup()
	//3. 初始化数据库链接
	database.Setup(config.DatabaseConfig.Driver)
	//4. 数据库迁移
	_ = migrateModel()
	fmt.Println("数据库结构初始化成功！")
	//5. 数据初始化完成
	if err := models.InitDb(); err != nil {
		global.Logger.Fatal("数据库基础数据初始化失败！")
	}
	usage = `数据库基础数据初始化成功`
	fmt.Println(usage)
}

func migrateModel() error {
	if config.DatabaseConfig.Driver == "mysql" {
		global.Eloquent = global.Eloquent.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	}
	return gorm.AutoMigrate(global.Eloquent)
}
