package checkup

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
)

func TestCheckJobs(t *testing.T) {
	dummyResources := batchv1.JobList{
		Items: []batchv1.Job{
			{
				Status: batchv1.JobStatus{
					Failed:         1,
					CompletionTime: new(metav1.NewTime(time.Now())),
					Conditions: []batchv1.JobCondition{
						{
							Status:  "True",
							Type:    "Failed",
							Message: "SomeMessage",
							Reason:  "SomeReason",
						},
					},
				},
			},
			{
				Status: batchv1.JobStatus{
					Failed:         0,
					CompletionTime: new(metav1.NewTime(time.Now())),
					Conditions: []batchv1.JobCondition{
						{
							Status:  "True",
							Type:    "Failed",
							Message: "SomeMessage",
							Reason:  "SomeReason",
						},
					},
				},
			},
			{
				Status: batchv1.JobStatus{
					Failed:         1,
					CompletionTime: new(metav1.NewTime(time.Now().Add(-2 * time.Hour))),
					Conditions: []batchv1.JobCondition{
						{
							Status:  "True",
							Type:    "Failed",
							Message: "SomeMessage",
							Reason:  "SomeReason",
						},
					},
				},
			},
		},
	}

	result := CheckJobs(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "SomeReason: SomeMessage", result.Symptoms[0].Message)
	assert.Equal(t, symptoms.Critical, result.Symptoms[0].Severity)
}
