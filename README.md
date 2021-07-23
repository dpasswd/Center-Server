## 创建example目录

	mkdir dh-passwd
	cd dh-passwd
	go env -w GO111MODULE=on
	go env -w GOPROXY=https://goproxy.cn,direct


## 初始化dh-passwd，生成go.mod文件

	go mod init dh-passwd
	('module dh-passwd
		go 1.14)

## 拉去gin框架

	go get -u github.com/gin-gonic/gin
	完成后项目下多了个新文件《go.sum》
	文件详细罗列了当前项目直接或间接依赖的所有模块版本，并写明了那些模块版本的 SHA-256 哈希值以备 Go 在今后的操作中保证项目所依赖的那些模块版本不会被篡改

## go mod文件
	《go.mod》 文件是启用了 Go modules 的项目所必须的最重要的文件，描述了当前项目的元信息，每一行都以一个动词开头，目前有以下 5 个动词:

	module：用于定义当前项目的模块路径。
	go：用于设置预期的 Go 版本。
	require：用于设置一个特定的模块版本。
	exclude：用于从使用中排除一个特定的模块版本。
	replace：用于将一个模块版本替换为另外一个模块版本。
	indirect 的意思是传递依赖，也就是非直接依赖。

## go拉取命令行模块
	
	go get -u github.com/spf13/cobra

## 编译模块
	
	go build (注意package命名，根据自身项目定义)
## dh-passwd目录文件

	── apis(api的入口，存放router请求的方法)
	│   └── system
	│       └── test.go(当前模块可以引入models的go文件，使用数据库查询，然后返回接口)
	├── assets(静态资源)
	├── cmd(命令行自定义工具，引入cobra包，自定义启动命令,go build后可以直接使用)
	│   ├── api(web端程序启动配置）
	│   │   └── server.go
	│   ├── cobra.go(cmd引入配置）
	│   ├── config(引入config.yaml配置文件)
	│   │   └── server.go
	│   └── version(程序的版本号)
	│       └── server.go
	├── config(程序端配置文件，可以自定义端口和数据库连接)
	│   └── settings.yml
	├── data（存放日志的目录）
	│   └── logs
	│       ├── dh-passwd
	│       │   └── dh-passwd-20210409.log
	│       └── request
	├── database(引入数据库资源，mysql,pgsql,sqlite等等）
	│   ├── initialize.go
	│   ├── interface.go
	│   ├── mysql-drive.go
	│   ├── pgsql-driver.go
	│   └── sqlite3-driver.go
	├── docs(接口文档，一般引入swagger)
	├── global(全局配置，如日志global.logger引用)
	│   ├── adm.go
	│   └── logger.go
	├── dh-passwd(项目生成程序)
	├── go.mod(go的模块包)
	├── go.sum(go的依赖模块)
	├── handler(处理器，可以写task任务)
	│   └── auth.go
	├── main.go(go程序的入口，目前接入后引用到cmd模块)
	├── middleware(中间件，auth认证，header，权限认证等等)
	│   ├── auth.go
	│   ├── customerror.go
	│   ├── header.go
	│   ├── init.go
	│   ├── logger.go
	│   ├── permission.go
	│   └── requestid.go
	├── models(数据库表定义，每个go文件为一个表，在表里面定义增删改查操作)
	│   ├── login.go
	│   ├── loginlog.go
	│   ├── model.go
	│   ├── sysuser.go
	│   └── test.go
	├── pkg(存放引入包的使用)
	│   ├── casbin
	│   │   └── mycasbin.go
	│   ├── jwtauth
	│   │   └── jwtauth.go
	│   └── logger
	│       └── logger.go
	├── routers(自定义路由，可以定义需要权限认证后的，或者不需要权限认证的)
	│   ├── initrouter.go
	│   └── sysrouter.go
	├── test(单元测试模块)
	└── tools(工具，可以自定义命令，在程序接入的时候使用，比如tools.intToString转换变量类型等等)
		├── app
		│   ├── model.go
		│   ├── msg
		│   │   └── message.go
		│   └── return.go
		├── config
		│   ├── application.go
		│   ├── config.go
		│   ├── database.go
		│   ├── gen.go
		│   ├── jwt.go
		│   ├── logger.go
		│   └── ssl.go
		├── env.go
		├── file.go
		├── float64.go
		├── int.go
		├── int64.go
		├── ip.go
		├── projectDir.go
		├── string.go
		├── textcolor.go
		├── user.go
		└── utils.go
		
## 启动命令
#### ./dh-passwd config  	查看配置文件
#### ./dh-passwd version  	查看程序版本
#### ./dh-passwd help     	查看帮助

	./dh-passwd
	欢迎进入 这个测试项目 可以使用 -h 查看命令
	Error: requires at least one arg
	Usage:
	  dh-passwd [flags]
	  dh-passwd [command]

	Available Commands:
	  config      Get Application config info
	  help        Help about any command
	  init        Initialize the database
	  server      Start API server
	  version     Get version info

	Flags:
	  -h, --help   help for dh-passwd

#### 启动命令 ./dh-passwd server
	./dh-passwd server
	Start API server

	Usage:
	  dh-passwd server [flags]

	Examples:
	dh-passwd server -c config/settings.yml

	Flags:
	  -c, --config string   Start server with provided configuration file (default "config/settings.yml")
	  -h, --help            help for server
	  -m, --mode string     server mode ; eg:dev,test,prod (default "dev")
	  -p, --port string     Tcp port server listening on (default "8080")

#### 初始化数据库命令 ./dh-passwd init
	初始化数据库表
	func AutoMigrate(db *gorm.DB) error {
		db.SingularTable(true)
		err := db.AutoMigrate(new(models.LoginLog)).Error
		if err != nil {
			return err
		}
		err = db.AutoMigrate(new(models.SysUser)).Error
		if err != nil {
			return err
		}
		err = db.AutoMigrate(new(models.Test)).Error
		if err != nil {
			return err
		}

		return err
	}


## 路由配置

	无需认证的路由
	func sysRoleRouter(r *gin.RouterGroup) {
		v1 := r.Group("/api/v1")
		// 自定义方法
		registerRouter(v1)
	}
	// 路由自定义名称，自定义url地址 /api/v1/test，指向/api/v1/test/list，方法为GET
	func registerRouter(api *gin.RouterGroup) {
		r := api.Group("/test")
		{
			r.GET("/list", system.GetTest)
		}
	}
	需要认证的路由
	// 需要帐号认证的api
	func sysCheckRoleRouterInit(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
		v1 := r.Group("/api/v1")
	
		// 添加你开发的路由，经过帐号认证
		registerJobRouter(v1, authMiddleware)
	
	}
	// 路由自定义名称，自定义url地址 /api/v1/test，指向/api/v1/test/authList，方法为GET
	func registerJobRouter(api *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
		r := api.Group("/test").Use(authMiddleware.MiddlewareFunc())
		{
			// system.GetTest为指向的方法
			r.GET("/authList", system.GetTest)
		}
	}
	
## 接口定义（apis目录）
#### 定义你的方法，返回接口的内容

	（注释是在swagger中会显示）
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
		// 前端给过来的查询字段
		data.Name = c.GetString("name")
		// 进入model查询
		result, _, err := data.GetList()
		// 异常捕获
		tools.HasError(err, "获取失败", 500)
		// 返回接口信息
		app.OK(c, result, "")
	}

## 表定义 （models目录）

	package models

	import (
		orm "dh-passwd/global"
		_ "time"
	)
	// 定义你的数据表字段
	type Test struct {
		Id       int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"` //ID
		Name     string `json:"name" gorm:"name"` 					  //名称
	}
	// 定义你的数据表名称
	func (Test) TableName() string {
		return "test"
	}
	// 创建数据的方法
	func (e *Test) Create() (Test, error) {
		var doc Test
		result := orm.Eloquent.Table(e.TableName()).Create(&e)
		if result.Error != nil {
			err := result.Error
			return doc, err
		}
		doc = *e
		return doc, nil
	}

	// 获取你的单个数据
	func (e *Test) Get() (Test, error) {
		var doc Test
		table := orm.Eloquent.Table(e.TableName())

		if e.Id != 0 {
			table = table.Where("id = ?", e.Id)
		}

		if err := table.First(&doc).Error; err != nil {
			return doc, err
		}
		return doc, nil
	}
	
	// 获取你的批量数据
	func (e *Test) GetList() ([]Test, int, error) {
		var doc []Test
		var count int

		table := orm.Eloquent.Select("*").Table(e.TableName())
		if e.Id != 0 {
			table = table.Where("id = ?", e.Id)
		}

		table.Where("`deleted_at` IS NULL").Count(&count)
		if err := table.Find(&doc).Error; err != nil {
			return nil, count, err
		}
		return doc, count, nil
	}

	// 更新你的数据，通过ID
	func (e *Test) Update(id int) (update Test, err error) {
		if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
			return
		}
		//参数1:是要修改的数据
		//参数2:是修改的数据
		if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
			return
		}
		return
	}

	// 删除你的数据，根据你的ID
	func (e *Test) Delete(id int) (success bool, err error) {
		if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&Test{}).Error; err != nil {
			success = false
			return
		}
		success = true
		return
	}

	//批量删除你的数据，根据ID列表
	func (e *Test) BatchDelete(id []int) (Result bool, err error) {
		if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&Test{}).Error; err != nil {
			return
		}
		Result = true
		return
	}
