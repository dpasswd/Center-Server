package middleware

import (
	"dh-passwd/global"
	config2 "dh-passwd/tools/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		fmt.Printf("%s [INFO] %s %s %3d %13v %15s \r\n",
			startTime.Format("2006-01-02 15:04:05"),
			reqMethod,
			reqUri,
			statusCode,
			latencyTime,
			clientIP,
		)

		global.RequestLogger.Info(statusCode, latencyTime, clientIP, reqMethod, reqUri)

		if c.Request.Method != "GET" && c.Request.Method != "OPTIONS" && config2.LoggerConfig.EnabledDB {
			//SetDBOperLog(c, clientIP, statusCode, reqUri, reqMethod, latencyTime)
		}
	}
}
