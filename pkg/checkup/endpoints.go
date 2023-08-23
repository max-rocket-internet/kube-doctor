package checkup

import (
	"fmt"

	v1 "k8s.io/api/core/v1"

	"github.com/max-rocket-internet/kube-doctor/pkg/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/pkg/log"
)

func CheckEndpoints(resources *v1.EndpointsList) (results symptoms.SymptomList) {
	resourceType := "Endpoint"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, endpoint := range resources.Items {
		log.Debug(fmt.Sprintf("Examining Endpoint %s/%s", endpoint.Namespace, endpoint.Name))

		for _, s := range endpoint.Subsets {
			readyAddressCount := len(s.Addresses)
			notReadyAddressCount := len(s.NotReadyAddresses)

			if readyAddressCount == 0 {
				results.Add(symptoms.Symptom{
					Message:      "no ready addresses in subsets",
					Severity:     "critical",
					ResourceName: endpoint.ObjectMeta.Name,
					ResourceType: resourceType,
					Namespace:    endpoint.ObjectMeta.Namespace,
				})
			} else if len(s.NotReadyAddresses) > 0 {
				results.Add(symptoms.Symptom{
					Message:      fmt.Sprintf("%d/%d addresses in subsets are NotReady", notReadyAddressCount, notReadyAddressCount+readyAddressCount),
					Severity:     "warning",
					ResourceName: endpoint.ObjectMeta.Name,
					ResourceType: resourceType,
					Namespace:    endpoint.ObjectMeta.Namespace,
				})
			}
		}
	}

	log.PrintEnd(len(resources.Items), results.CountSymptomsSeverity())

	return results
}
