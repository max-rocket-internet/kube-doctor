package checkup

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

// The amount of noise in events is very high and it's hard to filter accurately into meaningful symptoms.
// The logic in here is the most opinionated of all symptoms.
// PRs welcome.

func CheckEvents(resources *corev1.EventList) (results symptoms.SymptomList) {
	const resourceType = "Event"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, event := range resources.Items {
		if !isSymptomaticEvent(event) {
			continue
		}

		results.Add(symptoms.Symptom{
			Message:      fmt.Sprintf("(%s) %.1f minutes ago: %s", event.InvolvedObject.Kind, time.Since(event.LastTimestamp.Time).Minutes(), event.Message),
			Severity:     symptoms.Critical,
			ResourceName: event.InvolvedObject.Name,
			ResourceType: resourceType,
			Namespace:    event.InvolvedObject.Namespace,
		})
	}

	log.PrintEnd(len(resources.Items), results.Count())

	return results
}

func isSymptomaticEvent(event corev1.Event) bool {
	abnormal := event.Type != corev1.EventTypeNormal

	switch event.Source.Component {
	case "cluster-autoscaler":
		return abnormal || event.Reason == "ScaleDown" || event.Reason == "TriggeredScaleUp"
	case "service-controller":
		return abnormal
	case "default-scheduler":
		return abnormal && event.Reason != "FailedScheduling"
	case "kubelet":
		return abnormal && event.Reason != "Unhealthy"
	default:
		return false
	}
}
