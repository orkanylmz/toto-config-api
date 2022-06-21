package ports

import (
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"toto-config-api/internal/common/server/httperr"
	"toto-config-api/internal/skuconfig/app"
	"toto-config-api/internal/skuconfig/app/query"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app}
}

func (h HttpServer) GetSKU(w http.ResponseWriter, r *http.Request, params GetSKUParams) {

	countryCode := r.Header.Get("X-Appengine-Country")

	if countryCode == "" {
		countryCode = "ZZ"
	}

	sku, err := h.app.Queries.SKUForConfig.Handle(r.Context(), query.SKUForConfig{
		PackageName: params.Package,
		CountryCode: countryCode,
	})

	fmt.Println("AFTER QUERY: ", sku, err)

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	res := skuToResponse(sku)

	render.Respond(w, r, res)
}

func skuToResponse(sku string) SKUResponse {
	return SKUResponse{
		MainSku: sku,
	}
}
