package svc

import (
	"dtm-zero/service/stock/internal/config"
	"dtm-zero/service/stock/internal/dal/query"
)

type ServiceContext struct {
	Config config.Config
	Query  *query.Query
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Query:  query.Use(c.MysqlDBConf.NewDriver()),
	}
}
