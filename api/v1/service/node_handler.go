package service

import (
	"net/http"

	"github.com/kobeHub/Pegasus-engine/pkg/common/k8s"
)

func GetNodes(w http.ResponseWriter, r *http.Request) {
	pods, err := k8s.ListRescheduleLabelPods("reschedulable=false")
	if err != nil {
		respondError(w, apiError{errorInternal, err}, "")
		return
	}

	if res, err := k8s.ListReschedulableNodes(pods); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, res)
	}
}

func GetNS(w http.ResponseWriter, r *http.Request) {
	if res, err := k8s.GetWorkNS(); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, res)
	}
}

func GetReschedulablePods(w http.ResponseWriter, r *http.Request) {
	if res, err := k8s.ListRescheduleLabelPods("reschedulable=true"); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, res)
	}
}
