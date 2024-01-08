package statuses

import (
	"fmt"
	"strings"
)

type KubeApiHealthEndpointStatus struct {
	Name   string
	Path   string
	Status string
}

type KubeApiHealthEndpointStatusList struct {
	Items []KubeApiHealthEndpointStatus
}

func getApiStatusName(line string) string {
	trimmedLine := line[3:]
	words := strings.Fields(trimmedLine)
	wordCount := len(words)
	return strings.Join(words[:wordCount-1], ",")
}

func (l *KubeApiHealthEndpointStatusList) Add(s KubeApiHealthEndpointStatus) {
	l.Items = append(l.Items, s)
}

func (l *KubeApiHealthEndpointStatusList) AddRawLine(line string, path string) {
	if len(path) == 0 || len(line) == 0 || line == fmt.Sprintf("%s check passed", path[1:]) {
		return
	}

	var status string
	var name string

	switch line[0:3] {
	case "[+]":
		status = "ok"
		name = getApiStatusName(line)
	case "[-]":
		status = "bad"
		name = getApiStatusName(line)
	default:
		status = "unknown"
		name = line
	}

	l.Add(KubeApiHealthEndpointStatus{
		Name:   name,
		Path:   path,
		Status: status,
	})
}
