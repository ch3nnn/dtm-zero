package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	OrderRpcConf zrpc.RpcClientConf
	StockRpcConf zrpc.RpcClientConf
	DTMRpcConf   zrpc.RpcClientConf
}
