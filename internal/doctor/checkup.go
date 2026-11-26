// Package doctor runs every checkup against a cluster and prints the results.
package doctor

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/max-rocket-internet/kube-doctor/internal/checkup"
	"github.com/max-rocket-internet/kube-doctor/internal/checkup/symptoms"
	"github.com/max-rocket-internet/kube-doctor/internal/kubernetes"
	"github.com/max-rocket-internet/kube-doctor/internal/log"
)

// Options configures a checkup run.
type Options struct {
	// Namespace limits namespaced checks to a single namespace. Empty means all namespaces.
	Namespace string
	// LabelSelector filters namespaced resources, e.g. "app.kubernetes.io/name=prometheus".
	LabelSelector string
	// NonNamespaced enables checks of cluster-scoped resources and API server health.
	NonNamespaced bool
	Debug         bool
	// WarningSymptoms prints warning symptoms in addition to critical ones.
	WarningSymptoms bool
}

// Run connects to the cluster from the current kubeconfig context and runs every checkup.
// Checks whose resources cannot be listed are reported and skipped; Run returns an error
// if the client cannot be created or any check could not be run.
func Run(ctx context.Context, opts Options) error {
	log.Setup(opts.Debug, opts.WarningSymptoms)

	client, err := kubernetes.NewClient()
	if err != nil {
		return err
	}
	log.Debugf("Connected to cluster from context %s running version %s", client.ContextName, client.ServerVersion)

	r := runner{ctx: ctx}

	if opts.NonNamespaced {
		run(&r, "Nodes", func(ctx context.Context) (*corev1.NodeList, error) {
			return client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		}, checkup.CheckNodes)
		run(&r, "PersistentVolumes", func(ctx context.Context) (*corev1.PersistentVolumeList, error) {
			return client.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
		}, checkup.CheckPersistentVolumes)
		run(&r, "Kubernetes API health", client.KubeAPIHealth, checkup.CheckKubeAPIHealth)
	}

	ns, listOpts := opts.Namespace, metav1.ListOptions{LabelSelector: opts.LabelSelector}

	run(&r, "DaemonSets", func(ctx context.Context) (*appsv1.DaemonSetList, error) {
		return client.AppsV1().DaemonSets(ns).List(ctx, listOpts)
	}, checkup.CheckDaemonSets)
	run(&r, "Deployments", func(ctx context.Context) (*appsv1.DeploymentList, error) {
		return client.AppsV1().Deployments(ns).List(ctx, listOpts)
	}, checkup.CheckDeployments)
	run(&r, "EndpointSlices", func(ctx context.Context) (*discoveryv1.EndpointSliceList, error) {
		return client.DiscoveryV1().EndpointSlices(ns).List(ctx, listOpts)
	}, checkup.CheckEndpointSlices)
	run(&r, "Events", func(ctx context.Context) (*corev1.EventList, error) {
		return client.CoreV1().Events(ns).List(ctx, listOpts)
	}, checkup.CheckEvents)
	run(&r, "HorizontalPodAutoscalers", func(ctx context.Context) (*autoscalingv2.HorizontalPodAutoscalerList, error) {
		return client.AutoscalingV2().HorizontalPodAutoscalers(ns).List(ctx, listOpts)
	}, checkup.CheckHPAs)
	run(&r, "Jobs", func(ctx context.Context) (*batchv1.JobList, error) {
		return client.BatchV1().Jobs(ns).List(ctx, listOpts)
	}, checkup.CheckJobs)
	run(&r, "PersistentVolumeClaims", func(ctx context.Context) (*corev1.PersistentVolumeClaimList, error) {
		return client.CoreV1().PersistentVolumeClaims(ns).List(ctx, listOpts)
	}, checkup.CheckPersistentVolumeClaims)
	run(&r, "Pods", func(ctx context.Context) (*corev1.PodList, error) {
		return client.CoreV1().Pods(ns).List(ctx, listOpts)
	}, checkup.CheckPods)
	run(&r, "Services", func(ctx context.Context) (*corev1.ServiceList, error) {
		return client.CoreV1().Services(ns).List(ctx, listOpts)
	}, checkup.CheckServices)

	if r.failed > 0 {
		return fmt.Errorf("%d check(s) could not be completed", r.failed)
	}
	return nil
}

type runner struct {
	ctx    context.Context
	failed int
}

// run fetches resources with list and prints the symptoms found by check.
// It is a function rather than a method because methods cannot have type parameters.
func run[T any](r *runner, name string, list func(context.Context) (T, error), check func(T) symptoms.SymptomList) {
	if r.ctx.Err() != nil {
		return
	}

	resources, err := list(r.ctx)
	if err != nil {
		r.failed++
		log.Errorf("Error getting %s: %v", name, err)
		return
	}

	log.LogSymptoms(check(resources))
}
