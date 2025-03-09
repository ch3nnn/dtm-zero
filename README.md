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
## 三、快速开始
1. docker-compose 运行 dtm、etcd、mysql
    ```shell
    cd deplpy
    docker-compose up -d
    ```
2. 编译并运行 api、rpc 服务
   - 运行 rpc 服务
   ```shell
    # RPC 订单服务   
    cd service/order
    go run order.go -c etc/order.yaml
   ```
   ```shell
    # RPC 库存服务
    cd service/stock
    go run stock.go -c etc/stock.yaml
    ```
   - 运行 api 服务
    ```shell
    # HTTP 
    cd restful
    go run order.go -c etc/order.yaml
    ```
3. 测试

   ```shell
   # 创建订单 - ✅成功
   curl --request POST \
     --url http://127.0.0.1:8888/order/create \
     --header 'Accept: */*' \
     --header 'Accept-Encoding: gzip, deflate, br' \
     --header 'Connection: keep-alive' \
     --header 'Content-Type: application/json' \
     --header 'User-Agent: PostmanRuntime-ApipostRuntime/1.1.0' \
     --data '{
       "user_id": 1,
       "goods_id": 1,
       "num": 1
   }'
   ```

   ```shell
   # 创建订单 - ❌错误 回滚
   curl --request POST \
     --url http://127.0.0.1:8888/order/create \
     --header 'Accept: */*' \
     --header 'Accept-Encoding: gzip, deflate, br' \
     --header 'Connection: keep-alive' \
     --header 'Content-Type: application/json' \
     --header 'User-Agent: PostmanRuntime-ApipostRuntime/1.1.0' \
     --data '{
       "user_id": 1,
       "goods_id": 1,
       "num": 9999
   }'
   ```


## 四、参考
- [dtm-labs/dtm-examples](https://github.com/dtm-labs/dtm-examples)
- [go-zero对接分布式事务dtm保姆式教程](https://github.com/Mikaelemmmm/gozerodtm)
- [gorm](https://gorm.io/zh_CN/docs/)
- [dtm-异常与子事务屏障](https://dtm.pub/practice/barrier.html)
- [dtm-gozero支持](https://dtm.pub/ref/gozero.html)
- [订单系统 dtm-解决方案](https://www.dtm.pub/app/order.html#dtm-%E8%A7%A3%E5%86%B3%E6%96%B9%E6%A1%88)
