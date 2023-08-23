package checkup

import (
	"fmt"
	"time"

	"github.com/max-rocket-internet/kube-doctor/pkg/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/pkg/log"
	v1 "k8s.io/api/core/v1"
)

func CheckPods(resources *v1.PodList) (results symptoms.SymptomList) {
	resourceType := "Pod"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, pod := range resources.Items {
		log.Debug(fmt.Sprintf("Examining Pod %s/%s", pod.Namespace, pod.Name))

		if pod.Status.Phase == "Succeeded" {
			return
		}

		if pod.Status.Phase != "Running" {
			results.Add(symptoms.Symptom{
				Message:      "not running",
				Severity:     "critical",
				ResourceName: pod.Name,
				ResourceType: resourceType,
				Namespace:    pod.Namespace,
			})
		}

		for _, sc := range pod.Status.Conditions {
			if sc.Status != "True" {
				results.Add(symptoms.Symptom{
					Message:      fmt.Sprintf("status condition %s is %s", sc.Type, sc.Status),
					Severity:     "critical",
					ResourceName: pod.Name,
					ResourceType: resourceType,
					Namespace:    pod.Namespace,
				})
			}
		}

		for _, scs := range pod.Status.ContainerStatuses {
			if !scs.Ready {
				if time.Since(pod.Status.StartTime.Time).Minutes() < 3 {
					results.Add(symptoms.Symptom{
						Message:      fmt.Sprintf("container '%s' is not ready but pod started %.1f mins ago", scs.Name, time.Since(pod.Status.StartTime.Time).Minutes()),
						Severity:     "warning",
						ResourceName: pod.Name,
						ResourceType: resourceType,
						Namespace:    pod.Namespace,
					})
				} else {
					results.Add(symptoms.Symptom{
						Message:      fmt.Sprintf("container '%s' is not ready", scs.Name),
						Severity:     "critical",
						ResourceName: pod.Name,
						ResourceType: resourceType,
						Namespace:    pod.Namespace,
					})
				}
			}

			if scs.RestartCount != 0 {
				if time.Since(scs.LastTerminationState.Terminated.FinishedAt.Time).Hours() > 1 {
					results.Add(symptoms.Symptom{
						Message:      fmt.Sprintf("container '%s' has been restarted %d times", scs.Name, scs.RestartCount),
						Severity:     "warning",
						ResourceName: pod.Name,
						ResourceType: resourceType,
						Namespace:    pod.Namespace,
					})
				} else {
					results.Add(symptoms.Symptom{
						Message: fmt.Sprintf("container '%s' was restarted %.1f mins ago: %d (exit code) %s (reason)",
							scs.Name,
							time.Since(scs.LastTerminationState.Terminated.FinishedAt.Time).Minutes(),
							scs.LastTerminationState.Terminated.ExitCode,
							scs.LastTerminationState.Terminated.Reason,
						),
						Severity:     "critical",
						ResourceName: pod.Name,
						ResourceType: resourceType,
						Namespace:    pod.Namespace,
					})
				}
			}
		}

		if len(pod.OwnerReferences) == 0 {
			results.Add(symptoms.Symptom{
				Message:      "has no owner",
				Severity:     "warning",
				ResourceName: pod.Name,
				ResourceType: resourceType,
				Namespace:    pod.Namespace,
			})
		}
	}

	log.PrintEnd(len(resources.Items), len(results.Symptoms))

	return results
}
