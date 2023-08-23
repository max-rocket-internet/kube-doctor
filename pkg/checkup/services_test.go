package checkup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCheckServices(t *testing.T) {
	dummyResources := v1.ServiceList{
		Items: []v1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1.ServiceSpec{
					Type: "LoadBalancer",
				},
				Status: v1.ServiceStatus{
					LoadBalancer: v1.LoadBalancerStatus{
						Ingress: []v1.LoadBalancerIngress{},
					},
				},
			},
		},
	}

	result := CheckServices(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "LoadBalancer service has no ingress points", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestCheckServicesIgnored(t *testing.T) {
	dummyResources := v1.ServiceList{
		Items: []v1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: v1.ServiceSpec{
					Type: "ClusterIP",
				},
				Status: v1.ServiceStatus{},
			},
		},
	}

	result := CheckServices(&dummyResources)

	assert.Len(t, result.Symptoms, 0)
}
