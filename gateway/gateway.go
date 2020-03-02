package gateway

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// APIGateway 是web服务器的入口网关。
type APIGateway struct {
	router      *http.ServeMux // 路由
	preFilters  []Filter       // 前置过滤器
	postFilters []Filter       // 后置过滤器
}

// New 会返回一个持有空ServeMux的APIGateway。
func New() *APIGateway {
	gw := new(APIGateway)
	gw.router = http.NewServeMux()
	return gw
}

// ServeHTTP 实现http.Handler接口。
func (gw *APIGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, filter := range gw.preFilters {
		rw, req, code, err := filter.FilterFunc(w, r)
		if err != nil {
			NewJSONResponse(w, err.Error(), code, nil)
			return
		}
		w = rw
		r = req
	}
	gw.router.ServeHTTP(w, r)
}

// Filter 过滤器，在接收到请求到之前和之后做某些操作。
type Filter struct {
	// Name 过滤器名字
	Name string
	// Type 过滤器类型，有前置和后置两种
	Type FilterType
	// Order 过滤器的执行顺序
	Order uint
	// FilterFunc 过滤器函数
	FilterFunc FilterFunc
}

// FilterType 是过滤器类型，分为前置过滤器和后置过滤器。
// 前置过滤器会在分发请求前执行。
// 后置过滤器会在发送响应前执行。
type FilterType int

const (
	// PRE_FILTER_TYPE 前置过滤器
	PRE_FILTER_TYPE FilterType = iota
	// POST_FILTER_TYPE 后置过滤器
	POST_FILTER_TYPE
)

// FilterFunc 过滤器函数。
// 基本思想是拦截请求和响应，加上自定义逻辑，并把额外信息通过context传递。
type FilterFunc func(w http.ResponseWriter, r *http.Request) (rw http.ResponseWriter, req *http.Request, statusCode int, err error)

// AddPreFilter 添加前置过滤器并按照执行顺序排序
func (gw *APIGateway) AddPreFilter(filter *Filter) {
	if filter.Type != PRE_FILTER_TYPE {
		log.Fatal("filter类型错误")
	}
	gw.preFilters = append(gw.preFilters, *filter)
	if len(gw.preFilters) > 1 {
		var last = len(gw.preFilters)
		for last > 1 {
			for i := 1; i < last; i++ {
				if gw.preFilters[i-1].Order > gw.preFilters[i].Order {
					gw.preFilters[i-1], gw.preFilters[i] = gw.preFilters[i], gw.preFilters[i-1]
				}
			}
			last--
		}

	}
}

type Route struct {
	URL     string
	Handler http.HandlerFunc
	Method  string
}

// AddRoute 添加路由
func (gw *APIGateway) AddRoute(route Route) {
	switch strings.ToUpper(route.Method) {
	case "POST":
		route.Handler = methodPost(route.Handler)
	case "GET":
		route.Handler = methodGet(route.Handler)
	default:
		log.Fatalf("不支持的method: %v", route.Method)
	}
	gw.router.HandleFunc(route.URL, route.Handler)
}

type HTTPServer = http.Server

func (gw *APIGateway) Start(server *HTTPServer) {
	log.Fatal(server.ListenAndServe())
}

func methodGet(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "use GET method", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}

func methodPost(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "use POST method", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}

// JSONResponse 是返回json类型的结构
type JSONResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Body       interface{} `json:"body"`
}

// NewJSONResponse 创建一个JSONResponse并用http.ResponseWriter返回
func NewJSONResponse(w http.ResponseWriter, message string, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(JSONResponse{
		StatusCode: statusCode,
		Message:    message,
		Body:       body,
	})
	if err != nil {
		log.Print(err)
	}
}
