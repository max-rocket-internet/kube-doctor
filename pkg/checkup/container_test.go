package checkup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestContainerNoResources(t *testing.T) {
	container := v1.Container{
		Name:      "test",
		Resources: v1.ResourceRequirements{},
	}

	result := checkContainer(container)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "no resources specified", result.Symptoms[0].Message)
	assert.Equal(t, "warning", result.Symptoms[0].Severity)
}

func TestContainerNoMemoryResources(t *testing.T) {
	container := v1.Container{
		Name: "test",
		Resources: v1.ResourceRequirements{
			Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("100m")},
		},
	}

	result := checkContainer(container)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "no memory resources specified", result.Symptoms[0].Message)
	assert.Equal(t, "warning", result.Symptoms[0].Severity)
}

func TestContainerNoMemoryLimit(t *testing.T) {
	container := v1.Container{
		Name: "test",
		Resources: v1.ResourceRequirements{
			Requests: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("100m"),
				v1.ResourceMemory: resource.MustParse("1Gi"),
			},
		},
	}

	result := checkContainer(container)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "no memory limit", result.Symptoms[0].Message)
	assert.Equal(t, "warning", result.Symptoms[0].Severity)
}

func TestContainerMemoryRequestLimitNotEqual(t *testing.T) {
	container := v1.Container{
		Name: "test",
		Resources: v1.ResourceRequirements{
			Requests: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("100m"),
				v1.ResourceMemory: resource.MustParse("1Gi"),
			},
			Limits: v1.ResourceList{
				v1.ResourceMemory: resource.MustParse("2Gi"),
			},
		},
	}

	result := checkContainer(container)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "memory request and limit are not equal", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}
