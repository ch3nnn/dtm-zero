package handler

import (
	"net/http"

	"dtm-zero/restful/internal/logic"
	"dtm-zero/restful/internal/svc"
	"dtm-zero/restful/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建订单
func createHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCreateLogic(r.Context(), svcCtx)
		resp, err := l.Create(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
