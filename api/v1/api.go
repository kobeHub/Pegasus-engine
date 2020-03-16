// +build jsoniter
/*Pegasus-engine api version 1.0*/

package v1

import (
	"net/http"
	"sync"
	_ "time"

	_ "github.com/sirupsen/logrus"

	"github.com/kobeHub/Pegasus-engine/api/v1/service"
	"github.com/kobeHub/Pegasus-engine/pkg/common/router"
)

type Api struct {
	mtx sync.RWMutex
}

var corsHeaders = map[string]string{
	"Access-Control-Allow-Headers":  "Accept, Authorization, Content-Type, Origin",
	"Access-Control-Allow-Methods":  "GET, POST, DELETE, OPTIONS",
	"Access-Control-Allow-Origin":   "*",
	"Access-Control-Expose-Headers": "Date",
	"Cache-Control":                 "no-cache, no-store, must-revalidate",
}

// Enable cross-site script calls
func setCORS(w http.ResponseWriter) {
	for h, v := range corsHeaders {
		w.Header().Set(h, v)
	}
}

func Register(r *router.Router) {
	wrap := func(f http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			setCORS(w)
			f(w, r)
		})
	}

	r.Get("/", wrap(service.Hello))

	r = r.WithPrefix("/api/v1")
	r.Get("/", wrap(service.ApiHealthy))
	r.Get("/nodes", wrap(service.GetNodes))
}
