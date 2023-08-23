package checkup

import (
	"fmt"
	"time"

	"github.com/max-rocket-internet/kube-doctor/pkg/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/pkg/log"

	v1 "k8s.io/api/core/v1"
)

func CheckNodes(resources *v1.NodeList) (results symptoms.SymptomList) {
	nodeVersions := make(map[string]int)
	resourceType := "Node"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, node := range resources.Items {
		log.Debug(fmt.Sprintf("Examining Node %s", node.Name))

		nodeVersions[node.Status.NodeInfo.KubeletVersion]++

		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" {
				if condition.Status != "True" {
					if time.Now().Sub(node.ObjectMeta.CreationTimestamp.Time).Minutes() > 5 {
						results.Add(symptoms.Symptom{
							Message:      "older than 5 minutes and not Ready",
							Severity:     "critical",
							ResourceName: node.Name,
							ResourceType: resourceType,
						})
					} else {
						results.Add(symptoms.Symptom{
							Message:      "not ready",
							Severity:     "warning",
							ResourceName: node.Name,
							ResourceType: resourceType,
						})
					}
				}
			} else if condition.Status != "False" {
				results.Add(symptoms.Symptom{
					Message:      fmt.Sprintf("has condition %s=%s: %s", condition.Type, condition.Status, condition.Reason),
					Severity:     "critical",
					ResourceName: node.Name,
					ResourceType: resourceType,
				})
			}
		}
	}

	if len(nodeVersions) > 1 {
		results.Add(symptoms.Symptom{
			Message:      "multiple node versions running",
			Severity:     "critical",
			ResourceType: resourceType,
		})
	}

	log.PrintEnd(len(resources.Items), len(results.Symptoms))

	return results
}
