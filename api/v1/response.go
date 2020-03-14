package v1

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type response struct {
	Status    status      `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	ErrorType errorType   `json:"errorType,omitempty"`
	Error     string      `json:"error,omitempty"`
}

func respond(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	b, err := json.Marshal(&response{
		Status: statusSuccess,
		Data:   data,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Error marshaling JSON")
		return
	}

	if _, err := w.Write(b); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Faild to write data to connection")
	}
}

func respondError(w http.ResponseWriter, apiErr apiError, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	switch apiErr.typ {
	case errorBadData:
		w.WriteHeader(http.StatusBadRequest)
	case errorInternal:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		log.Fatal(fmt.Sprintf("unknown error type %q", apiErr.Error()))
	}

	b, err := json.Marshal(&response{
		Status:    statusError,
		ErrorType: apiErr.typ,
		Error:     apiErr.err.Error(),
		Data:      data,
	})
	if err != nil {
		return
	}
	log.WithFields(log.Fields{
		"err": apiErr.Error(),
	}).Error("Api error")

	if _, err := w.Write(b); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Faild to write data to connection")
	}
}
