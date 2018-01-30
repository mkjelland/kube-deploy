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

	"k8s.io/kube-deploy/cluster-api-bosh/cloud/bosh/director"
	"k8s.io/kube-deploy/cluster-api/api/cluster/v1alpha1"
)

func NewManifestGenerator() *ManifestGenerator {
	return &ManifestGenerator{}
}

type ManifestGenerator struct{}

func (ManifestGenerator) InstanceGroup(spec v1alpha1.MachineSpec) (director.Job, error) {
	return director.Job{}, errors.New("BOKU: NYI")
}

func (ManifestGenerator) Releases(manifest *director.Manifest) ([]director.Release, error) {
	return nil, errors.New("BOKU: NYI")
}
