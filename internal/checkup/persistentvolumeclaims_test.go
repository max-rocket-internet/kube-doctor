package checkup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCheckPersistentVolumeClaimsOldAndNotBound(t *testing.T) {
	dummyResources := v1.PersistentVolumeClaimList{
		Items: []v1.PersistentVolumeClaim{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test",
					CreationTimestamp: metav1.NewTime(time.Now().Add(-10 * time.Minute)),
				},
				Status: v1.PersistentVolumeClaimStatus{
					Phase: "NotBound",
				},
			},
		},
	}

	result := CheckPersistentVolumeClaims(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "older than 5 minutes and status is not bound", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestCheckPersistentVolumeClaimsNewAndNotBound(t *testing.T) {
	dummyResources := v1.PersistentVolumeClaimList{
		Items: []v1.PersistentVolumeClaim{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test",
					CreationTimestamp: metav1.NewTime(time.Now().Add(-1 * time.Minute)),
				},
				Status: v1.PersistentVolumeClaimStatus{
					Phase: "NotBound",
				},
			},
		},
	}

	result := CheckPersistentVolumeClaims(&dummyResources)

	assert.Len(t, result.Symptoms, 0)
}
