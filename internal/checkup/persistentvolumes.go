package checkup

import (
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

func CheckPersistentVolumes(resources *corev1.PersistentVolumeList) (results symptoms.SymptomList) {
	const resourceType = "PersistentVolume"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, volume := range resources.Items {
		log.Debugf("Examining PersistentVolume %s", volume.Name)

		if volume.Status.Phase != corev1.VolumeBound && time.Since(volume.CreationTimestamp.Time) > volumeBindGracePeriod {
			results.Add(symptoms.Symptom{
				Message:      "older than 5 minutes and status is not bound",
				Severity:     symptoms.Critical,
				ResourceName: volume.Name,
				ResourceType: resourceType,
			})
		}
	}

	log.PrintEnd(len(resources.Items), results.Count())

	return results
}
