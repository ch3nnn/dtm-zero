package logic

import (
	"context"
	"fmt"

	"dtm-zero/restful/internal/svc"
	"dtm-zero/restful/internal/types"
	"dtm-zero/service/order/client/order"
	orderpb "dtm-zero/service/order/pb"
	"dtm-zero/service/stock/client/stock"
	stockpb "dtm-zero/service/stock/pb"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/dtm-labs/logger"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateLogic 创建订单
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx).WithFields(
			logx.Field("service", svcCtx.Config.Name),
			logx.Field("method", "order.Create"),
		),
	}
}

func (l *CreateLogic) Create(req *types.OrderCreateReq) (resp *types.OrderCreateResp, err error) {
	// 获取订单服务地址
	orderRpcTarget, err := l.svcCtx.Config.OrderRpcConf.BuildTarget()
	if err != nil {
		return nil, fmt.Errorf("下单异常超时")
	}

	// 获取库存服务地址
	stockRpcTarget, err := l.svcCtx.Config.StockRpcConf.BuildTarget()
	if err != nil {
		return nil, fmt.Errorf("下单异常超时")
	}

	// 获取DTM服务地址
	dtmRpcTarget, err := l.svcCtx.Config.DTMRpcConf.BuildTarget()
	if err != nil {
		return nil, fmt.Errorf("下单异常超时")
	}

	// 创建一个saga事务
	m := dtmgrpc.NewSagaGrpc(dtmRpcTarget, dtmgrpc.MustGenGid(dtmRpcTarget)).
		// 添加订单服务
		Add(
			orderRpcTarget+orderpb.Order_Create_FullMethodName,         // 创建订单
			orderRpcTarget+orderpb.Order_CreateRollback_FullMethodName, // 补偿操作-创建订单回滚
			&order.CreateReq{
				UserId:  req.UserID,
				GoodsId: req.GoodsID,
				Num:     req.Num,
			},
		).
		// 添加库存服务
		Add(
			stockRpcTarget+stockpb.Stock_Deduct_FullMethodName,         // 扣减库存
			stockRpcTarget+stockpb.Stock_DeductRollback_FullMethodName, // 补偿操作-扣减库存回滚
			&stock.DeductReq{
				GoodsId: req.GoodsID,
				Num:     req.Num,
			},
		)

	// 提交数据到dtm
	err = m.Submit()
	logger.FatalIfError(err)

	return
}
