package service

import "net/http"

func Hello(w http.ResponseWriter, r *http.Request) {
	respond(w, "Welcome use pegasus-engine")
}

func ApiHealthy(w http.ResponseWriter, r *http.Request) {
	respond(w, "The pegasus-engine api is healthy")
}
