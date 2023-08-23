package doctor

import (
	"fmt"

	"github.com/max-rocket-internet/kube-doctor/pkg/checkup"
	"github.com/max-rocket-internet/kube-doctor/pkg/kubernetes"
	"github.com/max-rocket-internet/kube-doctor/pkg/log"
	"github.com/urfave/cli/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DoCheckUp(cCtx *cli.Context) {
	kubernetes.Init()

	log.Setup(cCtx.Bool("debug"), cCtx.Bool("warning-symptoms"))
	log.Debug(fmt.Sprintf("Connected to cluster from context %s running version %s", kubernetes.ContextName, kubernetes.ServerVersion))

	checkNonNamespaced := cCtx.Bool("non-namespaced-resources")
	namespace := cCtx.String("namespace")
	labelSelector := cCtx.String("label-selector")

	if checkNonNamespaced {
		log.LogSymptoms(checkup.CheckNodes(kubernetes.GetNodes()))
		log.LogSymptoms(checkup.CheckPersistentVolumes(kubernetes.GetPersistentVolumes()))
		log.LogSymptoms(checkup.KubeApiHealthStatuses(kubernetes.GetKubeApiHealth()))
	}

	log.LogSymptoms(checkup.CheckDaemonSets(kubernetes.GetDaemonSets(namespace, metav1.ListOptions{LabelSelector: labelSelector})))
	log.LogSymptoms(checkup.CheckDeployments(kubernetes.GetDeployments(namespace, metav1.ListOptions{LabelSelector: labelSelector})))
	log.LogSymptoms(checkup.CheckEndpoints(kubernetes.GetEndpoints(namespace, metav1.ListOptions{LabelSelector: labelSelector})))
	log.LogSymptoms(checkup.CheckEvents(kubernetes.GetEvents(namespace, metav1.ListOptions{LabelSelector: labelSelector})))
	log.LogSymptoms(checkup.CheckHpas(kubernetes.GetHpas(namespace, metav1.ListOptions{LabelSelector: labelSelector})))
	log.LogSymptoms(checkup.CheckJobs(kubernetes.GetJobs(namespace, metav1.ListOptions{LabelSelector: labelSelector})))
	log.LogSymptoms(checkup.CheckPersistentVolumeClaims(kubernetes.GetPersistentVolumeClaims(namespace, metav1.ListOptions{LabelSelector: labelSelector})))
	log.LogSymptoms(checkup.CheckPods(kubernetes.GetPods(namespace, metav1.ListOptions{LabelSelector: labelSelector})))
	log.LogSymptoms(checkup.CheckServices(kubernetes.GetServices(namespace, metav1.ListOptions{LabelSelector: labelSelector})))
}
