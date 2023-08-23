package checkup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCheckPersistentVolumeOldAndNotBound(t *testing.T) {
	dummyResources := v1.PersistentVolumeList{
		Items: []v1.PersistentVolume{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test",
					CreationTimestamp: metav1.NewTime(time.Now().Add(-10 * time.Minute)),
				},
				Status: v1.PersistentVolumeStatus{
					Phase: "NotBound",
				},
			},
		},
	}

	result := CheckPersistentVolumes(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "older than 5 minutes and status is not bound", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestCheckPersistentVolumeNewAndNotBound(t *testing.T) {
	dummyResources := v1.PersistentVolumeList{
		Items: []v1.PersistentVolume{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test",
					CreationTimestamp: metav1.NewTime(time.Now().Add(-1 * time.Minute)),
				},
				Status: v1.PersistentVolumeStatus{
					Phase: "NotBound",
				},
			},
		},
	}

	result := CheckPersistentVolumes(&dummyResources)

	assert.Len(t, result.Symptoms, 0)
}
