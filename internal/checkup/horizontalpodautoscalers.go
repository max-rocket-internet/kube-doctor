package checkup

import (
	"fmt"

	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

func CheckHPAs(resources *autoscalingv2.HorizontalPodAutoscalerList) (results symptoms.SymptomList) {
	const resourceType = "HorizontalPodAutoscaler"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, hpa := range resources.Items {
		log.Debugf("Examining HorizontalPodAutoscaler %s/%s", hpa.Namespace, hpa.Name)

		// The API server defaults spec.minReplicas to 1 when unset.
		minReplicas := ptr.Deref(hpa.Spec.MinReplicas, 1)

		for _, condition := range hpa.Status.Conditions {
			if isBadHPACondition(condition, hpa.Status.CurrentReplicas, minReplicas) {
				results.Add(symptoms.Symptom{
					Message:      fmt.Sprintf("has condition %s=%s and reason %s", condition.Type, condition.Status, condition.Reason),
					Severity:     symptoms.Warning,
					ResourceName: hpa.Name,
					ResourceType: resourceType,
					Namespace:    hpa.Namespace,
				})
			}
		}
	}

	log.PrintEnd(len(resources.Items), results.Count())

	return results
}

func isBadHPACondition(c autoscalingv2.HorizontalPodAutoscalerCondition, currentReplicas, minReplicas int32) bool {
	switch {
	case c.Status == corev1.ConditionUnknown:
		return true
	case c.Type == autoscalingv2.AbleToScale, c.Type == autoscalingv2.ScalingActive:
		return c.Status == corev1.ConditionFalse
	case c.Type == autoscalingv2.ScalingLimited:
		return c.Status == corev1.ConditionTrue && currentReplicas != minReplicas
	default:
		return false
	}
}
