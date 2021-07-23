package router

import (
	"dh-passwd/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	var r *gin.Engine
	r = gin.New()
	middleware.InitMiddleware(r)
	// the jwt middleware
	authMiddleware, _ := middleware.AuthInit()

	// 注册系统路由
	InitSysRouter(r, authMiddleware)

	return r
}
