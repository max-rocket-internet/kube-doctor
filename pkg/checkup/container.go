package checkup

import (
	"fmt"

	"github.com/max-rocket-internet/kube-doctor/pkg/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/pkg/log"
	v1 "k8s.io/api/core/v1"
)

func checkContainer(container v1.Container) (results symptoms.ContainerSymptomList) {
	log.Debug(fmt.Sprintf("Examining Container %s", container.Name))

	if len(container.Resources.Requests) == 0 && len(container.Resources.Limits) == 0 {
		results.Add(symptoms.ContainerSymptom{
			Name:     container.Name,
			Message:  "no resources specified",
			Severity: "warning",
		})
	} else if container.Resources.Limits.Memory().AsApproximateFloat64() == 0 && container.Resources.Requests.Memory().AsApproximateFloat64() == 0 {
		results.Add(symptoms.ContainerSymptom{
			Name:     container.Name,
			Message:  "no memory resources specified",
			Severity: "warning",
		})
	} else if container.Resources.Limits.Memory().AsApproximateFloat64() == 0 {
		results.Add(symptoms.ContainerSymptom{
			Name:     container.Name,
			Message:  "no memory limit",
			Severity: "warning",
		})
	} else if container.Resources.Limits.Memory().AsApproximateFloat64() != container.Resources.Requests.Memory().AsApproximateFloat64() {
		results.Add(symptoms.ContainerSymptom{
			Name:     container.Name,
			Message:  "memory request and limit are not equal",
			Severity: "critical",
		})
	}

	return results
}
