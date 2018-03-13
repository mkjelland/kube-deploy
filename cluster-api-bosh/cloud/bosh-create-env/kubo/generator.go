/*
Copyright 2018 The Kubernetes Authors.

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
package kubo

import (
	"encoding/json"
	"errors"
	"fmt"

	yaml "gopkg.in/yaml.v2"

	"regexp"

	"k8s.io/kube-deploy/cluster-api-bosh/cloud/bosh-create-env/director"
	"k8s.io/kube-deploy/cluster-api/api/cluster/v1alpha1"
	apiutil "k8s.io/kube-deploy/cluster-api/util"
)

func NewManifestGenerator() *ManifestGenerator {
	return &ManifestGenerator{}
}

type ManifestGenerator struct{}

const kubo_worker_1_9_2 = `
---
azs:
- z1
instances: 1
jobs:
- name: kubo-dns-aliases
  release: kubo-1.9.2
- name: secure-var-vcap
  release: kubo-1.9.2
- name: flanneld
  release: kubo-1.9.2
- name: docker
  properties:
    bip: 172.17.0.1/24
    default_ulimits:
    - nofile=65536
    env: {}
    flannel: true
    ip_masq: false
    iptables: false
    log_level: error
    storage_driver: overlay
    store_dir: /var/vcap/data
    tls_cacert: ((tls-docker.ca))
    tls_cert: ((tls-docker.certificate))
    tls_key: ((tls-docker.private_key))
  release: docker
- name: cloud-provider
  properties:
    cloud-provider:
      type: gce
  provides:
    cloud-provider:
      as: worker
  release: kubo-1.9.2
- name: kubelet
  properties:
    api-token: ((kubelet-password))
    tls:
      kubelet: ((tls-kubelet))
      kubernetes: ((tls-kubernetes))
  release: kubo-1.9.2
- name: kube-proxy
  properties:
    api-token: ((kube-proxy-password))
    tls:
      kubernetes: ((tls-kubernetes))
  release: kubo-1.9.2
name: worker
networks:
- name: default
persistent_disk: 10240
stemcell: trusty
vm_type: worker
`

const kubo_worker_1_8_6 = `
---
azs:
- z1
instances: 1
jobs:
- name: kubo-dns-aliases
  release: kubo-1.8.6
- name: secure-var-vcap
  release: kubo-1.8.6
- name: flanneld
  release: kubo-1.8.6
- name: docker
  properties:
    bip: 172.17.0.1/24
    default_ulimits:
    - nofile=65536
    env: {}
    flannel: true
    ip_masq: false
    iptables: false
    log_level: error
    storage_driver: overlay
    store_dir: /var/vcap/data
    tls_cacert: ((tls-docker.ca))
    tls_cert: ((tls-docker.certificate))
    tls_key: ((tls-docker.private_key))
  release: docker
- name: cloud-provider
  properties:
    cloud-provider:
      type: gce
  provides:
    cloud-provider:
      as: worker
  release: kubo-1.8.6
- name: kubelet
  properties:
    api-token: ((kubelet-password))
    tls:
      kubelet: ((tls-kubelet))
      kubernetes: ((tls-kubernetes))
  release: kubo-1.8.6
- name: kube-proxy
  properties:
    api-token: ((kube-proxy-password))
    tls:
      kubernetes: ((tls-kubernetes))
  release: kubo-1.8.6
name: worker
networks:
- name: default
persistent_disk: 10240
stemcell: trusty
vm_type: worker
`

const kubo_master_1_8_6 = `
---
azs:
  - z1
instances: 1
jobs:
- name: kubo-dns-aliases
  release: kubo-1.8.6
- name: docker
  release: docker
  properties:
    bip: 172.17.0.1/24
    default_ulimits:
    - nofile=65536
    env: {}
    flannel: true
    ip_masq: false
    iptables: false
    log_level: error
    storage_driver: overlay
    store_dir: "/var/vcap/data"
    tls_cacert: "((tls-docker.ca))"
    tls_cert: "((tls-docker.certificate))"
    tls_key: "((tls-docker.private_key))"
- name: kubelet
  release: kubo-1.8.6
  properties:
    api-token: "((kubelet-password))"
    register_with_taints: "node-role.kubernetes.io/master=true:NoSchedule"
    labels:
      type: "controller"
    tls:
      kubelet: "((tls-kubelet))"
      kubernetes: "((tls-kubernetes))"
- name: secure-var-vcap
  release: kubo-1.8.6
- name: cloud-provider
  properties:
    cloud-provider:
      type: gce
  provides:
    cloud-provider:
      as: master
  release: kubo-1.8.6
- name: flanneld
  release: kubo-1.8.6
- name: kube-apiserver
  properties:
    admin-password: ((kubo-admin-password))
    admin-username: admin
    authorization-mode: rbac
    backend_port: 8443
    kube-controller-manager-password: ((kube-controller-manager-password))
    kube-proxy-password: ((kube-proxy-password))
    kube-scheduler-password: ((kube-scheduler-password))
    kubelet-password: ((kubelet-password))
    port: 8443
    route-sync-password: ((route-sync-password))
    tls:
      kubernetes:
        ca: ((tls-kubernetes.ca))
        certificate: ((tls-kubernetes.certificate))
        private_key: ((tls-kubernetes.private_key))
  release: kubo-1.8.6
- name: kube-controller-manager
  properties:
    api-token: ((kube-controller-manager-password))
    tls:
      kubernetes: ((tls-kubernetes))
  release: kubo-1.8.6
- name: kube-scheduler
  properties:
    api-token: ((kube-scheduler-password))
    tls:
      kubernetes: ((tls-kubernetes))
  release: kubo-1.8.6
- name: etcd
  properties:
    etcd:
      advertise_urls_dns_suffix: etcd.cfcr.internal
      ca_cert: ((tls-etcd-server.ca))
      client_cert: ((tls-etcd-client.certificate))
      client_key: ((tls-etcd-client.private_key))
      delete_data_dir_on_stop: false
      dns_health_check_host: 169.254.0.2
      peer_ca_cert: ((tls-etcd-peer.ca))
      peer_cert: ((tls-etcd-peer.certificate))
      peer_key: ((tls-etcd-peer.private_key))
      peer_require_ssl: true
      require_ssl: true
      server_cert: ((tls-etcd-server.certificate))
      server_key: ((tls-etcd-server.private_key))
  release: kubo-etcd

name: master
networks:
- name: default
persistent_disk: 5120
stemcell: trusty
vm_type: master
`

const kubo_master_1_9_2 = `
---
azs:
  - z1
instances: 1
jobs:
- name: kubo-dns-aliases
  release: kubo-1.9.2
- name: secure-var-vcap
  release: kubo-1.9.2
- name: cloud-provider
  properties:
    cloud-provider:
      type: gce
  provides:
    cloud-provider:
      as: master
  release: kubo-1.9.2
- name: docker
  properties:
    bip: 172.17.0.1/24
    default_ulimits:
    - nofile=65536
    env: {}
    flannel: true
    ip_masq: false
    iptables: false
    log_level: error
    storage_driver: overlay
    store_dir: /var/vcap/data
    tls_cacert: ((tls-docker.ca))
    tls_cert: ((tls-docker.certificate))
    tls_key: ((tls-docker.private_key))
  release: docker
- name: kubelet
  release: kubo-1.9.2
  properties:
    api-token: ((kubelet-password))
    register_with_taints: "node-role.kubernetes.io/master=true:NoSchedule"
    labels:
      type: "controller"
    tls:
      kubelet: ((tls-kubelet))
      kubernetes: ((tls-kubernetes))
- name: flanneld
  release: kubo-1.9.2
- name: kube-apiserver
  properties:
    admin-password: ((kubo-admin-password))
    admin-username: admin
    authorization-mode: rbac
    backend_port: 8443
    kube-controller-manager-password: ((kube-controller-manager-password))
    kube-proxy-password: ((kube-proxy-password))
    kube-scheduler-password: ((kube-scheduler-password))
    kubelet-password: ((kubelet-password))
    port: 8443
    route-sync-password: ((route-sync-password))
    service-account-public-key: ((service-account-key.public_key))
    tls:
      kubernetes:
        ca: ((tls-kubernetes.ca))
        certificate: ((tls-kubernetes.certificate))
        private_key: ((tls-kubernetes.private_key))
  release: kubo-1.9.2
- name: kube-controller-manager
  properties:
    api-token: ((kube-controller-manager-password))
    service-account-private-key: ((service-account-key.private_key))
    tls:
      kubernetes: ((tls-kubernetes))
  release: kubo-1.9.2
- name: kube-scheduler
  properties:
    api-token: ((kube-scheduler-password))
    tls:
      kubernetes: ((tls-kubernetes))
  release: kubo-1.9.2
- name: etcd
  properties:
    etcd:
      advertise_urls_dns_suffix: etcd.cfcr.internal
      ca_cert: ((tls-etcd-server.ca))
      client_cert: ((tls-etcd-client.certificate))
      client_key: ((tls-etcd-client.private_key))
      delete_data_dir_on_stop: false
      dns_health_check_host: 169.254.0.2
      peer_ca_cert: ((tls-etcd-peer.ca))
      peer_cert: ((tls-etcd-peer.certificate))
      peer_key: ((tls-etcd-peer.private_key))
      peer_require_ssl: true
      require_ssl: true
      server_cert: ((tls-etcd-server.certificate))
      server_key: ((tls-etcd-server.private_key))
  release: kubo-etcd
name: master
networks:
- name: default
persistent_disk: 5120
stemcell: trusty
vm_type: master
`

const kubo_1_8_6_variables = `
- name: kubo-admin-password
  type: password
- name: kubelet-password
  type: password
- name: kube-proxy-password
  type: password
- name: kube-controller-manager-password
  type: password
- name: kube-scheduler-password
  type: password
- name: route-sync-password
  type: password
- name: kubo_ca
  options:
    common_name: ca
    is_ca: true
  type: certificate
- name: tls-kubelet
  options:
    alternative_names: []
    ca: kubo_ca
    common_name: kubelet.cfcr.internal
    organization: system:nodes
  type: certificate
- name: tls-kubernetes
  options:
    alternative_names:
    - 10.100.200.1
    - kubernetes
    - kubernetes.default
    - kubernetes.default.svc
    - kubernetes.default.svc.cluster.local
    - master.cfcr.internal
    ca: kubo_ca
    organization: system:masters
  type: certificate
- name: tls-docker
  options:
    ca: kubo_ca
    common_name: docker.cfcr.internal
  type: certificate
- name: tls-etcd-server
  options:
    alternative_names:
    - etcd.cfcr.internal
    - '*.etcd.cfcr.internal'
    ca: kubo_ca
    common_name: etcd.cfcr.internal
  type: certificate
- name: tls-etcd-client
  options:
    ca: kubo_ca
    common_name: etcdClient
  type: certificate
- name: tls-etcd-peer
  options:
    alternative_names:
    - '*.etcd.cfcr.internal'
    ca: kubo_ca
    common_name: etcd.cfcr.internal
  type: certificate
- name: kubernetes-dashboard-ca
  options:
    common_name: ca
    is_ca: true
  type: certificate
- name: tls-kubernetes-dashboard
  options:
    alternative_names: []
    ca: kubernetes-dashboard-ca
    common_name: kubernetesdashboard.cfcr.internal
  type: certificate
`

const kubo_1_9_2_variables = `
- name: kubo-admin-password
  type: password
- name: kubelet-password
  type: password
- name: kube-proxy-password
  type: password
- name: kube-controller-manager-password
  type: password
- name: kube-scheduler-password
  type: password
- name: route-sync-password
  type: password
- name: kubo_ca
  options:
    common_name: ca
    is_ca: true
  type: certificate
- name: tls-kubelet
  options:
    alternative_names: []
    ca: kubo_ca
    common_name: kubelet.cfcr.internal
    organization: system:nodes
  type: certificate
- name: tls-kubernetes
  options:
    alternative_names:
    - 10.100.200.1
    - kubernetes
    - kubernetes.default
    - kubernetes.default.svc
    - kubernetes.default.svc.cluster.local
    - master.cfcr.internal
    ca: kubo_ca
    organization: system:masters
  type: certificate
- name: service-account-key
  type: rsa
- name: tls-docker
  options:
    ca: kubo_ca
    common_name: docker.cfcr.internal
  type: certificate
- name: tls-etcd-server
  options:
    alternative_names:
    - etcd.cfcr.internal
    - '*.etcd.cfcr.internal'
    ca: kubo_ca
    common_name: etcd.cfcr.internal
  type: certificate
- name: tls-etcd-client
  options:
    ca: kubo_ca
    common_name: etcdClient
  type: certificate
- name: tls-etcd-peer
  options:
    alternative_names:
    - '*.etcd.cfcr.internal'
    ca: kubo_ca
    common_name: etcd.cfcr.internal
  type: certificate
- name: tls-heapster
  options:
    alternative_names:
    - heapster.kube-system.svc.cluster.local
    ca: kubo_ca
    common_name: heapster
  type: certificate
- name: tls-influxdb
  options:
    alternative_names: []
    ca: kubo_ca
    common_name: monitoring-influxdb
  type: certificate
- name: kubernetes-dashboard-ca
  options:
    common_name: ca
    is_ca: true
  type: certificate
- name: tls-kubernetes-dashboard
  options:
    alternative_names: []
    ca: kubernetes-dashboard-ca
    common_name: kubernetesdashboard.cfcr.internal
  type: certificate
`

// TODO: configurable machine type for vsphere
// TODO: service account for controller
const cloud_config = `
resource_pools:
- name: default
  network: default
  cloud_properties:
    machine_type: n1-standard-2
    root_disk_size_gb: 100
    root_disk_type: pd-ssd
    tags:
    - no-ip
    - internal
`

var workerInstanceGroups map[string]director.InstanceGroup
var masterInstanceGroups map[string]director.InstanceGroup

var variables map[string][]director.Variable
var releases map[string]director.Release

func init() {
	workerInstanceGroups = map[string]director.InstanceGroup{}
	instanceGroup := director.InstanceGroup{}
	if err := yaml.Unmarshal([]byte(kubo_worker_1_9_2), &instanceGroup); err != nil {
		panic(fmt.Errorf("unmarshalling kubo worker 1.9.2: %v", err))
	}
	workerInstanceGroups["1.9.2"] = instanceGroup

	instanceGroup = director.InstanceGroup{}
	if err := yaml.Unmarshal([]byte(kubo_worker_1_8_6), &instanceGroup); err != nil {
		panic(fmt.Errorf("unmarshalling kubo worker 1.8.6: %v", err))
	}
	workerInstanceGroups["1.8.6"] = instanceGroup

	masterInstanceGroups = map[string]director.InstanceGroup{}
	instanceGroup = director.InstanceGroup{}
	if err := yaml.Unmarshal([]byte(kubo_master_1_9_2), &instanceGroup); err != nil {
		panic(fmt.Errorf("unmarshalling kubo master 1.9.2: %v", err))
	}
	masterInstanceGroups["1.9.2"] = instanceGroup

	instanceGroup = director.InstanceGroup{}
	if err := yaml.Unmarshal([]byte(kubo_master_1_8_6), &instanceGroup); err != nil {
		panic(fmt.Errorf("unmarshalling kubo master 1.8.6: %v", err))
	}
	masterInstanceGroups["1.8.6"] = instanceGroup

	releases = map[string]director.Release{}
	releases["kubo-1.8.6"] = director.Release{Name: "kubo-1.8.6",
		Url:     "https://storage.googleapis.com/test-boku-kubo-releases/kubo-release-1.8.6.tgz",
		Version: "0+dev.4",
		Sha1:    "71ad779845ed7d444c0a12647ac35a777148f40d"}
	releases["kubo-1.9.2"] = director.Release{Name: "kubo-1.9.2",
		Url:     "https://storage.googleapis.com/test-boku-kubo-releases/kubo-release-1.9.2.tgz",
		Version: "0+dev.6",
		Sha1:    "8f97ad894ea58471de2e0d0dfac885206bcabc68"}

	variables = map[string][]director.Variable{}
	tmpVariables := []director.Variable{}
	if err := yaml.Unmarshal([]byte(kubo_1_9_2_variables), &tmpVariables); err != nil {
		panic(fmt.Errorf("unmarshalling kubo variables 1.9.2: %v", err))
	}
	variables["kubo-1.9.2"] = tmpVariables
	tmpVariables = []director.Variable{}
	if err := yaml.Unmarshal([]byte(kubo_1_8_6_variables), &tmpVariables); err != nil {
		panic(fmt.Errorf("unmarshalling kubo variables 1.8.6: %v", err))
	}
	variables["kubo-1.8.6"] = tmpVariables
}

func (ManifestGenerator) Generate(machine *v1alpha1.Machine, cluster *v1alpha1.Cluster, ip string) (string, error) {
	name := machine.Name
}

func (ManifestGenerator) InstanceGroup(machine v1alpha1.Machine) (director.InstanceGroup, error) {

	ig, ok := workerInstanceGroups[machine.Spec.Versions.Kubelet]
	if apiutil.IsMaster(&machine) {
		ig, ok = masterInstanceGroups[machine.Spec.Versions.ControlPlane]
	}
	if !ok {
		return director.InstanceGroup{}, errors.New("unsupported version")
	}

	// Specialize the generic InstanceGroup for spec
	ig.Name = machine.ObjectMeta.Name

	providerConfig := map[string]string{}
	// Ignoring error, these values aren't strictly required
	json.Unmarshal([]byte(machine.Spec.ProviderConfig), &providerConfig)

	for i := range ig.Jobs {
		if ig.Jobs[i].Name == "cloud-provider" {
			ig.Jobs[i].Properties = map[string]interface{}{
				"cloud-provider": map[string]interface{}{
					"type": "gce",
					"gce": map[string]string{
						"project-id":      providerConfig["project"],
						"network-name":    providerConfig["networkName"],
						"worker-node-tag": providerConfig["workerNodeTag"],
					},
				},
			}
		}
	}

	return ig, nil
}

func (ManifestGenerator) ReleasesAndVariables(manifest *director.Manifest) ([]director.Release, []director.Variable, error) {
	releaseMap := map[string]director.Release{}
	variableMap := map[string]director.Variable{}

	// Add non-kubo releases that were in the manifest already
	for _, release := range manifest.Releases {
		if ok, _ := regexp.MatchString(`kubo-\d+\.\d+\.\d+`, release.Name); !ok {
			releaseMap[release.Name] = release
		}
	}

	// Add any kubo releases that the instance groups in the manifest need
	// Then add any variables that are needed for that release
	for _, ig := range manifest.InstanceGroups {
		for _, job := range ig.Jobs {
			if job.Name == "kubelet" || job.Name == "kube-apiserver" {
				releaseName := job.Release
				releaseMap[releaseName] = releases[releaseName]
				for _, variable := range variables[releaseName] {
					variableMap[variable.Name] = variable
				}
			}
		}
	}

	manifestReleases := []director.Release{}
	for _, release := range releaseMap {
		manifestReleases = append(manifestReleases, release)
	}
	manifestVariables := []director.Variable{}
	for _, variable := range variableMap {
		manifestVariables = append(manifestVariables, variable)
	}

	return manifestReleases, manifestVariables, nil
}
