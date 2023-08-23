package checkup

import (
	"testing"

	"github.com/max-rocket-internet/kube-doctor/pkg/kubernetes/statuses"
	"github.com/stretchr/testify/assert"
)

func TestKubeApiHealthStatusesWithError(t *testing.T) {
	dummyResources := statuses.KubeApiHealthEndpointStatusList{}
	dummyResources.AddRawLine("[+]etcd-readiness ok", "/readyz")
	dummyResources.AddRawLine("[-]poststarthook/apiservice-status-available-controller error", "/readyz")

	result := KubeApiHealthStatuses(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "Kubnernetes API health endpoint '/readyz' poststarthook/apiservice-status-available-controller=bad", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestKubeApiHealthStatusesAllOK(t *testing.T) {
	dummyResources := statuses.KubeApiHealthEndpointStatusList{}
	dummyResources.AddRawLine("[+]etcd-readiness ok", "/readyz")
	dummyResources.AddRawLine("[+]poststarthook/apiservice-status-available-controller ok", "/readyz")

	result := KubeApiHealthStatuses(&dummyResources)

	assert.Len(t, result.Symptoms, 0)
}
