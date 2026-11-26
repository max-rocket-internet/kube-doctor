package checkup

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

func CheckServices(resources *corev1.ServiceList) (results symptoms.SymptomList) {
	const resourceType = "Service"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, service := range resources.Items {
		log.Debugf("Examining Service %s/%s", service.Namespace, service.Name)

		if service.Spec.Type == corev1.ServiceTypeLoadBalancer && len(service.Status.LoadBalancer.Ingress) == 0 {
			results.Add(symptoms.Symptom{
				Message:      "LoadBalancer service has no ingress points",
				Severity:     symptoms.Critical,
				ResourceName: service.Name,
				ResourceType: resourceType,
				Namespace:    service.Namespace,
			})
		}
	}

	log.PrintEnd(len(resources.Items), results.Count())

	return results
}
