package checkup

import (
	"fmt"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/kubernetes/statuses"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

func CheckKubeAPIHealth(resources *statuses.KubeAPIHealthEndpointStatusList) (results symptoms.SymptomList) {
	const resourceType = "KubeApiHealthEndpointStatus"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, status := range resources.Items {
		log.Debugf("Examining KubeApiHealthEndpointStatus %s", status.Name)

		if status.Status != statuses.StatusOK {
			results.Add(symptoms.Symptom{
				Message:      fmt.Sprintf("Kubernetes API health endpoint '%s' %s=%s", status.Path, status.Name, status.Status),
				Severity:     symptoms.Critical,
				ResourceType: resourceType,
			})
		}
	}

	log.PrintEnd(len(resources.Items), results.Count())

	return results
}
