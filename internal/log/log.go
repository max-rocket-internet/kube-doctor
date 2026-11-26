// Package log prints checkup progress and symptoms to the terminal.
package log

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
)

const messageCharacterLimit = 120

var (
	out    io.Writer = os.Stdout
	errOut io.Writer = os.Stderr

	logWarningSymptoms bool
	logDebug           bool

	colorResourceTypeBold  = color.New(color.FgMagenta, color.Bold).SprintFunc()
	colorResourceType      = color.New(color.FgMagenta).SprintFunc()
	colorResourceName      = color.New(color.FgBlue, color.Bold).SprintFunc()
	colorResourceNamespace = color.New(color.FgGreen).SprintFunc()
	colorDebug             = color.New(color.Faint).SprintFunc()
)

// Setup configures which messages are printed.
func Setup(debugEnabled, warningSymptoms bool) {
	logDebug = debugEnabled
	logWarningSymptoms = warningSymptoms
}

// Debugf prints a debug message when debug logging is enabled.
func Debugf(format string, args ...any) {
	if !logDebug {
		return
	}
	_, _ = fmt.Fprintln(out, colorDebug("DEBUG: "+fmt.Sprintf(format, args...)))
}

// Info prints a message to stdout.
func Info(message string) {
	_, _ = fmt.Fprintln(out, message)
}

// Infof prints a formatted message to stdout.
func Infof(format string, args ...any) {
	Info(fmt.Sprintf(format, args...))
}

// Errorf prints a formatted error message to stderr.
func Errorf(format string, args ...any) {
	_, _ = fmt.Fprintf(errOut, "❗ %s\n", fmt.Sprintf(format, args...))
}

func trimMessage(message string) string {
	runes := []rune(message)
	if len(runes) <= messageCharacterLimit {
		return message
	}
	return string(runes[:messageCharacterLimit]) + "..."
}

// LogSymptoms prints every symptom in l, honouring the warning-symptoms setting.
func LogSymptoms(l symptoms.SymptomList) {
	for _, s := range l.Symptoms {
		var message string
		switch {
		case s.ResourceName == "":
			message = fmt.Sprintf("%s: %s", colorResourceType(s.ResourceType), trimMessage(s.Message))
		case s.Namespace == "":
			message = fmt.Sprintf("%s %s: %s", colorResourceType(s.ResourceType), s.ResourceName, trimMessage(s.Message))
		default:
			message = fmt.Sprintf("%s %s/%s: %s", colorResourceType(s.ResourceType), colorResourceNamespace(s.Namespace), colorResourceName(s.ResourceName), trimMessage(s.Message))
		}

		switch s.Severity {
		case symptoms.Warning:
			if logWarningSymptoms {
				Info("👀 " + message)
			}
		case symptoms.Critical:
			Info("❌ " + message)
		default:
			Errorf("unknown symptom severity: %q", s.Severity)
		}
	}
}

// PrintBegin prints the header for a resource type check.
func PrintBegin(resourceCount int, resourceType string) {
	Infof("== Checking %d %s resources", resourceCount, colorResourceTypeBold(resourceType))
}

// PrintEnd prints a summary line for a resource type check when appropriate.
func PrintEnd(resourceCount int, counts symptoms.Counts) {
	switch {
	case resourceCount == 0:
		Info("⭕️ No resources found")
	case counts.Total() == 0:
		Info("🎉 No symptoms found")
	case counts.Critical == 0 && !logWarningSymptoms:
		Infof("👀 No critical symptoms found but %d warning symptoms", counts.Warning)
	}
}
