package bosh

import (
	"errors"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
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
func (ns *nameResolver) Get(cid string) (string, error) {
	ig, err := ns.machineToInstanceGroup(cid)
	if err != nil {
		return cid, nil
	}

	return ig, nil
}
