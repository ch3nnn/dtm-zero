package config

import (
	"dtm-zero/pkg/config"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	MysqlDBConf config.DatabaseConf
}
