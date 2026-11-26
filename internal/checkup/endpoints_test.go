package checkup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
)

func endpointSlice(service string, ready ...bool) discoveryv1.EndpointSlice {
	slice := discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service + "-abcde",
			Namespace: "default",
			Labels:    map[string]string{discoveryv1.LabelServiceName: service},
		},
	}
	for _, r := range ready {
		slice.Endpoints = append(slice.Endpoints, discoveryv1.Endpoint{
			Conditions: discoveryv1.EndpointConditions{Ready: new(r)},
		})
	}
	return slice
}

func TestEndpointSlices(t *testing.T) {
	tests := []struct {
		name         string
		slices       []discoveryv1.EndpointSlice
		wantMessage  string
		wantSeverity symptoms.Severity
	}{
		{
			name:         "no ready endpoints",
			slices:       []discoveryv1.EndpointSlice{endpointSlice("app", false)},
			wantMessage:  "no ready endpoints",
			wantSeverity: symptoms.Critical,
		},
		{
			name:         "some endpoints not ready",
			slices:       []discoveryv1.EndpointSlice{endpointSlice("app", true, true, false)},
			wantMessage:  "1/3 endpoints are not ready",
			wantSeverity: symptoms.Warning,
		},
		{
			name:         "aggregated across slices of one service",
			slices:       []discoveryv1.EndpointSlice{endpointSlice("app", false), endpointSlice("app", true)},
			wantMessage:  "1/2 endpoints are not ready",
			wantSeverity: symptoms.Warning,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckEndpointSlices(&discoveryv1.EndpointSliceList{Items: tt.slices})

			assert.Len(t, result.Symptoms, 1)
			assert.Equal(t, "app", result.Symptoms[0].ResourceName)
			assert.Equal(t, tt.wantMessage, result.Symptoms[0].Message)
			assert.Equal(t, tt.wantSeverity, result.Symptoms[0].Severity)
		})
	}
}

func TestEndpointSlicesHealthy(t *testing.T) {
	terminating := endpointSlice("app", true)
	terminating.Endpoints = append(terminating.Endpoints, discoveryv1.Endpoint{
		Conditions: discoveryv1.EndpointConditions{Ready: new(false), Terminating: new(true)},
	})

	resources := discoveryv1.EndpointSliceList{Items: []discoveryv1.EndpointSlice{
		endpointSlice("ready", true, true),
		endpointSlice("empty"),
		{Endpoints: []discoveryv1.Endpoint{{}}}, // nil Ready means ready
		terminating,
	}}

	assert.Empty(t, CheckEndpointSlices(&resources).Symptoms)
}
