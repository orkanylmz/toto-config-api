package ports

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"toto-config-api/internal/common/middleware"
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
	fmt.Println("REAL IP: ", r.RemoteAddr)
	countryCode := "ZZ"

	if r.RemoteAddr != "" {
		countryCode = GetCountryCodeFromIP(r.RemoteAddr)
	}

	useCache := r.Context().Value(middleware.UseCacheKey)

	fmt.Printf("%T", useCache)

	sku, err := h.app.Queries.SKUForConfig.Handle(r.Context(), query.SKUForConfig{
		PackageName: params.Package,
		CountryCode: countryCode,
		UseCache:    useCache.(bool),
	})

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

type GetIPResponse struct {
	CountryCode string `json:"countryCode"`
}

func GetCountryCodeFromIP(ip string) string {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	resp, err := http.Get(url)
	if err != nil {
		return ""
	}

	if resp.StatusCode != 200 {
		return ""
	}

	var res GetIPResponse

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return ""
	}
	return res.CountryCode
}
