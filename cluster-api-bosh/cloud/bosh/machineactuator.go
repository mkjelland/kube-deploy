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

package bosh

import (
	"errors"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"k8s.io/kube-deploy/cluster-api-bosh/cloud/bosh/director"

	clusterv1 "k8s.io/kube-deploy/cluster-api/api/cluster/v1alpha1"
	"k8s.io/kube-deploy/cluster-api/client"
	apiutil "k8s.io/kube-deploy/cluster-api/util"
)

type BOSHClient struct {
	machineClient client.MachinesInterface
	deployment    boshdir.Deployment
}

func (b *BOSHClient) CreateMachineController(cluster *clusterv1.Cluster, initialMachines []*clusterv1.Machine) error {
	return errors.New("NYI")
}

func NewMachineActuator(deployment boshdir.Deployment, machineClient client.MachinesInterface) (*BOSHClient, error) {
	return &BOSHClient{
		deployment:    deployment,
		machineClient: machineClient,
	}, nil
}

func (b *BOSHClient) getManifest() (*director.Manifest, error) {
	manifestStr, err := b.deployment.Manifest()
	if err != nil {
		return nil, err
	}

	glog.Infof("fetched manifest: \n%s", manifestStr)
	return director.Parse(manifestStr)
}

func (b *BOSHClient) deployManifest(manifest *director.Manifest) error {
	manifestBytes, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}
	glog.Infof("attempting to deploy: \n%s", string(manifestBytes))

	diff, err := b.deployment.Diff(manifestBytes, false)
	if err != nil {
		return err
	}

	glog.Infof("deployment diff: \n%v", diff)

	return b.deployment.Update(manifestBytes, boshdir.UpdateOpts{})
}

func (b *BOSHClient) Create(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	if apiutil.IsMaster(machine) {
		return errors.New("master node creation NYI")
	}

	manifest, err := b.getManifest()
	if err != nil {
		return err
	}

	err = manifest.AddWorker(machine.ObjectMeta.Name, machine.Spec)
	if err != nil {
		return err
	}

	return b.deployManifest(manifest)
}

func (b *BOSHClient) machineToInstanceGroup(cid string) (string, error) {
	vms, err := b.deployment.VMInfos()
	if err != nil {
		return "", err
	}
	for _, vm := range vms {
		if vm.VMID == cid {
			return vm.JobName, nil
		}
	}

	return "", errors.New("instance not found")
}

func (b *BOSHClient) Delete(machine *clusterv1.Machine) error {
	//ig, err := b.machineToInstanceGroup(machine.ObjectMeta.Name)
	//if err != nil {
	//	return err
	//}

	ig := machine.ObjectMeta.Name

	manifest, err := b.getManifest()
	if err != nil {
		return err
	}

	if err := manifest.DeleteWorker(ig); err != nil {
		return err
	}

	return b.deployManifest(manifest)
}

func (b *BOSHClient) PostDelete(cluster *clusterv1.Cluster, machines []*clusterv1.Machine) error {
	return nil
}

func (b *BOSHClient) Update(cluster *clusterv1.Cluster, goalMachine *clusterv1.Machine) error {
	if apiutil.IsMaster(goalMachine) {
		return errors.New("master node updating NYI")
	}
	manifest, err := b.getManifest()
	if err != nil {
		return err
	}

	err = manifest.UpdateWorker(goalMachine.ObjectMeta.Name, goalMachine.Spec)
	if err != nil {
		return err
	}

	return b.deployManifest(manifest)
}

func (b *BOSHClient) Exists(machine *clusterv1.Machine) (bool, error) {
	_, err := b.machineToInstanceGroup(machine.ObjectMeta.Name)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (b *BOSHClient) GetIP(machine *clusterv1.Machine) (string, error) {
	return "", errors.New("NYI")
}

func (b *BOSHClient) GetKubeConfig(master *clusterv1.Machine) (string, error) {
	return "", errors.New("NYI")
}
