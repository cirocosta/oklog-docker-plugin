package http

import (
	"encoding/json"
	"net/http"

	"github.com/cirocosta/oklog-docker-plugin/docker"
	"github.com/cirocosta/oklog-docker-plugin/driver"
	"github.com/docker/go-plugins-helpers/sdk"
)

func Handlers(h *sdk.Handler, d *driver.Driver) {
	h.HandleFunc("/LogDriver.StartLogging", func(w http.ResponseWriter, r *http.Request) {
		var (
			req startLoggingRequest
			err error
		)

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = d.StartLogging(req.File, req.Info)
		respond(err, w)
	})

	h.HandleFunc("/LogDriver.StopLogging", func(w http.ResponseWriter, r *http.Request) {
		var (
			req stopLoggingRequest
			err error
		)

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = d.StopLogging(req.File)
		respond(err, w)
	})

	h.HandleFunc("/LogDriver.Capabilities", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(&capabilitiesResponse{
			Cap: docker.Capability{ReadLogs: false},
		})
	})
}

func respond(err error, w http.ResponseWriter) {
	var res response

	if err != nil {
		res.Err = err.Error()
	}

	json.NewEncoder(w).Encode(&res)
}
