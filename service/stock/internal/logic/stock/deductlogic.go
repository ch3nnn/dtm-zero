package stocklogic

import (
	"context"
	"database/sql"
	"errors"

	"dtm-zero/service/stock/internal/dal/query"
	"dtm-zero/service/stock/internal/svc"
	"dtm-zero/service/stock/pb"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type DeductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductLogic {
	return &DeductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Deduct 扣减库存
//
// 该函数用于处理库存扣减逻辑。它首先查询指定商品的库存信息，然后根据请求中的扣减数量进行库存扣减操作。
// 如果库存不足或操作失败，函数会返回相应的错误信息。该函数还使用了DTM（分布式事务管理器）的Barrier机制，
// 以防止空补偿和空悬挂等问题。
//
// 参数:
//   - in: *pb.DeductReq，包含扣减库存的请求信息，如商品ID和扣减数量。
//
// 返回值:
//   - *pb.DeductResp: 扣减库存的响应信息，通常为空。
//   - error: 如果操作成功，返回nil；否则返回相应的错误信息。
func (l *DeductLogic) Deduct(in *pb.DeductReq) (*pb.DeductResp, error) {
	l.Infof("扣库存start....")

	// 查询指定商品的库存信息
	stockQuery := l.svcCtx.Query.Stock
	stock, err := stockQuery.WithContext(l.ctx).Where(stockQuery.GoodsID.Eq(in.GoodsId)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果查询库存时发生错误（非记录未找到错误），返回内部错误，但不触发DTM回滚
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 如果库存不足，返回Aborted错误，触发DTM回滚
	if stock != nil && int64(stock.Num) < in.Num {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	// 从gRPC上下文中获取Barrier对象，用于防止空补偿、空悬挂等问题
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		// 如果获取Barrier失败，返回内部错误，但不触发DTM回滚
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 获取底层数据库连接 gorm.DB
	underlyingDB := stockQuery.WithContext(l.ctx).UnderlyingDB()
	db, err := underlyingDB.DB()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 在事务中执行库存扣减操作
	err = barrier.CallWithDB(db, func(tx *sql.Tx) error {
		underlyingDB.ConnPool = tx
		sqlResult, err := query.Use(underlyingDB).Stock.WithContext(l.ctx).Where(stockQuery.GoodsID.Eq(in.GoodsId)).Update(stockQuery.Num, stockQuery.Num.Sub(int32(in.Num)))
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
		// 如果事务执行失败，返回错误，但不触发DTM回滚
		return nil, err
	}

	return &pb.DeductResp{}, nil
}
