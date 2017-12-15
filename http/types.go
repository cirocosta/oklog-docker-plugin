package http

import (
	"github.com/cirocosta/oklog-docker-plugin/docker"
)

type startLoggingRequest struct {
	File string
	Info docker.Info
}

type stopLoggingRequest struct {
	File string
}

type capabilitiesResponse struct {
	Err string
	Cap docker.Capability
}

type response struct {
	Err string
}
