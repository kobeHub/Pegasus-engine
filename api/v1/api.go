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
	r.Get("/ns", wrap(service.GetNS))
	r.Get("/pods", wrap(service.GetReschedulablePods))

	repoRouter := r.WithPrefix("/repo")
	repoRouter.Post("/getRepo", wrap(service.GetRepo))
	repoRouter.Post("/getRepoList", wrap(service.GetRepoList))
	repoRouter.Post("/getRepoTags", wrap(service.GetRepoTags))
	repoRouter.Post("/createRepo", wrap(service.CreateRepo))
	repoRouter.Post("/createRepoRule", wrap(service.CreateRepoBuildRule))
	repoRouter.Post("/getRepoBuildRule", wrap(service.GetRepoBuildRule))
	repoRouter.Post("/startRepoBuild", wrap(service.StartRepoBuildByRule))
	repoRouter.Post("/deleteRepoBuildRule", wrap(service.DeleteRepoBuildRule))
	repoRouter.Post("/deleteImage", wrap(service.DeleteImage))
	repoRouter.Post("/deleteRepo", wrap(service.DeleteRepo))
}
