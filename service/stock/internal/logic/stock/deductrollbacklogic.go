package stocklogic

import (
	"context"
	"database/sql"

	"dtm-zero/service/stock/internal/dal/query"
	"dtm-zero/service/stock/internal/svc"
	"dtm-zero/service/stock/pb"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeductRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductRollbackLogic {
	return &DeductRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeductRollback 扣减库存回滚
// 该函数用于在分布式事务中回滚库存扣减操作。它接收一个扣减请求，并尝试将库存数量恢复到扣减前的状态。
//
// 参数:
//   - in: *pb.DeductReq, 包含需要回滚的库存扣减请求信息，如商品ID和扣减数量。
//
// 返回值:
//   - *pb.DeductResp: 返回一个空的扣减响应，表示回滚操作成功。
//   - error: 如果回滚过程中发生错误，返回相应的错误信息。
func (l *DeductRollbackLogic) DeductRollback(in *pb.DeductReq) (*pb.DeductResp, error) {
	l.Infof("库存回滚 in : %+v ", in)

	// 获取底层数据库连接 gorm.DB
	stockQuery := l.svcCtx.Query.Stock
	underlyingDB := stockQuery.WithContext(l.ctx).UnderlyingDB()
	db, err := underlyingDB.DB()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 从 gRPC 上下文中获取 Barrier 对象，用于分布式事务管理
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		// 如果获取Barrier失败，返回内部错误，但不触发DTM回滚
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 在数据库事务中执行库存回滚操作
	err = barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 更新库存数量，将库存数量增加回扣减前的值
		underlyingDB.ConnPool = tx
		sqlResult, err := query.Use(underlyingDB).Stock.WithContext(l.ctx).Where(stockQuery.GoodsID.Eq(in.GoodsId)).Update(stockQuery.Num, stockQuery.Num.Add(int32(in.Num)))
		if err != nil {
			// 如果更新库存时发生错误，返回内部错误，但不触发DTM回滚
			return status.Error(codes.Internal, err.Error())
		}

		// 如果影响行数为0，返回 Aborted 错误，触发 DTM 回滚
		if sqlResult.RowsAffected <= 0 {
			return status.Error(codes.Aborted, dtmcli.ResultFailure)
		}

		return nil
	})
	if err != nil {
		l.Errorf("err : %v \n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeductResp{}, nil
}
