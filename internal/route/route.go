// The route package provides functionalities for routing and handling HTTP requests.
package route

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// AppReq represents the structure of an application request.
type AppReq struct {
	Body    map[string]interface{}
	Method  string
	Headers map[string]string
	Query   map[string]string
	Params  map[string]string
}

// AppResp represents the structure of an application response.
type AppResp map[string]interface{}

// RouteHandler defines the type for a route handling function.
type RouteHandler func(ctx context.Context, req *AppReq) AppResp

// AppMiddleWare defines the type for a middleware function.
type AppMiddleWare func(g *gin.Context)

// RouteDef defines the structure of a route definition.
type RouteDef struct {
	Path        string
	Version     string
	Method      string
	Handler     RouteHandler
	Middlewares []gin.HandlerFunc // middleware needed to execute before executing request
}

// GetPath constructs and returns the complete path for the route definition.
func (r *RouteDef) GetPath() string {
	// path creation logic
	return r.Version + r.Path
}

var clientroutes []RouteDef = []RouteDef{}

// RegisterRoutes adds a route definition to the list of client routes.
func RegisterRoutes(r RouteDef) {
	clientroutes = append(clientroutes, r)
}

// InitializeRoutes initializes the routes on the given gin engine.
func InitializeRoutes(server *gin.Engine) {
	//common middleware that sits in between framework and service and do transformation request set to app and response received from service
	// component
	for _, route := range clientroutes {
		r := route
		ginHandlerFunc := func(ctx *gin.Context) {

			// create service request
			appReq := &AppReq{
				Body:    make(map[string]interface{}),
				Method:  ctx.Request.Method,
				Headers: make(map[string]string),
				Query:   make(map[string]string),
				Params:  make(map[string]string),
			}

			extractData(ctx, appReq)

			// call service
			resp := r.Handler(ctx.Request.Context(), appReq)
			json.MarshalIndent(resp, "	", "  ")
			if resp["token"] != nil {
				ctx.Writer.Header().Set("Authorization", "Bearer "+resp["token"].(string))
			}
			ctx.JSON(resp["status"].(int), resp)
			return
		}
		routeHandlers := append(r.Middlewares, ginHandlerFunc)
		server.Handle(r.Method, r.GetPath(), routeHandlers...)
	}
}

// extractData extracts data from the gin context and adds it to the application request.
func extractData(ctx *gin.Context, appReq *AppReq) {

	for k, v := range ctx.Request.Header {
		if len(v) > 0 {
			appReq.Headers[k] = v[0]
		}
	}
	for k, v := range ctx.Request.URL.Query() {
		if len(v) > 0 {
			appReq.Query[k] = v[0]
		}
	}
	for _, p := range ctx.Params {
		appReq.Params[p.Key] = p.Value
	}
	body, exists := ctx.Get("body")
	if !exists {
		var jsonInput map[string]interface{}
		if err := ctx.BindJSON(&jsonInput); err == nil {
			appReq.Body = jsonInput
		}
		for k, v := range ctx.Keys {
			appReq.Body[k] = v
		}
	} else {
		appReq.Body, _ = StructToMapStringInterface(body)
		appReq.Body["confirmPassword"], _ = ctx.Get("confirmPassword")
		for k, v := range ctx.Keys {
			appReq.Body[k] = v
		}
	}
}

// StructToMapStringInterface converts a struct to a map[string]interface{}.
func StructToMapStringInterface(s interface{}) (map[string]interface{}, error) {
	marshalledData, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(marshalledData, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
