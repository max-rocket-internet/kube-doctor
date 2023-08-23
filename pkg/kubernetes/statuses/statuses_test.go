package statuses

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckServices(t *testing.T) {
	r := KubeApiHealthEndpointStatusList{}
	r.AddRawLine("[+]etcd-readiness ok", "/readyz")

	assert.Len(t, r.Items, 1)
	assert.Equal(t, "ok", r.Items[0].Status)
	assert.Equal(t, "etcd-readiness", r.Items[0].Name)
	assert.Equal(t, "/readyz", r.Items[0].Path)

	r.AddRawLine("[-]poststarthook/apiservice-status-available-controller bad", "/readyz")

	assert.Len(t, r.Items, 2)
	assert.Equal(t, "bad", r.Items[1].Status)
	assert.Equal(t, "poststarthook/apiservice-status-available-controller", r.Items[1].Name)
	assert.Equal(t, "/readyz", r.Items[1].Path)

	r.AddRawLine("[]asasdasd", "")
	r.AddRawLine("", "/path")

	assert.Len(t, r.Items, 2)

	r.AddRawLine("some new unexpected check output here", "/readyz")

	assert.Len(t, r.Items, 3)
	assert.Equal(t, "unknown", r.Items[2].Status)
	assert.Equal(t, "some new unexpected check output here", r.Items[2].Name)
	assert.Equal(t, "/readyz", r.Items[2].Path)
}
