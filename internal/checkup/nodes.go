package checkup

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

// Nodes younger than this are allowed to be not Ready yet.
const nodeStartupGracePeriod = 5 * time.Minute

func CheckNodes(resources *corev1.NodeList) (results symptoms.SymptomList) {
	const resourceType = "Node"

	nodeVersions := make(map[string]int)

	log.PrintBegin(len(resources.Items), resourceType)

	for _, node := range resources.Items {
		log.Debugf("Examining Node %s", node.Name)

		add := func(severity symptoms.Severity, message string) {
			results.Add(symptoms.Symptom{
				Message:      message,
				Severity:     severity,
				ResourceName: node.Name,
				ResourceType: resourceType,
			})
		}

		nodeVersions[node.Status.NodeInfo.KubeletVersion]++

		for _, condition := range node.Status.Conditions {
			switch {
			case condition.Type == corev1.NodeReady:
				if condition.Status == corev1.ConditionTrue {
					continue
				}
				if time.Since(node.CreationTimestamp.Time) > nodeStartupGracePeriod {
					add(symptoms.Critical, "older than 5 minutes and not Ready")
				} else {
					add(symptoms.Warning, "not ready")
				}
			case condition.Status != corev1.ConditionFalse:
				add(symptoms.Critical, fmt.Sprintf("has condition %s=%s: %s", condition.Type, condition.Status, condition.Reason))
			}
		}
	}

	if len(nodeVersions) > 1 {
		results.Add(symptoms.Symptom{
			Message:      "multiple node versions running",
			Severity:     symptoms.Critical,
			ResourceType: resourceType,
		})
	}

	log.PrintEnd(len(resources.Items), results.Count())

	return results
}
