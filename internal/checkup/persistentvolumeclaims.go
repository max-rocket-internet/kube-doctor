package checkup

import (
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

// Volumes and claims younger than this are allowed to be unbound.
const volumeBindGracePeriod = 5 * time.Minute

func CheckPersistentVolumeClaims(resources *corev1.PersistentVolumeClaimList) (results symptoms.SymptomList) {
	const resourceType = "PersistentVolumeClaim"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, pvc := range resources.Items {
		log.Debugf("Examining PersistentVolumeClaim %s/%s", pvc.Namespace, pvc.Name)

		if pvc.Status.Phase != corev1.ClaimBound && time.Since(pvc.CreationTimestamp.Time) > volumeBindGracePeriod {
			results.Add(symptoms.Symptom{
				Message:      "older than 5 minutes and status is not bound",
				Severity:     symptoms.Critical,
				ResourceName: pvc.Name,
				ResourceType: resourceType,
				Namespace:    pvc.Namespace,
			})
		}
	}

	log.PrintEnd(len(resources.Items), results.Count())

	return results
}
