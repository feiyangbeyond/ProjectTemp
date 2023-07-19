package data

import (
	"deviceback/v3/pkg/config"
	"deviceback/v3/pkg/log"

	mysql "deviceback/v3/pkg/gorm"
	rds "deviceback/v3/pkg/redis"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var DataProvider = wire.NewSet(NewData, NewTestRepo)

type Data struct {
	Rds   *redis.Client
	Mysql *gorm.DB
}

func NewData(logger log.Logger, config *config.Config) (*Data, func(), error) {
	mysqlClient, err := mysql.NewMysqlClient(config.Mysql)
	if err != nil {
		panic(err)
	}

	rdsClient, err := rds.NewRdsClient(config.Redis)
	if err != nil {
		panic(err)
	}

	cleanup := func() {
		db, _ := mysqlClient.DB()
		_ = db.Close()
		_ = rdsClient.Close()
		log.NewHelper(logger).Info("datasource cleanup")
	}

	return &Data{
		Rds:   rdsClient,
		Mysql: mysqlClient,
	}, cleanup, nil
}
