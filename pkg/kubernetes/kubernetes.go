package kubernetes

import (
	"context"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v2"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/max-rocket-internet/kube-doctor/pkg/kubernetes/statuses"
	"github.com/max-rocket-internet/kube-doctor/pkg/log"
)

var (
	client        *kubernetes.Clientset
	ContextName   string
	ServerVersion string
)

func Init() {
	client = createClient()
}

func createClient() *kubernetes.Clientset {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Fatal("error creating kubernetes client config", err)
	}

	rawConfig, _ := kubeConfig.RawConfig()
	ContextName = rawConfig.CurrentContext

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("error creating kubernetes client", err)
	}

	ver, err := client.DiscoveryClient.ServerVersion()
	if err != nil {
		log.Fatal("error getting kubernetes version", err)
	}
	ServerVersion = ver.String()

	return client
}

func GetPods(namespace string, listOptions metav1.ListOptions) *v1.PodList {
	resources, err := client.CoreV1().Pods(namespace).List(context.TODO(), listOptions)

	if err != nil {
		log.Error("Error getting Pods", err)
	}

	return resources
}

func GetHpas(namespace string, listOptions metav1.ListOptions) *autoscaling.HorizontalPodAutoscalerList {
	resources, err := client.AutoscalingV2().HorizontalPodAutoscalers(namespace).List(context.TODO(), listOptions)

	if err != nil {
		log.Error("Error getting HorizontalPodAutoscalers", err)
	}

	return resources
}

func GetServices(namespace string, listOptions metav1.ListOptions) *v1.ServiceList {
	resources, err := client.CoreV1().Services(namespace).List(context.TODO(), listOptions)

	if err != nil {
		log.Error("Error getting Services: %s", err)
	}

	return resources
}

func GetPersistentVolumeClaims(namespace string, listOptions metav1.ListOptions) *v1.PersistentVolumeClaimList {
	resources, err := client.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), listOptions)

	if err != nil {
		log.Error("Error getting PersistentVolumeClaims", err)
	}

	return resources
}

func GetEndpoints(namespace string, listOptions metav1.ListOptions) *v1.EndpointsList {
	resources, err := client.CoreV1().Endpoints(namespace).List(context.TODO(), listOptions)

	if err != nil {
		log.Error("Error getting Endpoints", err)
	}

	return resources
}

func GetDaemonSets(namespace string, listOptions metav1.ListOptions) *appsv1.DaemonSetList {
	resources, err := client.AppsV1().DaemonSets(namespace).List(context.TODO(), listOptions)

	if err != nil {
		log.Error("Error getting DaemonSets", err)
	}

	return resources
}

func GetDaemonSetByName(namespace string, name string) *appsv1.DaemonSetList {
	resources, err := client.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		log.Error("Error getting DaemonSets", err)
	}

	return &appsv1.DaemonSetList{Items: []appsv1.DaemonSet{*resources}}
}

func GetDeployments(namespace string, listOptions metav1.ListOptions) *appsv1.DeploymentList {
	resources, err := client.AppsV1().Deployments(namespace).List(context.TODO(), listOptions)

	if err != nil {
		log.Error("Error getting Deployments", err)
	}

	return resources
}

func GetDeploymentByName(namespace string, name string) *appsv1.DeploymentList {
	resources, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		log.Error("Error getting Deployments", err)
	}

	return &appsv1.DeploymentList{Items: []appsv1.Deployment{*resources}}
}

func GetJobs(namespace string, listOptions metav1.ListOptions) *batchv1.JobList {
	resources, err := client.BatchV1().Jobs(namespace).List(context.TODO(), listOptions)

	if err != nil {
		log.Error("Error getting Jobs", err)
	}

	return resources
}

func GetEvents(namespace string, listOptions metav1.ListOptions) *v1.EventList {
	resources, err := client.CoreV1().Events(namespace).List(context.TODO(), listOptions)

	if err != nil {
		log.Error("Error getting Events", err)
	}

	return resources
}

func GetEventsNotNormal(namespace string) *v1.EventList {
	resources, err := client.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "type!=Normal"})

	if err != nil {
		log.Error("Error getting Events", err)
	}

	return resources
}

// Non-namespaced resources

func GetPersistentVolumes() *v1.PersistentVolumeList {
	resources, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Error("Error getting PersistentVolumes", err)
	}

	return resources
}

func GetKubeApiHealth() *statuses.KubeApiHealthEndpointStatusList {
	results := statuses.KubeApiHealthEndpointStatusList{}

	for _, path := range []string{"/readyz", "/livez"} {
		var requestStatusCode int
		r := client.RESTClient().Get().AbsPath(path).Param("verbose", "").Do(context.TODO())

		r.StatusCode(&requestStatusCode)
		if requestStatusCode != 200 {
			log.Error(fmt.Sprintf("Error getting %s with RESTClient", path), nil)
			continue
		}

		requestBodyRaw, err := r.Raw()
		if err != nil {
			log.Error(fmt.Sprintf("Error getting raw body from %s request", path), err)
			continue
		}

		requestBody := string(requestBodyRaw[:])
		for _, line := range strings.Split(strings.TrimSuffix(requestBody, "\n"), "\n") {
			results.AddRawLine(line, path)
		}
	}

	return &results
}

func GetNodes() *v1.NodeList {
	resources, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Error("Error getting Nodes", err)
	}

	return resources
}

func GetComponentStatuses() *v1.ComponentStatusList {
	resources, err := client.CoreV1().ComponentStatuses().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Error("Error getting ComponentStatuses", err)
	}

	return resources
}
