// Package statuses parses the verbose output of the Kubernetes API health endpoints.
package statuses

import "strings"

const (
	StatusOK      = "ok"
	StatusBad     = "bad"
	StatusUnknown = "unknown"
)

// KubeAPIHealthEndpointStatus is a single check reported by /readyz or /livez.
type KubeAPIHealthEndpointStatus struct {
	Name   string
	Path   string
	Status string
}

// KubeAPIHealthEndpointStatusList is a collection of health endpoint checks.
type KubeAPIHealthEndpointStatusList struct {
	Items []KubeAPIHealthEndpointStatus
}

// Add appends s to the list.
func (l *KubeAPIHealthEndpointStatusList) Add(s KubeAPIHealthEndpointStatus) {
	l.Items = append(l.Items, s)
}

// AddRawLine parses one line of verbose health endpoint output, e.g.
// "[+]etcd ok" or "[-]etcd failed: reason withheld", and adds it to the list.
// Blank lines and the trailing "<path> check passed" summary are ignored.
func (l *KubeAPIHealthEndpointStatusList) AddRawLine(line, path string) {
	if path == "" || line == "" || line == strings.TrimPrefix(path, "/")+" check passed" {
		return
	}

	status, name := StatusUnknown, line
	if rest, ok := strings.CutPrefix(line, "[+]"); ok {
		status, name = StatusOK, checkName(rest)
	} else if rest, ok := strings.CutPrefix(line, "[-]"); ok {
		status, name = StatusBad, checkName(rest)
	}

	l.Add(KubeAPIHealthEndpointStatus{
		Name:   name,
		Path:   path,
		Status: status,
	})
}

func checkName(s string) string {
	name, _, _ := strings.Cut(strings.TrimSpace(s), " ")
	return name
}
