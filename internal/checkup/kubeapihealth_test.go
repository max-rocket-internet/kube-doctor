package checkup

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/kubernetes/statuses"
)

func TestKubeApiHealthStatusesWithError(t *testing.T) {
	dummyResources := statuses.KubeAPIHealthEndpointStatusList{}
	dummyResources.AddRawLine("[+]etcd-readiness ok", "/readyz")
	dummyResources.AddRawLine("[-]poststarthook/apiservice-status-available-controller error", "/readyz")

	result := CheckKubeAPIHealth(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "Kubernetes API health endpoint '/readyz' poststarthook/apiservice-status-available-controller=bad", result.Symptoms[0].Message)
	assert.Equal(t, symptoms.Critical, result.Symptoms[0].Severity)
}

func TestKubeApiHealthStatusesAllOK(t *testing.T) {
	dummyResources := statuses.KubeAPIHealthEndpointStatusList{}
	dummyResources.AddRawLine("[+]etcd-readiness ok", "/readyz")
	dummyResources.AddRawLine("[+]poststarthook/apiservice-status-available-controller ok", "/readyz")

	result := CheckKubeAPIHealth(&dummyResources)

	assert.Len(t, result.Symptoms, 0)
}
