# go-zero 与 dtm 分布式事务框架集成

## 一、涉及技术
1. [go-zero](https://github.com/zeromicro/go-zero)
2. [dtm](https://github.com/dtm-labs/dtm)
3. [gorm-gen](https://github.com/go-gorm/gen)

## 二、目录
```shell
./
├── README.md
├── deplpy
│      ├── docker-compose.yaml  // docker-compose
│      └── init.sql // 初始化数据库
├── go.mod 
├── go.sum
├── pkg
├── restful  // HTTP
└── service  // RPC
    ├── order  // 订单服务
    └── stock  // 库存服务
```
