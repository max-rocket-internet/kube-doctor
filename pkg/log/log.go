package log

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/max-rocket-internet/kube-doctor/pkg/checkup/symptoms"
)

var (
	debugLogger            *log.Logger
	infoLogger             *log.Logger
	logWarningSymptoms     bool
	logDebug               bool
	colorResourceTypeBold  = color.New(color.FgMagenta, color.Bold).SprintFunc()
	colorResourceType      = color.New(color.FgMagenta).SprintFunc()
	colorResourceName      = color.New(color.FgBlue, color.Bold).SprintFunc()
	colorResourceNamespace = color.New(color.FgGreen).SprintFunc()
	colorDebug             = color.New(color.Faint).SprintFunc()
	messageCharacterLimit  = 120
)

func init() {
	debugLogger = log.New(os.Stdout, "", 0)
	infoLogger = log.New(os.Stdout, "", 0)
}

func Setup(debugEnabled bool, warningSymptoms bool) {
	if debugEnabled {
		logDebug = true
	}

	if warningSymptoms {
		logWarningSymptoms = true
	}
}

func Debug(message string) {
	if !logDebug {
		return
	}
	debugLogger.Println(colorDebug(fmt.Sprintf("DEBUG: %s", message)))
}

func Info(message string) {
	infoLogger.Println(message)
}

func Error(message string, e error) {
	infoLogger.Printf("❗ %s: %s\n", message, e)
}

func Fatal(message string, e error) {
	infoLogger.Fatalf("💣 %s: %s\n", message, e)
}

func trimMessage(message string) string {
	if len(message) > messageCharacterLimit {
		return fmt.Sprintf("%s...", string([]byte(message)[:messageCharacterLimit]))
	} else {
		return message
	}
}

func LogSymptoms(s symptoms.SymptomList) {
	for _, s := range s.Symptoms {
		var message string

		if s.ResourceName == "" {
			message = fmt.Sprintf("%s: %s", colorResourceType(s.ResourceType), trimMessage(s.Message))
		} else if s.Namespace == "" {
			message = fmt.Sprintf("%s %s: %s", colorResourceType(s.ResourceType), s.ResourceName, trimMessage(s.Message))
		} else {
			message = fmt.Sprintf("%s %s/%s: %s", colorResourceType(s.ResourceType), colorResourceNamespace(s.Namespace), colorResourceName(s.ResourceName), trimMessage(s.Message))
		}

		if s.Severity == "warning" {
			if logWarningSymptoms {
				Info(fmt.Sprintf("👀 %s", message))
			}
		} else if s.Severity == "critical" {
			Info(fmt.Sprintf("❌ %s", message))
		} else {
			Error(fmt.Sprintf("unknown symptom severity: %s", s.Severity), nil)
		}
	}
}

func PrintBegin(resourceCount int, resourceType string) {
	Info(fmt.Sprintf("== Checking %d %s resources", resourceCount, colorResourceTypeBold(resourceType)))
}

func PrintEnd(resourceCount int, symptomCounts [2]int) {
	criticalCount := symptomCounts[0]
	warningCount := symptomCounts[1]
	totalSymptomCount := criticalCount + warningCount

	if resourceCount == 0 {
		Info("⭕️ No resources found")
		return
	}

	if totalSymptomCount == 0 {
		Info("🎉 No symptoms found")
		return
	}

	if (warningCount == 0 && criticalCount > 0) || logWarningSymptoms {
		return
	}

	if criticalCount == 0 && warningCount > 0 && !logWarningSymptoms {
		Info(fmt.Sprintf("👀 No critcal symptoms found but %d warning symptoms", warningCount))
	}

}
