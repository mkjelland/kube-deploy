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

package manifest

import (
	"math/rand"
	"strings"

	"fmt"

	yaml "gopkg.in/yaml.v2"
)

// Manifest represents a BOSH manifest
type Manifest struct {
	Name           string
	Update         interface{}
	Networks       []interface{}
	ResourcePools  []interface{} `yaml:"resource_pools"`
	DiskPools      []interface{} `yaml:"disk_pools"`
	Jobs           []Job
	InstanceGroups []Job `yaml:"instance_groups"`
	Properties     map[interface{}]interface{}
	Tags           map[string]string
	Features       map[interface{}]interface{}
	Releases       []map[string]interface{}
	Stemcells      []map[string]interface{}
	Variables      []map[string]interface{}
}

// Job represents a definition of a BOSH Job or Instance Group
type Job struct {
	Name               string
	Instances          int
	Lifecycle          string
	Templates          []interface{}
	Jobs               []interface{} `yaml:"jobs"`
	Networks           []interface{}
	PersistentDisk     int    `yaml:"persistent_disk"`
	PersistentDiskType string `yaml:"persistent_disk_type"`
	PersistentDiskPool string `yaml:"persistent_disk_pool"`
	ResourcePool       string `yaml:"resource_pool"`
	Stemcell           interface{}
	VMType             interface{} `yaml:"vm_type"`
	Properties         map[interface{}]interface{}
	AvailbilityZones   []interface{} `yaml:"azs"`
}

// Parse hydrates a Manifest from a string containing a YAML manifest
// Strict parsing is used to ensure no data loss
func Parse(val string) (*Manifest, error) {
	parsed := &Manifest{}
	return parsed, yaml.UnmarshalStrict([]byte(val), parsed)
}

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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// duplicateInstanceGroup duplicates an instance group by finding the
// first matching instance group containing `name` and adding a random suffix
func (m *Manifest) duplicateInstanceGroup(name string) (string, error) {
	templates := m.findInstanceGroupsByType(name)
	if len(templates) == 0 {
		return "", fmt.Errorf("can not find template for: %s", name)
	}
	template := templates[0]
	template.Name = fmt.Sprintf("%s_%i_%s", name, len(templates), randStr(5))
	m.InstanceGroups = append(m.InstanceGroups, template)
	return template.Name, nil
}

// AddWorker adds a new Worker instance group to the Manifest and returns the instance group name
func (m *Manifest) AddWorker() (string, error) {
	return m.duplicateInstanceGroup("worker")
}

// DeleteWorker removes a Worker instance by instance group name
func (m *Manifest) DeleteWorker(name string) error {
	var jobs []Job
	for _, j := range m.InstanceGroups {
		if j.Name != name {
			jobs = append(jobs, j)
		}
	}

	if len(jobs) == len(m.Jobs) {
		return fmt.Errorf("no jobs found for deletion: %s", name)
	}

	m.InstanceGroups = jobs
	return nil
}
