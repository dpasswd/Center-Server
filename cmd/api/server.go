package api

import (
	"context"
	"dh-passwd/database"
	"dh-passwd/global"
	mycasbin "dh-passwd/pkg/casbin"
	"dh-passwd/pkg/logger"
	"dh-passwd/routers"
	"dh-passwd/tools"
	"dh-passwd/tools/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	configYml string
	port      string
	mode      string
	StartCmd  = &cobra.Command{
		Use:          "server",
		Short:        "Start API server",
		Example:      "dh-passwd server -c config/settings.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

var echoTimes int

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "Tcp port server listening on")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

func setup() {

	//1. 读取配置
	config.Setup(configYml)
	//2. 设置日志
	logger.Setup()
	//3. 初始化数据库链接
	database.Setup(config.DatabaseConfig.Driver)

	//4. 接口访问控制加载
	mycasbin.Setup()

	usageStr := `starting api server`
	global.Logger.Info(usageStr)

}

func run() error {
	if viper.GetString("settings.application.mode") == string(tools.ModeProd) {
		gin.SetMode(gin.ReleaseMode)
	}

	r := router.InitRouter()
	defer global.Eloquent.Close()

	srv := &http.Server{
		Addr:    config.ApplicationConfig.Host + ":" + config.ApplicationConfig.Port,
		Handler: r,
	}

	go func() {
		// 服务连接
		if config.SslConfig.Enable {
			if err := srv.ListenAndServeTLS(config.SslConfig.Pem, config.SslConfig.KeyStr); err != nil && err != http.ErrServerClosed {
				global.Logger.Fatal("listen: ", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				global.Logger.Fatal("listen: ", err)
			}
		}
	}()
	tip()
	fmt.Println(tools.Green("Server run at:"))
	fmt.Printf("-  Local:   http://localhost:%s/ \r\n", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: http://%s:%s/ \r\n", tools.GetLocaHonst(), config.ApplicationConfig.Port)
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Printf("%s Shutdown Server ... \r\n", tools.GetCurrentTimeStr())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Logger.Fatal("Server Shutdown:", err)
	}
	global.Logger.Println("Server exiting")

	return nil
}

func tip() {
	usageStr := `欢迎使用 ` + tools.Green(`这是一个测试系统 `+global.Version) + ` 可以使用 ` + tools.Red(`-h`) + ` 查看命令`
	fmt.Printf("%s \n\n", usageStr)
}
