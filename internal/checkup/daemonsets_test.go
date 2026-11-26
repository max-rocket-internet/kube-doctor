package checkup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
)

func TestDaemonSetPodsNotReady(t *testing.T) {
	dummyResources := v1.DaemonSetList{
		Items: []v1.DaemonSet{
			{
				Status: v1.DaemonSetStatus{
					DesiredNumberScheduled: 4,
					UpdatedNumberScheduled: 4,
					NumberUnavailable:      1,
				},
			},
		},
	}

	result := CheckDaemonSets(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "1/4 pods are not ready", result.Symptoms[0].Message)
	assert.Equal(t, symptoms.Critical, result.Symptoms[0].Severity)
}

func TestDaemonSetMisscheduled(t *testing.T) {
	dummyResources := v1.DaemonSetList{
		Items: []v1.DaemonSet{
			{
				Status: v1.DaemonSetStatus{
					UpdatedNumberScheduled: 4,
					DesiredNumberScheduled: 4,
					NumberMisscheduled:     1,
				},
			},
		},
	}

	result := CheckDaemonSets(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "1/4 pods are misscheduled", result.Symptoms[0].Message)
	assert.Equal(t, symptoms.Critical, result.Symptoms[0].Severity)
}

func TestDaemonSetRollingUpdateInProgress(t *testing.T) {
	dummyResources := v1.DaemonSetList{
		Items: []v1.DaemonSet{
			{
				Status: v1.DaemonSetStatus{
					DesiredNumberScheduled: 5,
					NumberReady:            4,
				},
			},
		},
	}

	result := CheckDaemonSets(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "rolling update in progress", result.Symptoms[0].Message)
	assert.Equal(t, symptoms.Warning, result.Symptoms[0].Severity)
}
