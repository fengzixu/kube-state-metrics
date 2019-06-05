/*
Copyright 2018 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collector

import (
	"sort"
	"strings"

	"k8s.io/klog"

	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v2beta1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	certv1beta1 "k8s.io/api/certificates/v1beta1"
	v1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	policy "k8s.io/api/policy/v1beta1"
	storagev1 "k8s.io/api/storage/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kube-state-metrics/pkg/metric"
	metricsstore "k8s.io/kube-state-metrics/pkg/metrics_store"
	"k8s.io/kube-state-metrics/pkg/options"
)

type whiteBlackLister interface {
	IsIncluded(string) bool
	IsExcluded(string) bool
}

// Builder helps to build collector. It follows the builder pattern
// (https://en.wikipedia.org/wiki/Builder_pattern).
type Builder struct {
	kubeClient        clientset.Interface
	namespaces        options.NamespaceList
	ctx               context.Context
	enabledCollectors []string
	whiteBlackList    whiteBlackLister
}

// NewBuilder returns a new builder.
func NewBuilder(
	ctx context.Context,
) *Builder {
	return &Builder{
		ctx: ctx,
	}
}

// WithEnabledCollectors sets the enabledCollectors property of a Builder.
func (b *Builder) WithEnabledCollectors(c []string) {
	var copy []string
	copy = append(copy, c...)

	sort.Strings(copy)

	b.enabledCollectors = copy
}

// WithNamespaces sets the namespaces property of a Builder.
func (b *Builder) WithNamespaces(n options.NamespaceList) {
	b.namespaces = n
}

// WithKubeClient sets the kubeClient property of a Builder.
func (b *Builder) WithKubeClient(c clientset.Interface) {
	b.kubeClient = c
}

// WithWhiteBlackList configures the white or blacklisted metric to be exposed
// by the collector build by the Builder
func (b *Builder) WithWhiteBlackList(l whiteBlackLister) {
	b.whiteBlackList = l
}

// Build initializes and registers all enabled collectors.
func (b *Builder) Build() []*metricsstore.MetricsStore {
	if b.whiteBlackList == nil {
		panic("whiteBlackList should not be nil")
	}

	collectors := []*metricsstore.MetricsStore{}
	activeCollectorNames := []string{}

	for _, c := range b.enabledCollectors {
		constructor, ok := availableCollectors[c]
		if ok {
			collector := constructor(b)
			activeCollectorNames = append(activeCollectorNames, c)
			collectors = append(collectors, collector)
		}
	}

	klog.Infof("Active collectors: %s", strings.Join(activeCollectorNames, ","))

	return collectors
}

<<<<<<< HEAD
var availableCollectors = map[string]func(f *Builder) *coll.Collector{
	"certificatesigningrequests": func(b *Builder) *coll.Collector { return b.buildCsrCollector() },
	"configmaps":                 func(b *Builder) *coll.Collector { return b.buildConfigMapCollector() },
	"cronjobs":                   func(b *Builder) *coll.Collector { return b.buildCronJobCollector() },
	"daemonsets":                 func(b *Builder) *coll.Collector { return b.buildDaemonSetCollector() },
	"deployments":                func(b *Builder) *coll.Collector { return b.buildDeploymentCollector() },
	"endpoints":                  func(b *Builder) *coll.Collector { return b.buildEndpointsCollector() },
	"horizontalpodautoscalers":   func(b *Builder) *coll.Collector { return b.buildHPACollector() },
	"ingresses":                  func(b *Builder) *coll.Collector { return b.buildIngressCollector() },
	"jobs":                       func(b *Builder) *coll.Collector { return b.buildJobCollector() },
	"limitranges":                func(b *Builder) *coll.Collector { return b.buildLimitRangeCollector() },
	"namespaces":                 func(b *Builder) *coll.Collector { return b.buildNamespaceCollector() },
	"nodes":                      func(b *Builder) *coll.Collector { return b.buildNodeCollector() },
	"persistentvolumeclaims":     func(b *Builder) *coll.Collector { return b.buildPersistentVolumeClaimCollector() },
	"persistentvolumes":          func(b *Builder) *coll.Collector { return b.buildPersistentVolumeCollector() },
	"poddisruptionbudgets":       func(b *Builder) *coll.Collector { return b.buildPodDisruptionBudgetCollector() },
	"pods":                       func(b *Builder) *coll.Collector { return b.buildPodCollector() },
	"replicasets":                func(b *Builder) *coll.Collector { return b.buildReplicaSetCollector() },
	"replicationcontrollers":     func(b *Builder) *coll.Collector { return b.buildReplicationControllerCollector() },
	"resourcequotas":             func(b *Builder) *coll.Collector { return b.buildResourceQuotaCollector() },
	"secrets":                    func(b *Builder) *coll.Collector { return b.buildSecretCollector() },
	"services":                   func(b *Builder) *coll.Collector { return b.buildServiceCollector() },
	"statefulsets":               func(b *Builder) *coll.Collector { return b.buildStatefulSetCollector() },
	"storageclasses":             func(b *Builder) *coll.Collector { return b.buildStorageClassCollector() },
}

func (b *Builder) buildConfigMapCollector() *coll.Collector {
=======
var availableCollectors = map[string]func(f *Builder) *metricsstore.MetricsStore{
	"certificatesigningrequests": func(b *Builder) *metricsstore.MetricsStore { return b.buildCsrCollector() },
	"configmaps":                 func(b *Builder) *metricsstore.MetricsStore { return b.buildConfigMapCollector() },
	"cronjobs":                   func(b *Builder) *metricsstore.MetricsStore { return b.buildCronJobCollector() },
	"daemonsets":                 func(b *Builder) *metricsstore.MetricsStore { return b.buildDaemonSetCollector() },
	"deployments":                func(b *Builder) *metricsstore.MetricsStore { return b.buildDeploymentCollector() },
	"endpoints":                  func(b *Builder) *metricsstore.MetricsStore { return b.buildEndpointsCollector() },
	"horizontalpodautoscalers":   func(b *Builder) *metricsstore.MetricsStore { return b.buildHPACollector() },
	"ingresses":                  func(b *Builder) *metricsstore.MetricsStore { return b.buildIngressCollector() },
	"jobs":                       func(b *Builder) *metricsstore.MetricsStore { return b.buildJobCollector() },
	"limitranges":                func(b *Builder) *metricsstore.MetricsStore { return b.buildLimitRangeCollector() },
	"namespaces":                 func(b *Builder) *metricsstore.MetricsStore { return b.buildNamespaceCollector() },
	"nodes":                      func(b *Builder) *metricsstore.MetricsStore { return b.buildNodeCollector() },
	"persistentvolumeclaims":     func(b *Builder) *metricsstore.MetricsStore { return b.buildPersistentVolumeClaimCollector() },
	"persistentvolumes":          func(b *Builder) *metricsstore.MetricsStore { return b.buildPersistentVolumeCollector() },
	"poddisruptionbudgets":       func(b *Builder) *metricsstore.MetricsStore { return b.buildPodDisruptionBudgetCollector() },
	"pods":                       func(b *Builder) *metricsstore.MetricsStore { return b.buildPodCollector() },
	"replicasets":                func(b *Builder) *metricsstore.MetricsStore { return b.buildReplicaSetCollector() },
	"replicationcontrollers":     func(b *Builder) *metricsstore.MetricsStore { return b.buildReplicationControllerCollector() },
	"resourcequotas":             func(b *Builder) *metricsstore.MetricsStore { return b.buildResourceQuotaCollector() },
	"secrets":                    func(b *Builder) *metricsstore.MetricsStore { return b.buildSecretCollector() },
	"services":                   func(b *Builder) *metricsstore.MetricsStore { return b.buildServiceCollector() },
	"statefulsets":               func(b *Builder) *metricsstore.MetricsStore { return b.buildStatefulSetCollector() },
}

func (b *Builder) buildConfigMapCollector() *metricsstore.MetricsStore {
	return b.buildCollector(configMapMetricFamilies, &v1.ConfigMap{}, createConfigMapListWatch)
}

func (b *Builder) buildCronJobCollector() *metricsstore.MetricsStore {
	return b.buildCollector(cronJobMetricFamilies, &batchv1beta1.CronJob{}, createCronJobListWatch)
}

func (b *Builder) buildDaemonSetCollector() *metricsstore.MetricsStore {
	return b.buildCollector(daemonSetMetricFamilies, &appsv1.DaemonSet{}, createDaemonSetListWatch)
}

func (b *Builder) buildDeploymentCollector() *metricsstore.MetricsStore {
	return b.buildCollector(deploymentMetricFamilies, &appsv1.Deployment{}, createDeploymentListWatch)
}

func (b *Builder) buildEndpointsCollector() *metricsstore.MetricsStore {
	return b.buildCollector(endpointMetricFamilies, &v1.Endpoints{}, createEndpointsListWatch)
}

func (b *Builder) buildHPACollector() *metricsstore.MetricsStore {
	return b.buildCollector(hpaMetricFamilies, &autoscaling.HorizontalPodAutoscaler{}, createHPAListWatch)
}

func (b *Builder) buildIngressCollector() *metricsstore.MetricsStore {
	return b.buildCollector(ingressMetricFamilies, &extensions.Ingress{}, createIngressListWatch)
}

func (b *Builder) buildJobCollector() *metricsstore.MetricsStore {
	return b.buildCollector(jobMetricFamilies, &batchv1.Job{}, createJobListWatch)
}

func (b *Builder) buildLimitRangeCollector() *metricsstore.MetricsStore {
	return b.buildCollector(limitRangeMetricFamilies, &v1.LimitRange{}, createLimitRangeListWatch)
}

func (b *Builder) buildNamespaceCollector() *metricsstore.MetricsStore {
	return b.buildCollector(namespaceMetricFamilies, &v1.Namespace{}, createNamespaceListWatch)
}

func (b *Builder) buildNodeCollector() *metricsstore.MetricsStore {
	return b.buildCollector(nodeMetricFamilies, &v1.Node{}, createNodeListWatch)
}

func (b *Builder) buildPersistentVolumeClaimCollector() *metricsstore.MetricsStore {
	return b.buildCollector(persistentVolumeClaimMetricFamilies, &v1.PersistentVolumeClaim{}, createPersistentVolumeClaimListWatch)
}

func (b *Builder) buildPersistentVolumeCollector() *metricsstore.MetricsStore {
	return b.buildCollector(persistentVolumeMetricFamilies, &v1.PersistentVolume{}, createPersistentVolumeListWatch)
}

func (b *Builder) buildPodDisruptionBudgetCollector() *metricsstore.MetricsStore {
	return b.buildCollector(podDisruptionBudgetMetricFamilies, &policy.PodDisruptionBudget{}, createPodDisruptionBudgetListWatch)
}

func (b *Builder) buildReplicaSetCollector() *metricsstore.MetricsStore {
	return b.buildCollector(replicaSetMetricFamilies, &extensions.ReplicaSet{}, createReplicaSetListWatch)
}

func (b *Builder) buildReplicationControllerCollector() *metricsstore.MetricsStore {
	return b.buildCollector(replicationControllerMetricFamilies, &v1.ReplicationController{}, createReplicationControllerListWatch)
}

func (b *Builder) buildResourceQuotaCollector() *metricsstore.MetricsStore {
	return b.buildCollector(resourceQuotaMetricFamilies, &v1.ResourceQuota{}, createResourceQuotaListWatch)
}

func (b *Builder) buildSecretCollector() *metricsstore.MetricsStore {
	return b.buildCollector(secretMetricFamilies, &v1.Secret{}, createSecretListWatch)
}

func (b *Builder) buildServiceCollector() *metricsstore.MetricsStore {
	return b.buildCollector(serviceMetricFamilies, &v1.Service{}, createServiceListWatch)
}

func (b *Builder) buildStatefulSetCollector() *metricsstore.MetricsStore {
	return b.buildCollector(statefulSetMetricFamilies, &appsv1.StatefulSet{}, createStatefulSetListWatch)
}

<<<<<<< HEAD
func (b *Builder) buildStorageClassCollector() *coll.Collector {
	return b.buildCollector(storageClassMetricFamilies, &storagev1.StorageClass{}, createStorageClassListWatch)
}

func (b *Builder) buildPodCollector() *metricsstore.MetricsStore {
	return b.buildCollector(podMetricFamilies, &v1.Pod{}, createPodListWatch)
}

func (b *Builder) buildCsrCollector() *metricsstore.MetricsStore {
	return b.buildCollector(csrMetricFamilies, &certv1beta1.CertificateSigningRequest{}, createCSRListWatch)
}

func (b *Builder) buildCollector(
	metricFamilies []metric.FamilyGenerator,
	expectedType interface{},
	listWatchFunc func(kubeClient clientset.Interface, ns string) cache.ListerWatcher,
) *metricsstore.MetricsStore {
	filteredMetricFamilies := metric.FilterMetricFamilies(b.whiteBlackList, metricFamilies)
	composedMetricGenFuncs := metric.ComposeMetricGenFuncs(filteredMetricFamilies)

	familyHeaders := metric.ExtractMetricFamilyHeaders(filteredMetricFamilies)

	store := metricsstore.NewMetricsStore(
		familyHeaders,
		composedMetricGenFuncs,
	)
	b.reflectorPerNamespace(expectedType, store, listWatchFunc)

	return store
}

// reflectorPerNamespace creates a Kubernetes client-go reflector with the given
// listWatchFunc for each given namespace and registers it with the given store.
func (b *Builder) reflectorPerNamespace(
	expectedType interface{},
	store cache.Store,
	listWatchFunc func(kubeClient clientset.Interface, ns string) cache.ListerWatcher,
) {
	for _, ns := range b.namespaces {
		lw := listWatchFunc(b.kubeClient, ns)
		reflector := cache.NewReflector(lw, expectedType, store, 0)
		go reflector.Run(b.ctx.Done())
	}
}
