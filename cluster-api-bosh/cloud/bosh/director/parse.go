/*
Copyright 2017 The Kubernetes Authors.

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

package director

import (
	"errors"
	"strings"

	"fmt"

	yaml "gopkg.in/yaml.v2"
	"k8s.io/kube-deploy/cluster-api/api/cluster/v1alpha1"
)

// BOKU: Define the real release struct here
type Release map[string]interface{}

// Manifest represents a BOSH Manifest
type Manifest struct {
	Name   string
	Update interface{}
	//Networks       []interface{}
	//ResourcePools  []interface{} `yaml:"resource_pools"`
	//DiskPools      []interface{} `yaml:"disk_pools"`
	//Jobs           []Job
	InstanceGroups []Job                       `yaml:"instance_groups"`
	Properties     map[interface{}]interface{} `yaml:"properties,omitempty"`
	//Tags           map[string]string
	Features map[interface{}]interface{}
	Releases []Release
	// Don't expose cloud-config properties. We need Marshall functions that will give us the right fields. For now
	// drop the fields not relevant to the deployment manifest.
	Stemcells []map[string]interface{} `yaml:",omitempty"`
	Variables []map[string]interface{}
	AddOns    []map[string]interface{} `yaml:"addons,omitempty"`
}

// Job represents a definition of a BOSH Job or Instance Group
type Job struct {
	Name      string
	Instances int
	//	Lifecycle string
	//Templates []interface{}
	Jobs               []interface{} `yaml:"jobs"`
	Networks           []interface{}
	PersistentDisk     int         `yaml:"persistent_disk,omitempty"`
	PersistentDiskType interface{} `yaml:"persistent_disk_type,omitempty"`
	//PersistentDiskPool string `yaml:"persistent_disk_pool"`
	//ResourcePool string `yaml:"resource_pool"`
	Stemcell         interface{}                 `yaml:"stemcell,omitempty"`
	VMType           interface{}                 `yaml:"vm_type"`
	Properties       map[interface{}]interface{} `yaml:"properties,omitempty"`
	AvailbilityZones []interface{}               `yaml:"azs"`
}

// Parse hydrates a Manifest from a string containing a YAML Manifest
// TODO: Strict parsing should be used to ensure no data loss
func Parse(val string) (*Manifest, error) {
	parsed := &Manifest{}
	return parsed, yaml.Unmarshal([]byte(val), parsed)
}

// Deletes an InstanceGroup by Name
func (m *Manifest) DeleteInstanceGroup(name string) error {
	return errors.New("BOKU: NYI")
}

// BOKU: Delete these
// findInstanceGroupsByType returns all matching instance groups denoted
// by the same prefix
func (m *Manifest) findInstanceGroupsByType(name string) []Job {
	var jobs []Job

	for _, job := range m.InstanceGroups {
		if strings.HasPrefix(job.Name, name) {
			jobs = append(jobs, job)
		}
	}

	return jobs
}

// createInstanceGroup uses the configuration options in machineSpec
// to generate a BOSH instance group
func (m *Manifest) createInstanceGroup(src, dest string, machineSpec v1alpha1.MachineSpec) (Job, error) {
	templates := m.findInstanceGroupsByType(src)
	if len(templates) == 0 {
		return Job{}, fmt.Errorf("can not find template for: %s", src)
	}
	template := templates[0]
	template.Name = dest
	return template, nil
}

// AddWorker adds a new Worker instance group to the Manifest and returns the instance group name
func (m *Manifest) UpdateWorker(name string, machineSpec v1alpha1.MachineSpec) error {
	err := m.DeleteWorker(name)
	if err != nil {
		return err
	}
	err = m.AddWorker(name, machineSpec)
	if err != nil {
		return err
	}
	return nil
}

// AddWorker adds a new Worker instance group to the Manifest and returns the instance group name
func (m *Manifest) AddWorker(name string, machineSpec v1alpha1.MachineSpec) error {
	worker, err := m.createInstanceGroup("worker", name, machineSpec)
	if err != nil {
		return err
	}

	// A single instance allows lookup of VM CID -> Instance Group to enable deletion of a specific VM
	worker.Instances = 1

	m.InstanceGroups = append(m.InstanceGroups, worker)

	return nil
}

// DeleteWorker removes a Worker instance by instance group name
func (m *Manifest) DeleteWorker(name string) error {
	var jobs []Job
	for _, j := range m.InstanceGroups {
		if j.Name != name {
			jobs = append(jobs, j)
		}
	}

	if len(jobs) == len(m.InstanceGroups) {
		return fmt.Errorf("no instance group found for deletion: %s", name)
	}

	m.InstanceGroups = jobs
	return nil
}
