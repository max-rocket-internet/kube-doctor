package checkup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCheckNodesNotReady(t *testing.T) {
	dummyResources := v1.NodeList{
		Items: []v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test1",
					CreationTimestamp: metav1.NewTime(time.Now()),
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:   "Ready",
							Status: "False",
						},
					},
					NodeInfo: v1.NodeSystemInfo{
						KubeletVersion: "v1.25.9",
					},
				},
			},
		},
	}

	result := CheckNodes(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "not ready", result.Symptoms[0].Message)
	assert.Equal(t, "warning", result.Symptoms[0].Severity)
}

func TestCheckNodesOldAndNotReady(t *testing.T) {
	dummyResources := v1.NodeList{
		Items: []v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test1",
					CreationTimestamp: metav1.NewTime(time.Now().Add(-10 * time.Minute)),
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:   "Ready",
							Status: "False",
						},
					},
					NodeInfo: v1.NodeSystemInfo{
						KubeletVersion: "v1.25.9",
					},
				},
			},
		},
	}

	result := CheckNodes(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "older than 5 minutes and not Ready", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestCheckNodesBadCondition(t *testing.T) {
	dummyResources := v1.NodeList{
		Items: []v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test1",
					CreationTimestamp: metav1.NewTime(time.Now()),
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:   "KernelHasNoDeadlock",
							Status: "True",
							Reason: "SomeReason",
						},
					},
					NodeInfo: v1.NodeSystemInfo{
						KubeletVersion: "v1.25.9",
					},
				},
			},
		},
	}

	result := CheckNodes(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "has condition KernelHasNoDeadlock=True: SomeReason", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestCheckNodesMultipleVersions(t *testing.T) {
	dummyResources := v1.NodeList{
		Items: []v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test1",
					CreationTimestamp: metav1.NewTime(time.Now()),
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:   "Ready",
							Status: "True",
						},
					},
					NodeInfo: v1.NodeSystemInfo{
						KubeletVersion: "v1.25.9",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test1",
					CreationTimestamp: metav1.NewTime(time.Now()),
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:   "Ready",
							Status: "True",
						},
					},
					NodeInfo: v1.NodeSystemInfo{
						KubeletVersion: "v1.24.7",
					},
				},
			},
		},
	}

	result := CheckNodes(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "multiple node versions running", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}
