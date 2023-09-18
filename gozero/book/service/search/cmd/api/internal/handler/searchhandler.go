package handler

import (
	"net/http"

	"StandardProject/gozero/book/service/search/cmd/api/internal/logic"
	"StandardProject/gozero/book/service/search/cmd/api/internal/svc"
	"StandardProject/gozero/book/service/search/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func searchHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewSearchLogic(r.Context(), ctx)
		resp, err := l.Search(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
