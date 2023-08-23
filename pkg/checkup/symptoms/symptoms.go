package symptoms

type Symptom struct {
	Severity     string `validate:"oneof=warning critical"`
	ResourceName string
	ResourceType string
	Namespace    string
	Message      string
}

type SymptomList struct {
	Symptoms []Symptom
}

func (l *SymptomList) Add(s Symptom) {
	l.Symptoms = append(l.Symptoms, s)
}

type ContainerSymptom struct {
	Name     string
	Severity string `validate:"oneof=warning critical"`
	Message  string
}

type ContainerSymptomList struct {
	Symptoms []ContainerSymptom
}

func (l *ContainerSymptomList) Add(s ContainerSymptom) {
	l.Symptoms = append(l.Symptoms, s)
}
