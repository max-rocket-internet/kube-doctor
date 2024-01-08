package checkup

import (
	"fmt"
	"time"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"

	appsv1 "k8s.io/api/apps/v1"
)

func CheckDeployments(resources *appsv1.DeploymentList) (results symptoms.SymptomList) {
	resourceType := "Deployment"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, deployment := range resources.Items {
		log.Debug(fmt.Sprintf("Examining Deployment %s/%s", deployment.Namespace, deployment.Name))

		for _, container := range deployment.Spec.Template.Spec.Containers {
			for _, s := range checkContainer(container).Symptoms {
				results.Add(symptoms.Symptom{
					Message:      fmt.Sprintf("container '%s' %s", container.Name, s.Message),
					Severity:     s.Severity,
					ResourceName: deployment.Name,
					ResourceType: resourceType,
					Namespace:    deployment.Namespace,
				})
			}
		}

		for _, condition := range deployment.Status.Conditions {
			if condition.Reason == "MinimumReplicasAvailable" && condition.Status != "True" {
				results.Add(symptoms.Symptom{
					Message:      "minimum availability not met",
					Severity:     "critical",
					ResourceName: deployment.Name,
					ResourceType: resourceType,
					Namespace:    deployment.Namespace,
				})
			}
			if condition.Reason == "ReplicaSetUpdated" && condition.Type == "Progressing" {
				if time.Since(condition.LastUpdateTime.Time).Minutes() > 10 {
					results.Add(symptoms.Symptom{
						Message:      "ReplicaSet update in progress but no progress for 10 minutes or longer",
						Severity:     "critical",
						ResourceName: deployment.Name,
						ResourceType: resourceType,
						Namespace:    deployment.Namespace,
					})
				} else {
					results.Add(symptoms.Symptom{
						Message:      "ReplicaSet update in progress",
						Severity:     "warning",
						ResourceName: deployment.Name,
						ResourceType: resourceType,
						Namespace:    deployment.Namespace,
					})
				}
			}
		}

		if deployment.Status.ReadyReplicas != *deployment.Spec.Replicas {
			results.Add(symptoms.Symptom{
				Message:      fmt.Sprintf("%d/%d pods are not ready", *deployment.Spec.Replicas-deployment.Status.ReadyReplicas, *deployment.Spec.Replicas),
				Severity:     "warning",
				ResourceName: deployment.Name,
				ResourceType: resourceType,
				Namespace:    deployment.Namespace,
			})
		}
	}

	log.PrintEnd(len(resources.Items), results.CountSymptomsSeverity())

	return results
}
