package checkup

import (
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

// Jobs that completed longer ago than this are ignored.
const jobMaxAge = time.Hour

func CheckJobs(resources *batchv1.JobList) (results symptoms.SymptomList) {
	const resourceType = "Job"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, job := range resources.Items {
		log.Debugf("Examining Job %s/%s", job.Namespace, job.Name)

		if job.Status.CompletionTime != nil && time.Since(job.Status.CompletionTime.Time) > jobMaxAge {
			continue
		}

		if job.Status.Failed == 0 {
			continue
		}

		for _, condition := range job.Status.Conditions {
			if condition.Type == batchv1.JobFailed && condition.Status == corev1.ConditionTrue {
				results.Add(symptoms.Symptom{
					Message:      fmt.Sprintf("%s: %s", condition.Reason, condition.Message),
					Severity:     symptoms.Critical,
					ResourceName: job.Name,
					ResourceType: resourceType,
					Namespace:    job.Namespace,
				})
			}
		}
	}

	log.PrintEnd(len(resources.Items), results.Count())

	return results
}
