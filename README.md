# XiaoMi_Gin

## 1、环境配置

- win11
- Go 1.22.2
- Gin v1.10.0
- MySql 14.14
- Redis

## 2、生成数据表

在MySql数据库中创建名为xiaomi的数据，接着运行sql文件。等数据导入成功后，需要修改conf/app.ini配置文件信息。


## 3、安装第三方包

使用Goland（或者其它IDE）打开此项目，运行main.go文件。如果第三方包没有自动安装，可以使用go mod tidy，会自动下载缺少的依赖。

## 4、运行项目

在终端输入go run main.go，打开项目。

后台地址：http://localhost:8080/admin/ 管理员账号：admin  密码：123456（可修改）

前台地址：http://localhost:8080/ 默认用户账号：15201686666  密码：123456（可自己注册多个用户）
