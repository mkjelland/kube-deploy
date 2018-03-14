package boshcreateenv

import (
	"encoding/json"
	"errors"

	"fmt"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kube-deploy/cluster-api-bosh/cloud/bosh-create-env/config"
	clusterapiclient "k8s.io/kube-deploy/cluster-api/client"
)

func NewNameResolver(deployment boshdir.Deployment) *nameResolver {
	return &nameResolver{deployment: deployment}
}

type nameResolver struct {
	deployment boshdir.Deployment
}

func (ns *nameResolver) machineToInstanceGroup(cid string) (string, error) {
	vms, err := ns.deployment.VMInfos()
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

// Map the cid to the instance group name
// HACK: If it's in the worker instance group we will return the cid as these aren't managed by the Machine
func (ns *nameResolver) Get(node *corev1.Node, machineClient clusterapiclient.MachinesInterface) (string, error) {
	ips := node.Status.Addresses
	var internalIp string
	for _, ip := range ips {
		if ip.Type == "InternalIP" {
			internalIp = ip.Address
		}
	}

	machines, err := machineClient.List(metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("getting machines in name resolver %v", err)
	}
	for _, m := range machines.Items {
		vmState := config.VMState{}

		err := json.Unmarshal([]byte(m.Status.ProviderState), &vmState)
		if err != nil {
			return "", fmt.Errorf("unmarshalling ProviderState: %v", err)
		}
		if vmState.IP == internalIp {
			return m.Name, nil
		}
	}
	return "", fmt.Errorf("no machine found for ip %v", internalIp)
}
