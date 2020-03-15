package service

import (
	"net/http"

	"github.com/kobeHub/Pegasus-engine/pkg/common/k8s"
)

func GetNodes(w http.ResponseWriter, r *http.Request) {
	if res, err := k8s.ListNodes(); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, res)
	}
}
