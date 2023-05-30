package route

import (
	"encoding/json"
	"fmt"

	// "fmt"
	"net/http"

	"context"

	"github.com/gin-gonic/gin"
)

type AppReq struct {
	Body    map[string]interface{}
	method  string
	headers map[string]string
	Query   map[string]string
	Params  map[string]string
}

type AppResp interface{}

type RouteHandler func(ctx context.Context, req *AppReq) AppResp
type AppMiddleWare func(g *gin.Context)

type RouteDef struct {
	Path        string
	Version     string
	Method      string
	Handler     RouteHandler
	Middlewares []gin.HandlerFunc // middleware needed to execute before executing request
}

func (r *RouteDef) GetPath() string {
	// path creation logic
	return r.Version + r.Path + "/"
}

var clientRoutes []RouteDef = []RouteDef{}

func RegisterRoutes(r RouteDef) {
	clientRoutes = append(clientRoutes, r)
	fmt.Println("registering routest", clientRoutes)
}

// route handler laydr
func InitializeRoutes(server *gin.Engine) {
	//common middleware that sits in between framework and service and do transaformatino request set to app and response received from service
	// component
	fmt.Println(clientRoutes, "clined routes")
	for _, route := range clientRoutes {
		r := route
		ginHandlerFunc := func(ctx *gin.Context) {

			// create service request

			appReq := &AppReq{
				Body:    make(map[string]interface{}),
				method:  ctx.Request.Method,
				headers: make(map[string]string),
				Query:   make(map[string]string),
				Params:  make(map[string]string),
			}

			extractData(ctx, appReq)

			// call service
			// fmt.Println(appReq)
			resp := r.Handler(ctx.Request.Context(), appReq)
			json.MarshalIndent(resp, "	", "\n")

			// fmt.Println(resp)
			ctx.JSON(http.StatusOK, resp)
			return

		}
		routeHandlers := append(r.Middlewares, ginHandlerFunc)
		// fmt.Println(r.Method,r.Handler)
		server.Handle(r.Method, r.GetPath(), routeHandlers...)
	}
}

func extractData(ctx *gin.Context, appReq *AppReq) {

	for k, v := range ctx.Request.Header {
		if len(v) > 0 {
			appReq.headers[k] = v[0]
		}
	}
	for k, v := range ctx.Request.URL.Query() {
		if len(v) > 0 {
			appReq.Query[k] = v[0]
		}
	}
	for _, p := range ctx.Params {
		{
			appReq.Params[p.Key] = p.Value
		}
	}
	var jsonInput map[string]interface{}
	if err := ctx.BindJSON((&jsonInput)); err == nil {
		appReq.Body = jsonInput
	}
}
