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

package boshcreateenv

import (
	"encoding/json"
	"fmt"

	"errors"

	"k8s.io/kube-deploy/cluster-api-bosh/cloud/bosh-create-env/config"
	"k8s.io/kube-deploy/cluster-api-bosh/cloud/bosh-create-env/kubo"

	clusterv1 "k8s.io/kube-deploy/cluster-api/api/cluster/v1alpha1"
	"k8s.io/kube-deploy/cluster-api/client"
	apiutil "k8s.io/kube-deploy/cluster-api/util"
)

type BOSHClient struct {
	machineClient client.MachinesInterface
	clusterClient client.ClustersInterface

	generator ManifestGenerator
}

type ManifestGenerator interface {
	Generate(machine *clusterv1.Machine, cluster *clusterv1.Cluster, ip string) (string, error)
}

func (b *BOSHClient) CreateMachineController(cluster *clusterv1.Cluster, initialMachines []*clusterv1.Machine) error {
	return errors.New("NYI")
}

func NewMachineActuator(clusterClient client.ClustersInterface, machineClient client.MachinesInterface) (*BOSHClient, error) {
	return &BOSHClient{
		machineClient: machineClient,
		clusterClient: clusterClient,
		generator:     kubo.NewManifestGenerator(),
	}, nil
}

func (b *BOSHClient) Create(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	if apiutil.IsMaster(machine) {
		return errors.New("master node creation NYI")
	}

	boshState, err := b.deployWorker(machine)
	if err != nil {
		return fmt.Errorf("deploy worker: %v", err)
	}
	stateBytes, err := json.Marshal(boshState)
	if err != nil {
		return fmt.Errorf("marshaling state: %v", err)
	}
	machine.Status.ProviderState = string(stateBytes)
	_, err = b.machineClient.Update(machine)
	return err
}

func (b *BOSHClient) Delete(machine *clusterv1.Machine) error {
	return errors.New("NYI")
}

func (b *BOSHClient) PostDelete(cluster *clusterv1.Cluster, machines []*clusterv1.Machine) error {
	return nil
}

func (b *BOSHClient) Update(cluster *clusterv1.Cluster, goalMachine *clusterv1.Machine) error {
	if apiutil.IsMaster(goalMachine) {
		return errors.New("master node updating NYI")
	}
	return errors.New("NYI")
}

func (b *BOSHClient) Exists(machine *clusterv1.Machine) (bool, error) {
	if machine.Status.ProviderState == "" {
		return false, nil
	}

	vmState := config.VMState{}

	err := json.Unmarshal([]byte(machine.Status.ProviderState), &vmState)
	if err != nil {
		return false, fmt.Errorf("unmarshalling ProviderState: %v", err)
	} else if vmState.IP != "" {
		// TODO better to review bosh state for vm_cid or something more concrete
		return true, nil
	}

	return false, nil
}

func (b *BOSHClient) GetIP(machine *clusterv1.Machine) (string, error) {
	return "", errors.New("NYI")
}

func (b *BOSHClient) GetKubeConfig(master *clusterv1.Machine) (string, error) {
	return "", errors.New("NYI")
}
