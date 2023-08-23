package checkup

import (
	"fmt"
	"time"

	"github.com/max-rocket-internet/kube-doctor/pkg/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/pkg/log"
	v1 "k8s.io/api/core/v1"
)

func CheckPersistentVolumeClaims(resources *v1.PersistentVolumeClaimList) (results symptoms.SymptomList) {
	resourceType := "PersistentVolumeClaim"

	log.PrintBegin(len(resources.Items), resourceType)

	for _, pvc := range resources.Items {
		log.Debug(fmt.Sprintf("Examining PersistentVolumeClaim %s/%s", pvc.Name, pvc.Namespace))

		if pvc.Status.Phase != "Bound" && time.Now().Sub(pvc.CreationTimestamp.Time).Minutes() > 5 {
			results.Add(symptoms.Symptom{
				Message:      "older than 5 minutes and status is not bound",
				Severity:     "critical",
				ResourceName: pvc.Name,
				ResourceType: resourceType,
			})
		}
	}

	log.PrintEnd(len(resources.Items), len(results.Symptoms))

	return results
}
