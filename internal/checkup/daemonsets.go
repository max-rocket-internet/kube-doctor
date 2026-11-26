package checkup

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

func CheckDaemonSets(resources *appsv1.DaemonSetList) (results symptoms.SymptomList) {
	const resourceType = "DaemonSet"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, daemonset := range resources.Items {
		log.Debugf("Examining DaemonSet %s/%s", daemonset.Namespace, daemonset.Name)

		add := func(severity symptoms.Severity, message string) {
			results.Add(symptoms.Symptom{
				Message:      message,
				Severity:     severity,
				ResourceName: daemonset.Name,
				ResourceType: resourceType,
				Namespace:    daemonset.Namespace,
			})
		}

		for _, container := range daemonset.Spec.Template.Spec.Containers {
			for _, s := range checkContainer(container).Symptoms {
				add(s.Severity, fmt.Sprintf("%s %s", container.Name, s.Message))
			}
		}

		status := daemonset.Status
		rolloutComplete := status.UpdatedNumberScheduled == status.DesiredNumberScheduled

		if status.NumberUnavailable > 0 && rolloutComplete {
			add(symptoms.Critical, fmt.Sprintf("%d/%d pods are not ready", status.NumberUnavailable, status.DesiredNumberScheduled))
		}

		if status.NumberMisscheduled > 0 {
			add(symptoms.Critical, fmt.Sprintf("%d/%d pods are misscheduled", status.NumberMisscheduled, status.DesiredNumberScheduled))
		}

		if !rolloutComplete {
			add(symptoms.Warning, "rolling update in progress")
		}
	}

	log.PrintEnd(len(resources.Items), results.Count())

	return results
}
