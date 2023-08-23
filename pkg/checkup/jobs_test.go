package checkup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
)

func TestCheckJobs(t *testing.T) {
	dummyResources := batchv1.JobList{
		Items: []batchv1.Job{
			{
				Status: batchv1.JobStatus{
					Failed:         1,
					CompletionTime: createNewTimeStampPtr(time.Now()),
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
					CompletionTime: createNewTimeStampPtr(time.Now()),
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
					CompletionTime: createNewTimeStampPtr(time.Now().Add(-2 * time.Hour)),
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
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}
