package checkup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCheckPodsNotRunning(t *testing.T) {
	dummyResources := v1.PodList{
		Items: []v1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Deployment",
						},
					},
				},
				Status: v1.PodStatus{
					StartTime:  createNewTimeStampPtr(time.Now()),
					Phase:      "Pending",
					Conditions: []v1.PodCondition{},
				},
			},
		},
	}

	result := CheckPods(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "not running", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestCheckPodsBadCondition(t *testing.T) {
	dummyResources := v1.PodList{
		Items: []v1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Deployment",
						},
					},
				},
				Status: v1.PodStatus{
					StartTime: createNewTimeStampPtr(time.Now()),
					Phase:     "Running",
					Conditions: []v1.PodCondition{
						{
							Type:   "ContainersReady",
							Status: "False",
						},
					},
				},
			},
		},
	}

	result := CheckPods(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "status condition ContainersReady is False", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestCheckPodsBadContainerStatuses(t *testing.T) {
	dummyResources := v1.PodList{
		Items: []v1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Deployment",
						},
					},
				},
				Status: v1.PodStatus{
					StartTime:  createNewTimeStampPtr(time.Now()),
					Phase:      "Running",
					Conditions: []v1.PodCondition{},
					ContainerStatuses: []v1.ContainerStatus{
						{
							Ready:        false,
							Name:         "c1",
							RestartCount: 0,
						},
					},
				},
			},
		},
	}

	result := CheckPods(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "container 'c1' is not ready but pod started -0.0 mins ago", result.Symptoms[0].Message)
	assert.Equal(t, "warning", result.Symptoms[0].Severity)
}

func TestCheckPodsBadContainerStatusesOld(t *testing.T) {
	dummyResources := v1.PodList{
		Items: []v1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Deployment",
						},
					},
				},
				Status: v1.PodStatus{
					StartTime:  createNewTimeStampPtr(time.Now().Add(-10 * time.Minute)),
					Phase:      "Running",
					Conditions: []v1.PodCondition{},
					ContainerStatuses: []v1.ContainerStatus{
						{
							Ready:        false,
							Name:         "c1",
							RestartCount: 0,
						},
					},
				},
			},
		},
	}

	result := CheckPods(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "container 'c1' is not ready", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestCheckPodsWithRestartsButOld(t *testing.T) {
	dummyResources := v1.PodList{
		Items: []v1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Deployment",
						},
					},
				},
				Status: v1.PodStatus{
					StartTime:  createNewTimeStampPtr(time.Now()),
					Phase:      "Running",
					Conditions: []v1.PodCondition{},
					ContainerStatuses: []v1.ContainerStatus{
						{
							Ready:        true,
							Name:         "c1",
							RestartCount: 1,
							LastTerminationState: v1.ContainerState{
								Terminated: &v1.ContainerStateTerminated{
									FinishedAt: metav1.NewTime(time.Now().Add(-2 * time.Hour)),
									ExitCode:   1,
									Reason:     "",
								},
							},
						},
					},
				},
			},
		},
	}

	result := CheckPods(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "container 'c1' has been restarted 1 times", result.Symptoms[0].Message)
	assert.Equal(t, "warning", result.Symptoms[0].Severity)
}

func TestCheckPodsWithRestarts(t *testing.T) {
	dummyResources := v1.PodList{
		Items: []v1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "Deployment",
						},
					},
				},
				Status: v1.PodStatus{
					StartTime:  createNewTimeStampPtr(time.Now()),
					Phase:      "Running",
					Conditions: []v1.PodCondition{},
					ContainerStatuses: []v1.ContainerStatus{
						{
							Ready:        true,
							Name:         "c1",
							RestartCount: 1,
							LastTerminationState: v1.ContainerState{
								Terminated: &v1.ContainerStateTerminated{
									FinishedAt: metav1.NewTime(time.Now()),
									ExitCode:   1,
									Reason:     "Crashed",
								},
							},
						},
					},
				},
			},
		},
	}

	result := CheckPods(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "container 'c1' was restarted -0.0 mins ago: 1 (exit code) Crashed (reason)", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestCheckPodsNoOwner(t *testing.T) {
	dummyResources := v1.PodList{
		Items: []v1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "test",
					OwnerReferences: []metav1.OwnerReference{},
				},
				Status: v1.PodStatus{
					StartTime:  createNewTimeStampPtr(time.Now()),
					Phase:      "Running",
					Conditions: []v1.PodCondition{},
				},
			},
		},
	}

	result := CheckPods(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "has no owner", result.Symptoms[0].Message)
	assert.Equal(t, "warning", result.Symptoms[0].Severity)
}
