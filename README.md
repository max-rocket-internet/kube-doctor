# 🩺 kube-doctor

[![](img/k8s-logo-sick.png)](#)

[![license](https://img.shields.io/github/license/sebasrp/awslimitchecker)](https://tldrlegal.com/license/mit-license)
[![go Report Card](https://goreportcard.com/badge/github.com/max-rocket-internet/kube-doctor)](https://goreportcard.com/report/github.com/max-rocket-internet/kube-doctor)

Is your Kubernetes cluster unhealthy? Do your workloads have symptoms? Then maybe it needs a checkup with `kube-doctor` 🏥

```shell
$ kube-doctor --warning-symptoms --non-namespaced-resources
== Checking DaemonSet resources
👀 DaemonSet kube-system/efs-csi-node: efs-plugin no resources specified
== Checking Deployment resources
👀 Deployment opencost/opencost: container 'opencost' memory request and limit are not equal
👀 Deployment default/application-one-listener: 5/8 pods are not ready
== Checking Endpoint resources
❌ Endpoint default/application-two: no ready addresses in subsets
== Checking Event resources
❌ Event datadog/datadog-x62q2: (Pod) 43.4 minutes ago: network is not ready: container runtime network not ready: NetworkReady=fals...
❌ Event default/application-one-597f47458c-fdb4r: (Pod) 1.6 minutes ago: Back-off restarting failed container
❌ Event datadog/datadog-95q6n: (Pod) 18.6 minutes ago: deleting pod for node scale down
❌ Event ip-10-10-10-10.compute.internal: (Node) 9.5 minutes ago: marked the node as toBeDeleted/unschedulable
❌ Event kube-system/cluster-autoscaler-status: (ConfigMap) 26.6 minutes ago: Scale-down: node ip-10-10-10-20.compute.internal removed with drain
== Checking HorizontalPodAutoscaler resources
👀 HorizontalPodAutoscaler default/application-three: has condition ScalingActive=False and reason ScalingDisabled
👀 HorizontalPodAutoscaler default/application-four: has condition ScalingLimited=True and reason TooFewReplicas
== Checking Job resources
❌ Job production/train-model: BackoffLimitExceeded: Job has reached the specified backoff limit
❌ Job production/run-analysis: DeadlineExceeded: Job was active longer than specified deadline
== Checking PersistentVolume resources
❌ PersistentVolume pgwatch-storage-pv-database: older than 5 minutes and status is not bound
== Checking Pod resources
❌ Pod default/application-two-uje-h2bhq: not running
❌ Pod datadog/datadog-555h5: status condition Ready is False
❌ Pod default/application-six: container 'app' was restarted 3.1 mins ago: 1 (exit code) Error (reason)
👀 Pod default/application-two-lhu-4r7hn: container 'app' has been restarted 5 times
== Checking Service resources
🎉 No symptoms found
== Checking PersistentVolume resources
⭕️ No resources found
== Checking KubeApiHealthEndpointStatus resources
🎉 No symptoms found
```

This tool will check for the following symptoms:

- `Container`:
  - no resources specified
  - no memory resources specified
  - no memory limit
  - memory request and limit are not equal
- `DaemonSet`:
  - pods are not ready
  - pods are miss-scheduled
  - rolling update in progress
- `Deployment`:
  - minimum availability not met
  - ReplicaSet update in progress
  - ReplicaSet update in progress but no progress
- `Endpoint`:
  - no ready addresses in subsets
- `Event`:
  - `cluster-autoscaler` events that are:
    - events that are not `Type=Normal`
    - `ScaleUp` and `ScaleDown` events
  - `service-controller` events that are not `Type=Normal`
  - `default-scheduler` events that are not `Type=Normal` and not `Reason=FailedScheduling`
  - `kubelet` events that are not `Type=Normal` and not `Reason=Unhealthy`
- `HorizontalPodAutoscaler`:
  - various bad status conditions
- `Job`:
  - `Failed` jobs within last hour
- [Kubernetes API health endpoints](https://kubernetes.io/docs/reference/using-api/health-checks/):
  - any bad or unknown checks
- `Node`:
  - mixed kubelet versions
  - not `Ready`
  - older than 5 minutes and not `Ready`
  - any bad status conditions
- `PersistentVolumeClaim` & `PersistentVolume`
  - older than 5 minutes and not in `Bound` phase
- `Pod`:
  - phase that is not `Running`
  - any bad pod status conditions
  - various bad container status conditions
  - crashed containers in last hour
  - without owner (created from `kubectl run`)
- `Service`:
  - `LoadBalancer` type without bad status

## Running

By default `kube-doctor` will check all namespaces but it can also target a specific namespace:

```shell
kube-doctor --namespace kube-system
```

Or label selector;:

```shell
kube-doctor --label-selector app.kubernetes.io/name=prometheus
```

Or a combination of both:

```shell
kube-doctor --label-selector app.kubernetes.io/name=prometheus --namespace monitoring
```

Non-namespaced resources are checked separately and can be enalbed with the `--non-namespaced-resources` flag:

```shell
kube-doctor --non-namespaced-resources
```

To see other options, including debug logging, consult the help:

```shell
kube-doctor --help
```

## Installation

Check out code and build:

```shell
git clone git@github.com:max-rocket-internet/kube-doctor.git
cd kube-doctor
go build ./... && go install ./...
```

Run from `main` branch without `git`:

```shell
go install github.com/max-rocket-internet/kube-doctor@latest
cd $GOPATH/pkg/mod/github.com/max-rocket-internet/kube-doctor*/
go run main.go
```

To get a binary, check [the releases](https://github.com/max-rocket-internet/kube-doctor/releases).

## Contributing

Pull requests welcome 💙

To run all tests:

```shell
go test ./...
```

Or just a single package:

```shell
go test ./.../checkup
```
