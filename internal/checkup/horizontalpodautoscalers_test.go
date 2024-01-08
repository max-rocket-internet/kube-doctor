package checkup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	autoscaling "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCheckHpasUknown(t *testing.T) {
	dummyResources := autoscaling.HorizontalPodAutoscalerList{
		Items: []autoscaling.HorizontalPodAutoscaler{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "test",
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					Conditions: []autoscaling.HorizontalPodAutoscalerCondition{
						{
							Status: "Unknown",
							Type:   "AbleToScale",
							Reason: "SomeReason",
						},
					},
				},
			},
		},
	}

	result := CheckHpas(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "has condition AbleToScale=Unknown and reason SomeReason", result.Symptoms[0].Message)
	assert.Equal(t, "warning", result.Symptoms[0].Severity)
}

func TestCheckHpasNotAbleToScale(t *testing.T) {
	dummyResources := autoscaling.HorizontalPodAutoscalerList{
		Items: []autoscaling.HorizontalPodAutoscaler{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "test",
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					Conditions: []autoscaling.HorizontalPodAutoscalerCondition{
						{
							Status: "False",
							Type:   "AbleToScale",
							Reason: "FailedGetScale",
						},
					},
				},
			},
		},
	}

	result := CheckHpas(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "has condition AbleToScale=False and reason FailedGetScale", result.Symptoms[0].Message)
	assert.Equal(t, "warning", result.Symptoms[0].Severity)
}

func TestCheckHpasToBeIgnored(t *testing.T) {
	dummyResources := autoscaling.HorizontalPodAutoscalerList{
		Items: []autoscaling.HorizontalPodAutoscaler{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "test",
				},
				Status: autoscaling.HorizontalPodAutoscalerStatus{
					Conditions: []autoscaling.HorizontalPodAutoscalerCondition{
						{
							Status: "False",
							Type:   "ScalingLimited",
							Reason: "DesiredWithinRange",
						},
					},
				},
			},
		},
	}

	result := CheckHpas(&dummyResources)

	assert.Len(t, result.Symptoms, 0)
}
