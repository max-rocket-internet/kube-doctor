package checkup

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

const (
	// Containers in pods younger than this are allowed to be not ready yet.
	containerStartupGracePeriod = 3 * time.Minute
	// Restarts older than this are reported as warnings rather than critical.
	recentRestartWindow = time.Hour
)

func CheckPods(resources *corev1.PodList) (results symptoms.SymptomList) {
	const resourceType = "Pod"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, pod := range resources.Items {
		log.Debugf("Examining Pod %s/%s", pod.Namespace, pod.Name)

		if pod.Status.Phase == corev1.PodSucceeded {
			continue
		}

		add := func(severity symptoms.Severity, message string) {
			results.Add(symptoms.Symptom{
				Message:      message,
				Severity:     severity,
				ResourceName: pod.Name,
				ResourceType: resourceType,
				Namespace:    pod.Namespace,
			})
		}

		if pod.Status.Phase != corev1.PodRunning {
			add(symptoms.Critical, "not running")
		}

		for _, condition := range pod.Status.Conditions {
			if condition.Status != corev1.ConditionTrue {
				add(symptoms.Critical, fmt.Sprintf("status condition %s is %s", condition.Type, condition.Status))
			}
		}

		// StartTime is unset until the kubelet has accepted the pod, so treat it as just started.
		var podAge time.Duration
		if pod.Status.StartTime != nil {
			podAge = time.Since(pod.Status.StartTime.Time)
		}

		for _, cs := range pod.Status.ContainerStatuses {
			if !cs.Ready {
				if podAge < containerStartupGracePeriod {
					add(symptoms.Warning, fmt.Sprintf("container '%s' is not ready but pod started %.1f mins ago", cs.Name, podAge.Minutes()))
				} else {
					add(symptoms.Critical, fmt.Sprintf("container '%s' is not ready", cs.Name))
				}
			}

			if terminated := cs.LastTerminationState.Terminated; cs.RestartCount > 0 && terminated != nil {
				sinceRestart := time.Since(terminated.FinishedAt.Time)
				if sinceRestart > recentRestartWindow {
					add(symptoms.Warning, fmt.Sprintf("container '%s' has been restarted %d times", cs.Name, cs.RestartCount))
				} else {
					add(symptoms.Critical, fmt.Sprintf("container '%s' was restarted %.1f mins ago: %d (exit code) %s (reason)",
						cs.Name, sinceRestart.Minutes(), terminated.ExitCode, terminated.Reason))
				}
			}
		}

		if len(pod.OwnerReferences) == 0 {
			add(symptoms.Warning, "has no owner")
		}
	}

	log.PrintEnd(len(resources.Items), results.Count())

	return results
}
