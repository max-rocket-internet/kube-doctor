package checkup

import (
	"fmt"

	"github.com/max-rocket-internet/kube-doctor/pkg/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/pkg/log"
	v1 "k8s.io/api/core/v1"
)

func CheckServices(resources *v1.ServiceList) (results symptoms.SymptomList) {
	resourceType := "Service"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, service := range resources.Items {
		log.Debug(fmt.Sprintf("Examining Service %s/%s", service.Namespace, service.Name))

		if service.Spec.Type == "LoadBalancer" {
			if len(service.Status.LoadBalancer.Ingress) == 0 {
				results.Add(symptoms.Symptom{
					Message:      "LoadBalancer service has no ingress points",
					Severity:     "critical",
					ResourceName: service.ObjectMeta.Name,
					ResourceType: resourceType,
					Namespace:    service.ObjectMeta.Namespace,
				})
			}
		}
	}

	log.PrintEnd(len(resources.Items), len(results.Symptoms))

	return results
}
