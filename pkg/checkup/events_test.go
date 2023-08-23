package checkup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEventsClusterAutoscalerScaling(t *testing.T) {
	dummyResources := v1.EventList{
		Items: []v1.Event{
			{
				Reason:        "TriggeredScaleUp",
				LastTimestamp: metav1.NewTime(time.Now()),
				Message:       "a message",
				InvolvedObject: v1.ObjectReference{
					Kind: "Pod",
				},
				Type: "Warning",
				Source: v1.EventSource{
					Component: "cluster-autoscaler",
				},
			},
		},
	}

	result := CheckEvents(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "(Pod) 0.0 minutes ago: a message", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestEventsServiceControllerNotNormal(t *testing.T) {
	dummyResources := v1.EventList{
		Items: []v1.Event{
			{
				Reason:        "SomeReason",
				LastTimestamp: metav1.NewTime(time.Now()),
				Message:       "a message",
				InvolvedObject: v1.ObjectReference{
					Kind: "Pod",
				},
				Type: "Warning",
				Source: v1.EventSource{
					Component: "service-controller",
				},
			},
		},
	}

	result := CheckEvents(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "(Pod) 0.0 minutes ago: a message", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestEventsDefaultSchedulerNotNormal(t *testing.T) {
	dummyResources := v1.EventList{
		Items: []v1.Event{
			{
				Reason:        "SomeReason",
				LastTimestamp: metav1.NewTime(time.Now()),
				Message:       "a message",
				InvolvedObject: v1.ObjectReference{
					Kind: "Pod",
				},
				Type: "Warning",
				Source: v1.EventSource{
					Component: "default-scheduler",
				},
			},
		},
	}

	result := CheckEvents(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "(Pod) 0.0 minutes ago: a message", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestEventsKubeletNotNormal(t *testing.T) {
	dummyResources := v1.EventList{
		Items: []v1.Event{
			{
				Reason:        "SomeReason",
				LastTimestamp: metav1.NewTime(time.Now()),
				Message:       "a message",
				InvolvedObject: v1.ObjectReference{
					Kind: "Pod",
				},
				Type: "Warning",
				Source: v1.EventSource{
					Component: "kubelet",
				},
			},
		},
	}

	result := CheckEvents(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "(Pod) 0.0 minutes ago: a message", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestEventsNormalToBeIgnored(t *testing.T) {
	dummyResources := v1.EventList{
		Items: []v1.Event{
			{
				Reason:        "SomeReason",
				LastTimestamp: metav1.NewTime(time.Now()),
				Message:       "a message",
				InvolvedObject: v1.ObjectReference{
					Kind: "Pod",
				},
				Type: "Normal",
				Source: v1.EventSource{
					Component: "cluster-autoscaler",
				},
			},
			{
				Reason:        "SomeReason",
				LastTimestamp: metav1.NewTime(time.Now()),
				Message:       "a message",
				InvolvedObject: v1.ObjectReference{
					Kind: "Pod",
				},
				Type: "Normal",
				Source: v1.EventSource{
					Component: "service-controller",
				},
			},
			{
				Reason:        "SomeReason",
				LastTimestamp: metav1.NewTime(time.Now()),
				Message:       "a message",
				InvolvedObject: v1.ObjectReference{
					Kind: "Pod",
				},
				Type: "Normal",
				Source: v1.EventSource{
					Component: "default-scheduler",
				},
			},
			{
				Reason:        "SomeReason",
				LastTimestamp: metav1.NewTime(time.Now()),
				Message:       "a message",
				InvolvedObject: v1.ObjectReference{
					Kind: "Pod",
				},
				Type: "Normal",
				Source: v1.EventSource{
					Component: "kubelet",
				},
			},
		},
	}

	result := CheckEvents(&dummyResources)

	assert.Len(t, result.Symptoms, 0)
}
