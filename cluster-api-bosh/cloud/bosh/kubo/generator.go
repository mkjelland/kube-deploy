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
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
	"k8s.io/kube-deploy/cluster-api-bosh/cloud/bosh/director"
	"k8s.io/kube-deploy/cluster-api/api/cluster/v1alpha1"
)

func NewManifestGenerator() *ManifestGenerator {
	return &ManifestGenerator{}
}

type ManifestGenerator struct{}

const kubo_1_9_2 = `
---
- azs:
  - z1
  instances: 1
  jobs:
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

const kubo_1_8_6 = `
---
- azs:
  - z1
  instances: 1
  jobs:
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
var workerInstanceGroups map[string]director.InstanceGroup{}

var releases map[string]director.Release{}

func init() {
	instanceGroup := director.InstanceGroup{}
	if err := yaml.Unmarshal([]byte(kubo_1_9_2), &instanceGroup); err != nil {
		panic(fmt.Errorf("unmarshalling kubo 1.9.2: %v", err))
	}
	workerInstanceGroups["1.9.2"] = instanceGroup

	instanceGroup = director.InstanceGroup{}
	if err := yaml.Unmarshal([]byte(kubo_1_8_6), &instanceGroup); err != nil {
		panic(fmt.Errorf("unmarshalling kubo 1.8.6: %v", err))
	}
	workerInstanceGroups["1.8.6"] = instanceGroup
}

func (ManifestGenerator) InstanceGroup(spec v1alpha1.MachineSpec) (director.InstanceGroup, error) {
	ig, ok := workerInstanceGroups[spec.Versions.Kubelet]
	if !ok {
		return director.InstanceGroup{}, errors.New("unsupported version")
	}

	// Specialize the generic InstanceGroup for spec
	ig.Name = spec.Name

	for i := range ig.Jobs {
		if ig.Jobs[i].Name == "cloud-provider" {
			// BOKU: Fill in these properties from the spec.ProviderConfig
			gce := map[string]string{
				"project-id": "",
				"network-name": "",
				"worker-node-tag": "",
			}
			ig.Jobs[i].Properties = map[string]interface{}{"gce": gce}
		}
	}
	return ig, nil
}

func (ManifestGenerator) Releases(manifest *director.Manifest) ([]director.Release, error) {
	return nil, errors.New("BOKU: NYI")
}
