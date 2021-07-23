package router

import (
	"dh-passwd/apis/backup"
	"dh-passwd/apis/share"
	"dh-passwd/apis/system"
	"dh-passwd/handler"
	"dh-passwd/middleware"
	jwt "dh-passwd/pkg/jwtauth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitSysRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.RouterGroup {
	router := r.Group("")

	// 无需认证
	sysRoleRouter(router)

	sysCheckRoleRouterInit(router, authMiddleware)
	// registerTempShareRouter(router)

	return router
}

// 无需过帐号认证的api
func sysRoleRouter(r *gin.RouterGroup) {
	v1 := r.Group("/api/v1")
	// 自定义方法
	registerRouter(v1)
	r.StaticFS("static/uploadfile", http.Dir("static/uploadfile"))
}

// 需要帐号认证的api
func sysCheckRoleRouterInit(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1 := r.Group("/api/v1")
	share := r.Group("share")

	v1.POST("/auth/login", authMiddleware.LoginHandler)
	v1.POST("/register", system.AddUser)
	v1.GET("/activate", system.SetActivateUser)
	v1.GET("/refresh_token", authMiddleware.RefreshHandler)
	registerTempShareRouter(share)

	// 添加你开发的路由，经过帐号认证
	registerJobRouter(v1, authMiddleware)
	registerAuthRouter(v1, authMiddleware)
	registerPageRouter(v1, authMiddleware)
	registerUserCenterRouter(v1, authMiddleware)
	registerSettingsRouter(v1, authMiddleware)
	registerSysSettingRouter(v1, authMiddleware)
	registerBackupRouter(v1, authMiddleware)
	registerShareRouter(v1, authMiddleware)
}

// 路由自定义名称，自定义url地址 /api/v1/test，指向/api/v1/test/list，方法为GET
func registerRouter(api *gin.RouterGroup) {
	r := api.Group("/test")
	{
		r.GET("/list", system.GetTest)
		r.POST("/list", system.PostTest)
	}
}

func registerPageRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1auth := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		v1auth.GET("/deptList", system.GetDeptList)
		v1auth.GET("/deptTree", system.GetDeptTree)
		v1auth.GET("/roleList", system.GetRoleList)
		v1auth.GET("/menuList", system.GetMenuList)
	}
}

func registerUserCenterRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	user := v1.Group("/user").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		user.GET("/profile", system.GetSysUserProfile)
		user.GET("/list", system.GetUserListShare)
		user.POST("/avatar", system.InsetSysUserAvatar)
		user.POST("/info", system.UpdateUserInfoByUserId)
		user.PUT("/pwd", system.SysUserUpdatePwd)
	}
}

// 路由自定义名称，自定义url地址 /api/v1/test，指向/api/v1/test/authList，方法为GET
func registerJobRouter(api *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := api.Group("/test").Use(authMiddleware.MiddlewareFunc())
	{
		// system.GetTest为指向的方法
		r.GET("/authList", system.GetTest) // 265 485
	}
}

func registerAuthRouter(api *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := api.Group("/user").Use(authMiddleware.MiddlewareFunc())
	{
		// system.GetTest为指向的方法
		r.GET("/info", system.GetUserInfo)
		r.POST("/logout", handler.LogOut)
	}
}

// 基础系统信息管理
func registerSettingsRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1 := r.Group("/sysbasic").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		v1.GET("/info", system.GetSettingBasic)
		v1.PUT("/info", system.UpdateSettingBasic)
		v1.GET("/alluser", system.GetAllUserList)
		v1.DELETE("/users", system.DeleteSysUser)
		v1.POST("/userinfo", system.AdminUpdateUserInfoByUserId)
	}
}

// 系统信息管理
func registerSysSettingRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1 := r.Group("/syssets").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		v1.GET("/info", system.GetSysSettingsInfo)
		v1.PUT("/info", system.UpdateSysSettingsInfo)
	}
}

// 备份管理
func registerBackupRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1bak := r.Group("/backup").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		v1bak.GET("/list", backup.GetBackupList)
		v1bak.POST("/item", backup.UploadBackup)
		v1bak.GET("/download", backup.DownloadBackup)
		v1bak.DELETE("/item", backup.DelBackup)
		v1bak.GET("/size", backup.GetBackupSize)
	}
}

// 分享
func registerShareRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	aaa := r.Group("/share").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		aaa.GET("/list", share.GetShareList)
		aaa.PUT("/item", share.CreateShare)
		aaa.POST("/item", share.UpdateShare)
	}
}

// 临时分享
func registerTempShareRouter(r *gin.RouterGroup) {
	sr := r.Group("")
	{
		sr.GET("/key", share.GetShareItem)
	}
}
