/*Pegasus-engine api version 1.0*/

package v1

import (
	"net/http"
	"sync"
	_ "time"

	"github.com/json-iterator/go"
	_ "github.com/sirupsen/logrus"

	"github.com/kobeHub/Pegasus-engine/pkg/common/router"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Api struct {
	mtx sync.RWMutex
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

	r.Get("/", wrap(func(w http.ResponseWriter, r *http.Request) {
		respond(w, "Pegasus-engine is healty!")
	}))
}
