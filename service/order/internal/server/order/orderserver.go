// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.1
// Source: order.proto

package server

import (
	"context"

	"dtm-zero/service/order/internal/logic/order"
	"dtm-zero/service/order/internal/svc"
	"dtm-zero/service/order/pb"
)

type OrderServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedOrderServer
}

func NewOrderServer(svcCtx *svc.ServiceContext) *OrderServer {
	return &OrderServer{
		svcCtx: svcCtx,
	}
}

// 创建订单
func (s *OrderServer) Create(ctx context.Context, in *pb.CreateReq) (*pb.CreateResp, error) {
	l := orderlogic.NewCreateLogic(ctx, s.svcCtx)
	return l.Create(in)
}

// 创建订单回滚
func (s *OrderServer) CreateRollback(ctx context.Context, in *pb.CreateReq) (*pb.CreateResp, error) {
	l := orderlogic.NewCreateRollbackLogic(ctx, s.svcCtx)
	return l.CreateRollback(in)
}
