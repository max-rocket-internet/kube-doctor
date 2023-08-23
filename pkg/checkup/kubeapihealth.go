package checkup

import (
	"fmt"

	"github.com/max-rocket-internet/kube-doctor/pkg/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/pkg/kubernetes/statuses"
	"github.com/max-rocket-internet/kube-doctor/pkg/log"
)

func KubeApiHealthStatuses(resources *statuses.KubeApiHealthEndpointStatusList) (results symptoms.SymptomList) {
	resourceType := "KubeApiHealthEndpointStatus"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, status := range resources.Items {
		log.Debug(fmt.Sprintf("Examining KubeApiHealthEndpointStatus %s", status.Name))

		if status.Status != "ok" {
			results.Add(symptoms.Symptom{
				Message:      fmt.Sprintf("Kubnernetes API health endpoint '%s' %s=%s", status.Path, status.Name, status.Status),
				Severity:     "critical",
				ResourceType: resourceType,
			})
		}
	}

	log.PrintEnd(len(resources.Items), len(results.Symptoms))

	return results
}
