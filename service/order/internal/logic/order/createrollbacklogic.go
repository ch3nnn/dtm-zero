package orderlogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"dtm-zero/service/order/internal/dal/query"
	"dtm-zero/service/order/internal/svc"
	"dtm-zero/service/order/pb"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type CreateRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRollbackLogic {
	return &CreateRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateRollback 创建订单回滚
func (l *CreateRollbackLogic) CreateRollback(in *pb.CreateReq) (*pb.CreateResp, error) {
	l.Infof("订单回滚  , in: %+v ", in)

	// 查询订单信息
	var orderQuery = l.svcCtx.Query.Order
	order, err := orderQuery.WithContext(l.ctx).Where(orderQuery.UserID.Eq(in.UserId), orderQuery.GoodsID.Eq(in.GoodsId)).Last()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果查询订单时发生错误（非记录未找到错误），返回内部错误 但不触发DTM回滚
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 从gRPC上下文中获取DTM的Barrier
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		// 如果获取Barrier时发生错误，返回内部错误 但不触发DTM回滚
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 获取底层数据库连接  gorm.DB
	underlyingDB := orderQuery.WithContext(l.ctx).UnderlyingDB()
	db, err := underlyingDB.DB()
	if err != nil {
		// 如果获取数据库连接时发生错误，返回内部错误 但不触发DTM回滚
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 使用Barrier机制执行数据库事务
	err = barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 将订单状态标记为回滚
		underlyingDB.ConnPool = tx
		if _, err = query.Use(underlyingDB).Order.WithContext(l.ctx).Where(orderQuery.ID.Eq(order.ID)).Update(orderQuery.RowState, -1); err != nil {
			// 如果更新订单状态时发生错误，返回错误信息
			return fmt.Errorf("回滚订单失败  err : %v , userId:%d , goodsId:%d", err, in.UserId, in.GoodsId)
		}

		return nil
	})
	if err != nil {
		l.Errorf("err : %v \n", err)
		// 如果事务执行过程中发生错误，返回内部错误 但不触发DTM回滚
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 回滚成功，返回空响应
	return &pb.CreateResp{}, nil
}
