package checkup

import (
	"fmt"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
	v1 "k8s.io/api/apps/v1"
)

func CheckDaemonSets(resources *v1.DaemonSetList) (results symptoms.SymptomList) {
	resourceType := "DaemonSet"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, daemonset := range resources.Items {
		log.Debug(fmt.Sprintf("Examining DaemonSet %s/%s", daemonset.Namespace, daemonset.Name))

		for _, container := range daemonset.Spec.Template.Spec.Containers {
			for _, s := range checkContainer(container).Symptoms {
				results.Add(symptoms.Symptom{
					Message:      fmt.Sprintf("%s %s", container.Name, s.Message),
					Severity:     s.Severity,
					ResourceName: daemonset.Name,
					ResourceType: resourceType,
					Namespace:    daemonset.Namespace,
				})
			}
		}

		if daemonset.Status.NumberUnavailable > 0 && (daemonset.Status.UpdatedNumberScheduled == daemonset.Status.DesiredNumberScheduled) {
			results.Add(symptoms.Symptom{
				Message:      fmt.Sprintf("%d/%d pods are not ready", daemonset.Status.NumberUnavailable, daemonset.Status.DesiredNumberScheduled),
				Severity:     "critical",
				ResourceName: daemonset.Name,
				ResourceType: resourceType,
				Namespace:    daemonset.Namespace,
			})
		}

		if daemonset.Status.NumberMisscheduled > 0 {
			results.Add(symptoms.Symptom{
				Message:      fmt.Sprintf("%d/%d pods are miss scheduled", daemonset.Status.NumberMisscheduled, daemonset.Status.DesiredNumberScheduled),
				Severity:     "critical",
				ResourceName: daemonset.Name,
				ResourceType: resourceType,
				Namespace:    daemonset.Namespace,
			})
		}

		if daemonset.Status.UpdatedNumberScheduled != daemonset.Status.DesiredNumberScheduled {
			results.Add(symptoms.Symptom{
				Message:      "rolling update in progress",
				Severity:     "warning",
				ResourceName: daemonset.Name,
				ResourceType: resourceType,
				Namespace:    daemonset.Namespace,
			})
		}
	}

	log.PrintEnd(len(resources.Items), results.CountSymptomsSeverity())

	return results
}
