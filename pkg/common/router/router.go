package router

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type param string

// Param return context params or empty string
// is nothing exists
func Param(ctx context.Context, k string) string {
	if v := ctx.Value(param(k)); v != nil {
		return v.(string)
	}
	return ""
}

// Return a Context with param k set to v
func WithParam(ctx context.Context, k, v string) context.Context {
	return context.WithValue(ctx, param(p), v)
}

// Wrap httprouter and support prefix group
// also could spefic `HandlerFunc`
type Router struct {
	r      *httprouter.Router
	prefix string
	instrh func(handlerName string, handler http.HandlerFunc) http.HandlerFunc
}

func New() *Router {
	return &Router{
		r: httprouter.New(),
	}
}

// Return router with default handler func
func (info *Router) WithInstrumentation(instrh func(handlerName string, handler http.HandlerFunc) http.HandlerFunc) *Router {
	if r.instrh != nil {
		newInstrh := instrh
		instrh = func(handlerName string, handler http.HandlerFunc) http.HandlerFunc {
			return newInstrh(handlerName, info.instrh(handlerName, handler))
		}
	}
	return &Router{r: info.r, prefix: info.prefix, instrh: instrh}

}

// Return a group router with same prefix
func (info *Router) WithPrefix(prefix string) *Router {
	return &Router{r: info.r, prefix: r.prefix + prefix, instrh: info.instrh}
}

// Turn a HandlerFunc to httprouter.Handle
func (info *Router) handle(handleName string, h http.HandlerFunc) httprouter.Handle {
	if info.instrh != nil {
		h = r.instrh(handleName, h)
	}
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()

		for _, p := range params {
			ctx = context.WithValue(ctx, param(p.Key), p.Value)
		}
		h(w, req.WithContext(ctx))
	}
}

// Get registers a new GET route.
func (info *Router) Get(path string, h http.HandlerFunc) {
	info.r.GET(info.prefix+path, info.handle(path, h))
}

// Options registers a new OPTIONS route.
func (info *Router) Options(path string, h http.HandlerFunc) {
	info.r.OPTIONS(info.prefix+path, info.handle(path, h))
}

// Del registers a new DELETE route.
func (info *Router) Del(path string, h http.HandlerFunc) {
	info.r.DELETE(info.prefix+path, info.handle(path, h))
}

// Put registers a new PUT route.
func (info *Router) Put(path string, h http.HandlerFunc) {
	info.r.PUT(info.prefix+path, info.handle(path, h))
}

// Post registers a new POST route.
func (info *Router) Post(path string, h http.HandlerFunc) {
	info.r.POST(info.prefix+path, info.handle(path, h))
}

// Redirect takes an absolute path and sends an internal HTTP redirect for it,
// prefixed by the router's path prefix. Note that this method does not include
// functionality for handling relative paths or full URL redirects.
func (info *Router) Redirect(w http.ResponseWriter, req *http.Request, path string, code int) {
	http.Redirect(w, req, info.prefix+path, code)
}

// ServeHTTP implements http.Handler.
func (info *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	info.r.ServeHTTP(w, req)
}

// FileServe returns a new http.HandlerFunc that serves files from dir.
// Using routes must provide the *filepath parameter.
func FileServe(dir string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))

	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = Param(r.Context(), "filepath")
		fs.ServeHTTP(w, r)
	}
}
