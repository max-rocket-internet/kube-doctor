// Package symptoms defines the findings produced by health checks.
package symptoms

// Severity describes how serious a Symptom is.
type Severity string

const (
	Warning  Severity = "warning"
	Critical Severity = "critical"
)

// Symptom is a single problem found on a Kubernetes resource.
type Symptom struct {
	Severity     Severity
	ResourceName string
	ResourceType string
	Namespace    string
	Message      string
}

// SymptomList is a collection of Symptoms.
type SymptomList struct {
	Symptoms []Symptom
}

// Add appends s to the list.
func (l *SymptomList) Add(s Symptom) {
	l.Symptoms = append(l.Symptoms, s)
}

// Counts holds the number of symptoms per severity.
type Counts struct {
	Critical int
	Warning  int
}

// Total returns the number of symptoms across all severities.
func (c Counts) Total() int {
	return c.Critical + c.Warning
}

// Count returns the number of symptoms per severity.
func (l SymptomList) Count() (c Counts) {
	for _, s := range l.Symptoms {
		switch s.Severity {
		case Critical:
			c.Critical++
		case Warning:
			c.Warning++
		}
	}
	return c
}

// ContainerSymptom is a problem found on a single container spec.
type ContainerSymptom struct {
	Name     string
	Severity Severity
	Message  string
}

// ContainerSymptomList is a collection of ContainerSymptoms.
type ContainerSymptomList struct {
	Symptoms []ContainerSymptom
}

// Add appends s to the list.
func (l *ContainerSymptomList) Add(s ContainerSymptom) {
	l.Symptoms = append(l.Symptoms, s)
}
