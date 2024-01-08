package checkup

import (
	"fmt"
	"time"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
	v1 "k8s.io/api/core/v1"
)

func CheckPersistentVolumes(resources *v1.PersistentVolumeList) (results symptoms.SymptomList) {
	resourceType := "PersistentVolume"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, volume := range resources.Items {
		log.Debug(fmt.Sprintf("Examining PersistentVolume %s", volume.Name))

		if volume.Status.Phase != "Bound" && time.Since(volume.CreationTimestamp.Time).Minutes() > 5 {
			results.Add(symptoms.Symptom{
				Message:      "older than 5 minutes and status is not bound",
				Severity:     "critical",
				ResourceName: volume.Name,
				ResourceType: resourceType,
			})
		}
	}

	log.PrintEnd(len(resources.Items), results.CountSymptomsSeverity())

	return results
}
