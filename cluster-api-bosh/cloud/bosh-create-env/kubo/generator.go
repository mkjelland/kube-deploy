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

	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/cppforlife/go-patch/patch"
	"k8s.io/kube-deploy/cluster-api-bosh/cloud/bosh-create-env/config"
	"k8s.io/kube-deploy/cluster-api/api/cluster/v1alpha1"
	apiutil "k8s.io/kube-deploy/cluster-api/util"
)

func NewManifestGenerator() *ManifestGenerator {
	return &ManifestGenerator{}
}

type ManifestGenerator struct{}

func (g ManifestGenerator) Generate(machine *v1alpha1.Machine, cluster *v1alpha1.Cluster, ip string) (string, error) {
	// TODO merging versions?
	ops, ok := workerInstanceGroups[machine.Spec.Versions.Kubelet]
	if apiutil.IsMaster(machine) {
		ops, ok = masterInstanceGroups[machine.Spec.Versions.ControlPlane]
	}
	if !ok {
		return "", errors.New("unsupported version")
	}

	tpl := boshtpl.NewTemplate([]byte(base_manifest))

	providerConfig := config.ClusterProvider{}
	if err := json.Unmarshal([]byte(cluster.Spec.ProviderConfig), &providerConfig); err != nil {
		return "", fmt.Errorf("unmarshalling ClusterProviderConfig: %v", err)
	}

	// TODO the nested format for variables isn't working right now?
	vars := boshtpl.StaticVariables(providerConfig.DeploymentVars)
	vars["name"] = machine.Name
	vars["network_cidr"] = providerConfig.Network.Range
	vars["network_ip"] = ip
	vars["network_gw"] = providerConfig.Network.Gateway
	vars["network_dns"] = providerConfig.Network.DNS
	vars["network_cloud_properties"] = providerConfig.Network.CloudProperties
	vars["cloud_provider"] = providerConfig.CloudProvider
	vars["worker_service_account"] = providerConfig.DeploymentVars["worker_service_account"]

	// TODO ExpectAllVarsUsed should be true for strictness
	bytes, err := tpl.Evaluate(vars, patch.Ops(ops), boshtpl.EvaluateOpts{ExpectAllKeys: false, ExpectAllVarsUsed: false})
	if err != nil {
		return "", fmt.Errorf("evaluating template: %v", err)
	}

	return string(bytes), nil
}
