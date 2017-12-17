package docker

import (
	"strings"
	"time"
)

type Capability struct {
	ReadLogs bool
}

type Info struct {
	Config              map[string]string
	ContainerID         string
	ContainerName       string
	ContainerEntrypoint string
	ContainerArgs       []string
	ContainerImageID    string
	ContainerImageName  string
	ContainerCreated    time.Time
	ContainerEnv        []string
	ContainerLabels     map[string]string
	LogPath             string
	DaemonName          string
}

// GetOptLabel extracts a list of labels from `--log-opt`
// configuration and populates a map with their keys and
// values.
func (info *Info) GetOptLabels() (optLabels map[string]string) {
	optLabels = make(map[string]string)

	labels, ok := info.Config["labels"]
	if ok && len(labels) > 0 {
		for _, l := range strings.Split(labels, ",") {
			if v, ok := info.ContainerLabels[l]; ok {
				optLabels[l] = v
			}
		}
	}

	return
}
