package gorm

import (
	"fmt"

	"template/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlClient(config config.MysqlConfig) (*gorm.DB, error) {
	// 配置MySQL连接参数
	username := config.User     // 账号
	password := config.Password // 密码
	host := config.Host         // 数据库地址，可以是Ip或者域名
	port := config.Port         // 数据库端口
	database := config.Database // 数据库名
	timeout := config.Timeout   // 连接超时，10秒

	var err error

	// 拼接下dsn参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		username, password, host, port, database, timeout)
	// 连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。

	_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接mysql数据库失败, error =" + err.Error())
	}

	// 设置数据库连接池参数
	sqlDB, _ := _db.DB()
	sqlDB.SetMaxOpenConns(config.MaxConn)     // 设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(config.MaxIdleConn) // 连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。

	return _db, nil
}
