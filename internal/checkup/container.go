package checkup

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

func checkContainer(container corev1.Container) (results symptoms.ContainerSymptomList) {
	log.Debugf("Examining Container %s", container.Name)

	add := func(message string) {
		results.Add(symptoms.ContainerSymptom{
			Name:     container.Name,
			Message:  message,
			Severity: symptoms.Warning,
		})
	}

	requests, limits := container.Resources.Requests, container.Resources.Limits
	memRequest, memLimit := requests.Memory(), limits.Memory()

	switch {
	case len(requests) == 0 && len(limits) == 0:
		add("no resources specified")
	case memLimit.IsZero() && memRequest.IsZero():
		add("no memory resources specified")
	case memLimit.IsZero():
		add("no memory limit")
	case memLimit.Cmp(*memRequest) != 0:
		add("memory request and limit are not equal")
	}

	return results
}
