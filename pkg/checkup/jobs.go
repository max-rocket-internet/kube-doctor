package checkup

import (
	"fmt"
	"time"

	"github.com/max-rocket-internet/kube-doctor/pkg/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/pkg/log"
	batchv1 "k8s.io/api/batch/v1"
)

func CheckJobs(resources *batchv1.JobList) (results symptoms.SymptomList) {
	resourceType := "Job"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, job := range resources.Items {
		log.Debug(fmt.Sprintf("Examining Job %s/%s", job.Namespace, job.Name))

		// Ignore jobs older than 1 hour
		if job.Status.CompletionTime != nil && time.Since(job.Status.CompletionTime.Time).Minutes() > 60 {
			continue
		}

		if job.Status.Failed == 0 {
			continue
		}

		for _, condition := range job.Status.Conditions {
			if condition.Type == "Failed" && condition.Status == "True" {
				results.Add(symptoms.Symptom{
					Message:      fmt.Sprintf("%s: %s", condition.Reason, condition.Message),
					Severity:     "critical",
					ResourceName: job.Name,
					ResourceType: resourceType,
					Namespace:    job.Namespace,
				})
			}
		}
	}

	log.PrintEnd(len(resources.Items), len(results.Symptoms))

	return results
}
