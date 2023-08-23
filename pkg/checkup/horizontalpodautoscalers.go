package checkup

import (
	"fmt"

	"github.com/max-rocket-internet/kube-doctor/pkg/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/pkg/log"
	autoscaling "k8s.io/api/autoscaling/v2"
)

func CheckHpas(resources *autoscaling.HorizontalPodAutoscalerList) (results symptoms.SymptomList) {
	resourceType := "HorizontalPodAutoscaler"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, hpa := range resources.Items {
		log.Debug(fmt.Sprintf("Examining HorizontalPodAutoscaler %s/%s", hpa.Namespace, hpa.Name))

		for _, condition := range hpa.Status.Conditions {
			if (condition.Status == "Unknown") || (condition.Type == "AbleToScale" && condition.Status == "False") || (condition.Type == "ScalingActive" && condition.Status == "False") || (condition.Type == "ScalingLimited" && condition.Status == "True" && hpa.Status.CurrentReplicas != *hpa.Spec.MinReplicas) {
				results.Add(symptoms.Symptom{
					Message:      fmt.Sprintf("has condition %s=%s and reason %s", condition.Type, condition.Status, condition.Reason),
					Severity:     "warning",
					ResourceName: hpa.ObjectMeta.Name,
					ResourceType: resourceType,
					Namespace:    hpa.ObjectMeta.Namespace,
				})
			}
		}
	}

	log.PrintEnd(len(resources.Items), len(results.Symptoms))

	return results
}
