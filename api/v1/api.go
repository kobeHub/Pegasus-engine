/*Pegasus-engine api version 1.0*/

package v1

import (
	"net/http"
	"time"

	"github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"

	"github.com/kobeHub/Pegasus-engine/pkg/common/router"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
