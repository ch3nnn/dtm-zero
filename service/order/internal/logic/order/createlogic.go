package orderlogic

import (
	"context"
	"database/sql"
	"fmt"

	"dtm-zero/service/order/internal/dal/model"
	"dtm-zero/service/order/internal/dal/query"
	"dtm-zero/service/order/internal/svc"
	"dtm-zero/service/order/pb"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Create 创建订单
func (l *CreateLogic) Create(in *pb.CreateReq) (*pb.CreateResp, error) {
	l.Infof("创建订单 in : %+v", in)

	// 从gRPC上下文中获取Barrier对象，用于防止空补偿、空悬挂等问题
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		// 如果获取Barrier失败，返回内部错误，但不触发DTM回滚
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 获取底层数据库连接 gorm.DB
	underlyingDB := l.svcCtx.Query.Order.WithContext(l.ctx).UnderlyingDB()
	db, err := underlyingDB.DB()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 子事务屏障中执行订单创建操作 (空补偿、悬挂、重复)
	err = barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 创建订单记录
		underlyingDB.ConnPool = tx
		err := query.Use(underlyingDB).Order.WithContext(l.ctx).Create(&model.Order{
			UserID:  in.UserId,
			GoodsID: in.GoodsId,
			Num:     int32(in.Num),
		})
		if err != nil {
			// 如果创建订单失败，返回错误信息
			return fmt.Errorf("创建订单失败 err : %v ", err)
		}

		return nil
	})
	if err != nil {
		// 如果事务执行失败，返回内部错误，但不触发DTM回滚
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 返回空创建订单响应
	return &pb.CreateResp{}, nil
}
