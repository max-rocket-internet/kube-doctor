package checkup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
)

func TestDeploymentConditionMinimumReplicasAvailable(t *testing.T) {
	dummyResources := v1.DeploymentList{
		Items: []v1.Deployment{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1.DeploymentSpec{
					Replicas: new(int32(1)),
				},
				Status: v1.DeploymentStatus{
					Replicas:      1,
					ReadyReplicas: 1,
					Conditions: []v1.DeploymentCondition{
						{
							Status: "False",
							Reason: "MinimumReplicasAvailable",
						},
					},
				},
			},
		},
	}

	result := CheckDeployments(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "minimum availability not met", result.Symptoms[0].Message)
	assert.Equal(t, symptoms.Critical, result.Symptoms[0].Severity)
}

func TestDeploymentConditionReplicaSetUpdated(t *testing.T) {
	dummyResources := v1.DeploymentList{
		Items: []v1.Deployment{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1.DeploymentSpec{
					Replicas: new(int32(1)),
				},
				Status: v1.DeploymentStatus{
					Replicas:      1,
					ReadyReplicas: 1,
					Conditions: []v1.DeploymentCondition{
						{
							Reason:         "ReplicaSetUpdated",
							Type:           "Progressing",
							LastUpdateTime: metav1.NewTime(time.Now()),
						},
					},
				},
			},
		},
	}

	result := CheckDeployments(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "ReplicaSet update in progress", result.Symptoms[0].Message)
	assert.Equal(t, symptoms.Warning, result.Symptoms[0].Severity)
}

func TestDeploymentConditionReplicaSetUpdatedTooLong(t *testing.T) {
	dummyResources := v1.DeploymentList{
		Items: []v1.Deployment{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1.DeploymentSpec{
					Replicas: new(int32(1)),
				},
				Status: v1.DeploymentStatus{
					Replicas:      1,
					ReadyReplicas: 1,
					Conditions: []v1.DeploymentCondition{
						{
							Reason:         "ReplicaSetUpdated",
							Type:           "Progressing",
							LastUpdateTime: metav1.NewTime(time.Now().Add(-20 * time.Minute)),
						},
					},
				},
			},
		},
	}

	result := CheckDeployments(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "ReplicaSet update in progress but no progress for 10 minutes or longer", result.Symptoms[0].Message)
	assert.Equal(t, symptoms.Critical, result.Symptoms[0].Severity)
}

func TestDeploymentPodsNotReady(t *testing.T) {
	dummyResources := v1.DeploymentList{
		Items: []v1.Deployment{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1.DeploymentSpec{
					Replicas: new(int32(3)),
				},
				Status: v1.DeploymentStatus{
					Replicas:      1,
					ReadyReplicas: 1,
				},
			},
		},
	}

	result := CheckDeployments(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "2/3 pods are not ready", result.Symptoms[0].Message)
	assert.Equal(t, symptoms.Warning, result.Symptoms[0].Severity)
}
