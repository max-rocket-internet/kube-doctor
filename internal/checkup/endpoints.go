package checkup

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	discoveryv1 "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

type endpointCounts struct {
	ready, notReady int
}

// CheckEndpointSlices aggregates EndpointSlices per Service and reports
// Services whose endpoints are wholly or partially not ready.
func CheckEndpointSlices(resources *discoveryv1.EndpointSliceList) (results symptoms.SymptomList) {
	const resourceType = "EndpointSlice"

	log.PrintBegin(len(resources.Items), resourceType)

	services := make(map[types.NamespacedName]*endpointCounts)

	for _, slice := range resources.Items {
		log.Debugf("Examining EndpointSlice %s/%s", slice.Namespace, slice.Name)

		// Slices not managed for a Service are reported under their own name.
		name := slice.Name
		if svc, ok := slice.Labels[discoveryv1.LabelServiceName]; ok {
			name = svc
		}
		key := types.NamespacedName{Namespace: slice.Namespace, Name: name}

		counts, ok := services[key]
		if !ok {
			counts = &endpointCounts{}
			services[key] = counts
		}

		for _, endpoint := range slice.Endpoints {
			// Terminating endpoints are expected to be going away.
			if endpoint.Conditions.Terminating != nil && *endpoint.Conditions.Terminating {
				continue
			}
			// A nil Ready condition must be interpreted as ready.
			if endpoint.Conditions.Ready == nil || *endpoint.Conditions.Ready {
				counts.ready++
			} else {
				counts.notReady++
			}
		}
	}

	keys := slices.SortedFunc(maps.Keys(services), func(a, b types.NamespacedName) int {
		return strings.Compare(a.String(), b.String())
	})

	for _, key := range keys {
		counts := services[key]

		add := func(severity symptoms.Severity, message string) {
			results.Add(symptoms.Symptom{
				Message:      message,
				Severity:     severity,
				ResourceName: key.Name,
				ResourceType: resourceType,
				Namespace:    key.Namespace,
			})
		}

		switch {
		case counts.ready == 0 && counts.notReady > 0:
			add(symptoms.Critical, "no ready endpoints")
		case counts.notReady > 0:
			add(symptoms.Warning, fmt.Sprintf("%d/%d endpoints are not ready", counts.notReady, counts.notReady+counts.ready))
		}
	}

	log.PrintEnd(len(resources.Items), results.Count())

	return results
}
