package checkup

import (
	"fmt"
	"time"

	"github.com/max-rocket-internet/kube-doctor/pkg/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/pkg/log"
	v1 "k8s.io/api/core/v1"
)

// The amount of noise in events is very high and it's hard to filter accurately into meaningful symptoms.
// The logic in here is the most opinionated of all symptoms
// PRs welcome

func CheckEvents(resources *v1.EventList) (results symptoms.SymptomList) {
	resourceType := "Event"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, event := range resources.Items {
		if event.Source.Component == "cluster-autoscaler" {
			if event.Reason == "ScaleDown" || event.Reason == "TriggeredScaleUp" || event.Type != "Normal" {
				results.Add(symptoms.Symptom{
					Message:      fmt.Sprintf("(%s) %.1f minutes ago: %s", event.InvolvedObject.Kind, time.Since(event.LastTimestamp.Time).Minutes(), event.Message),
					Severity:     "critical",
					ResourceName: event.InvolvedObject.Name,
					ResourceType: resourceType,
					Namespace:    event.InvolvedObject.Namespace,
				})
			}
		}

		if event.Type != "Normal" && event.Source.Component == "service-controller" {
			results.Add(symptoms.Symptom{
				Message:      fmt.Sprintf("(%s) %.1f minutes ago: %s", event.InvolvedObject.Kind, time.Since(event.LastTimestamp.Time).Minutes(), event.Message),
				Severity:     "critical",
				ResourceName: event.InvolvedObject.Name,
				ResourceType: resourceType,
				Namespace:    event.InvolvedObject.Namespace,
			})
		}

		if event.Type != "Normal" && event.Source.Component == "default-scheduler" && event.Reason != "FailedScheduling" {
			results.Add(symptoms.Symptom{
				Message:      fmt.Sprintf("(%s) %.1f minutes ago: %s", event.InvolvedObject.Kind, time.Since(event.LastTimestamp.Time).Minutes(), event.Message),
				Severity:     "critical",
				ResourceName: event.InvolvedObject.Name,
				ResourceType: resourceType,
				Namespace:    event.InvolvedObject.Namespace,
			})
		}

		if event.Type != "Normal" && event.Source.Component == "kubelet" && event.Reason != "Unhealthy" {
			results.Add(symptoms.Symptom{
				Message:      fmt.Sprintf("(%s) %.1f minutes ago: %s", event.InvolvedObject.Kind, time.Since(event.LastTimestamp.Time).Minutes(), event.Message),
				Severity:     "critical",
				ResourceName: event.InvolvedObject.Name,
				ResourceType: resourceType,
				Namespace:    event.InvolvedObject.Namespace,
			})
		}
	}

	log.PrintEnd(len(resources.Items), results.CountSymptomsSeverity())

	return results
}
