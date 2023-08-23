package checkup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestEndpointNoReadyAdresses(t *testing.T) {
	dummyResources := v1.EndpointsList{
		Items: []v1.Endpoints{
			{
				Subsets: []v1.EndpointSubset{
					{
						Addresses: []v1.EndpointAddress{},
						NotReadyAddresses: []v1.EndpointAddress{
							{
								IP: "10.0.0.3",
							},
						},
					},
				},
			},
		},
	}

	result := CheckEndpoints(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "no ready addresses in subsets", result.Symptoms[0].Message)
	assert.Equal(t, "critical", result.Symptoms[0].Severity)
}

func TestEndpointSomeNotReadyAddresses(t *testing.T) {
	dummyResources := v1.EndpointsList{
		Items: []v1.Endpoints{
			{
				Subsets: []v1.EndpointSubset{
					{
						Addresses: []v1.EndpointAddress{
							{
								IP: "10.0.0.1",
							},
							{
								IP: "10.0.0.2",
							},
						},
						NotReadyAddresses: []v1.EndpointAddress{
							{
								IP: "10.0.0.3",
							},
						},
					},
				},
			},
		},
	}

	result := CheckEndpoints(&dummyResources)

	assert.Len(t, result.Symptoms, 1)
	assert.Equal(t, "1/3 addresses in subsets are NotReady", result.Symptoms[0].Message)
	assert.Equal(t, "warning", result.Symptoms[0].Severity)
}
