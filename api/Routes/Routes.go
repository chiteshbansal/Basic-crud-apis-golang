package route

import (
	"encoding/json"
	"fmt"

	// "fmt"

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

type AppResp map[string]interface{}

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
			resp := r.Handler(ctx.Request.Context(), appReq)
			json.MarshalIndent(resp, "	", "\n")

			// fmt.Println(resp)
			ctx.JSON(resp["status"].(int), resp)
			return

		}
		routeHandlers := append(r.Middlewares, ginHandlerFunc)
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
	body, exists := ctx.Get("body")

	if !exists {
		var jsonInput map[string]interface{}
		if err := ctx.BindJSON((&jsonInput)); err == nil {
			appReq.Body = jsonInput
		}
	} else {
		var err error
		appReq.Body, err = StructToMapStringInterface(body)
		fmt.Println(err)
	}
}

func StructToMapStringInterface(s interface{}) (map[string]interface{}, error) {

	marshalledDate, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(marshalledDate, &data)
	if err != nil {
		return nil, err
	}
	return data, nil

}
