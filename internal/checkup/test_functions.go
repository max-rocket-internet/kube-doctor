package checkup

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func int32Ptr(i int32) *int32 { return &i }

func int32Literal(i int32) int32 { return i }

func createNewTimeStampPtr(t time.Time) *metav1.Time {
	i := metav1.NewTime(t)
	return &i
}
