package service

import (
	"io/ioutil"
	"net/http"

	"github.com/kobeHub/Pegasus-engine/internal/registry"
)

// Basic imge struct
type repoImage struct {
	Name      string  `json:"name,required"`
}

func GetRepo(w http.ResponseWriter, r *http.Request) {
	var info repoImage
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}

	if err := json.Unmarshal(jsonData, &info); err != nil {
		respondError(w, apiError{errorBadData, err}, "")
		return
	}
	if err := registry.GetRepo(info.Name); err != nil {
		respondError(w, apiError{errorInternal, err}, "")
	} else {
		respond(w, "")
	}
}
