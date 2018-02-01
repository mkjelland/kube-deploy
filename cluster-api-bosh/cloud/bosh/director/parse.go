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
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

// BOKU: Define the real release struct here

type Release struct {
	Name    string `yaml:"name"`
	Url     string `yaml:"url"`
	Version string `yaml:"version"`
	Sha1    string `yaml:"sha1"`
}

// Manifest represents a BOSH Manifest
type Manifest struct {
	Name   string
	Update interface{}
	//Networks       []interface{}
	//ResourcePools  []interface{} `yaml:"resource_pools"`
	//DiskPools      []interface{} `yaml:"disk_pools"`
	//Jobs           []InstanceGroup
	InstanceGroups []InstanceGroup             `yaml:"instance_groups"`
	Properties     map[interface{}]interface{} `yaml:"properties,omitempty"`
	//Tags           map[string]string
	Features map[interface{}]interface{}
	Releases []Release `yaml:"releases,omitempty"`
	// Don't expose cloud-config properties. We need Marshall functions that will give us the right fields. For now
	// drop the fields not relevant to the deployment manifest.
	Stemcells []map[string]interface{} `yaml:",omitempty"`
	Variables []Variable
	AddOns    []map[string]interface{} `yaml:"addons,omitempty"`
}

type Variable struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options,omitempty"`
}

// InstanceGroup represents a definition of a BOSH InstanceGroup or Instance Group
type InstanceGroup struct {
	Name      string `yaml:"name"`
	Instances int    `yaml:"instances"`
	//	Lifecycle string
	//Templates []interface{}
	Jobs               []Job         `yaml:"jobs"`
	Networks           []interface{} `yaml:"networks"`
	PersistentDisk     int           `yaml:"persistent_disk,omitempty"`
	PersistentDiskType interface{}   `yaml:"persistent_disk_type,omitempty"`
	//PersistentDiskPool string `yaml:"persistent_disk_pool"`
	//ResourcePool string `yaml:"resource_pool"`
	Stemcell         interface{}                 `yaml:"stemcell,omitempty"`
	VMType           interface{}                 `yaml:"vm_type"`
	Properties       map[interface{}]interface{} `yaml:"properties,omitempty"`
	AvailbilityZones []interface{}               `yaml:"azs"`
}

type Job struct {
	Name       string                 `yaml:"name"`
	Release    string                 `yaml:"release"`
	Consumes   map[string]interface{} `yaml:"consumes,omitempty"`
	Provides   map[string]interface{} `yaml:"provides,omitempty"`
	Properties map[string]interface{} `yaml:"properties,omitempty"`
}

// Parse hydrates a Manifest from a string containing a YAML Manifest
// TODO: Strict parsing should be used to ensure no data loss
func Parse(val string) (*Manifest, error) {
	parsed := &Manifest{}
	return parsed, yaml.Unmarshal([]byte(val), parsed)
}

// Deletes an InstanceGroup by Name
func (m *Manifest) DeleteInstanceGroup(name string) error {
	var jobs []InstanceGroup
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
