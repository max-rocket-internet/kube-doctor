package checkup

import (
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

const deploymentProgressTimeout = 10 * time.Minute

func CheckDeployments(resources *appsv1.DeploymentList) (results symptoms.SymptomList) {
	const resourceType = "Deployment"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, deployment := range resources.Items {
		log.Debugf("Examining Deployment %s/%s", deployment.Namespace, deployment.Name)

		add := func(severity symptoms.Severity, message string) {
			results.Add(symptoms.Symptom{
				Message:      message,
				Severity:     severity,
				ResourceName: deployment.Name,
				ResourceType: resourceType,
				Namespace:    deployment.Namespace,
			})
		}

		for _, container := range deployment.Spec.Template.Spec.Containers {
			for _, s := range checkContainer(container).Symptoms {
				add(s.Severity, fmt.Sprintf("container '%s' %s", container.Name, s.Message))
			}
		}

		for _, condition := range deployment.Status.Conditions {
			if condition.Reason == "MinimumReplicasAvailable" && condition.Status != corev1.ConditionTrue {
				add(symptoms.Critical, "minimum availability not met")
			}

			if condition.Reason == "ReplicaSetUpdated" && condition.Type == appsv1.DeploymentProgressing {
				if time.Since(condition.LastUpdateTime.Time) > deploymentProgressTimeout {
					add(symptoms.Critical, "ReplicaSet update in progress but no progress for 10 minutes or longer")
				} else {
					add(symptoms.Warning, "ReplicaSet update in progress")
				}
			}
		}

		// The API server defaults spec.replicas to 1 when unset.
		desired := ptr.Deref(deployment.Spec.Replicas, 1)
		if deployment.Status.ReadyReplicas != desired {
			add(symptoms.Warning, fmt.Sprintf("%d/%d pods are not ready", desired-deployment.Status.ReadyReplicas, desired))
		}
	}

	log.PrintEnd(len(resources.Items), results.Count())

	return results
}
