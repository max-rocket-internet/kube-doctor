// Package kubernetes builds the Kubernetes client used by kube-doctor.
package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/max-rocket-internet/kube-doctor/internal/kubernetes/statuses"
)

// Client wraps a Kubernetes clientset along with details about the cluster it targets.
type Client struct {
	kubernetes.Interface

	ContextName   string
	ServerVersion string
}

// NewClient creates a client from the default kubeconfig loading rules
// ($KUBECONFIG, ~/.kube/config or in-cluster config) and verifies connectivity.
func NewClient() (*Client, error) {
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client: %w", err)
	}

	ver, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("getting kubernetes server version: %w", err)
	}

	// The raw config is only used for display; it is unavailable for in-cluster config.
	var contextName string
	if rawConfig, err := kubeConfig.RawConfig(); err == nil {
		contextName = rawConfig.CurrentContext
	}

	return &Client{
		Interface:     clientset,
		ContextName:   contextName,
		ServerVersion: ver.String(),
	}, nil
}

// KubeAPIHealth queries the verbose /readyz and /livez endpoints of the API server.
func (c *Client) KubeAPIHealth(ctx context.Context) (*statuses.KubeAPIHealthEndpointStatusList, error) {
	var (
		results statuses.KubeAPIHealthEndpointStatusList
		errs    []error
	)

	for _, path := range []string{"/readyz", "/livez"} {
		var code int
		body, err := c.Discovery().RESTClient().Get().AbsPath(path).Param("verbose", "").Do(ctx).StatusCode(&code).Raw()

		// An unhealthy endpoint responds with 500 and still lists every check in the body.
		if err != nil && (code != http.StatusInternalServerError || len(body) == 0) {
			errs = append(errs, fmt.Errorf("requesting %s: %w", path, err))
			continue
		}

		for line := range strings.Lines(string(body)) {
			results.AddRawLine(strings.TrimSpace(line), path)
		}
	}

	return &results, errors.Join(errs...)
}
