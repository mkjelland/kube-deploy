package boshcreatenv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"

	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cppforlife/go-patch/patch"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "k8s.io/kube-deploy/cluster-api/api/cluster/v1alpha1"
)

func (b *BOSHClient) deployWorker(machine *clusterv1.Machine) (*VMState, error) {
	ip, err := b.nextIp()
	if err != nil {
		return nil, err
	}

	clusterList, err := b.clusterClient.List(metav1.ListOptions{})
	if err != nil || len(clusterList.Items) != 1 {
		return nil, fmt.Errorf("invalid cluster list: %v", err)
	}

	manifest, err := b.generator.Generate(machine, &clusterList.Items[0], ip)
	if err != nil {
		return nil, err
	}

	logger := boshlog.NewLogger(boshlog.LevelDebug)
	ui := boshui.NewConfUI(logger)
	defer ui.Flush()

	deps := boshcmd.NewBasicDeps(ui, logger)
	envProvider := func(manifestPath string, statePath string, vars boshtpl.Variables, op patch.Op) boshcmd.DeploymentPreparer {
		return boshcmd.NewEnvFactory(deps, manifestPath, statePath, vars, op).Preparer()
	}

	opts := boshcmd.CreateEnvOpts{}
	opts.Args.Manifest.Bytes = []byte(manifest)

	file, err := ioutil.TempFile("", "bosh-state")
	if err != nil {
		return nil, fmt.Errorf("creating temp state file: %v", err)
	}
	file.Close()

	fmt.Printf("bosh_state.json path: %v", file.Name())

	opts.StatePath = file.Name()
	stage := boshui.NewStage(deps.UI, deps.Time, deps.Logger)
	err = boshcmd.NewCreateEnvCmd(deps.UI, envProvider).Run(stage, opts)
	if err != nil {
		return nil, err
	}

	stateContents, err := ioutil.ReadFile(file.Name())

	vmState := &VMState{}
	vmState.IP = ip
	vmState.State = stateContents
	_, err := b.machineClient.Update(machine)
	if err != nil {
		return nil, err
	}

	return vmState, nil
}

// HACK: Only support /24 networks
func incrementIp(ip net.IP) {
	ip[15] += 1
	if ip[15] == 255 {
		panic("overflow?")
	}
}

func isReserved(reserved []string, ip net.IP) bool {
	for _, r := range reserved {
		if net.ParseIP(r).Equal(ip) {
			return true
		}
	}
	return false
}

func (b *BOSHClient) nextIp() (string, error) {
	clusterList, err := b.clusterClient.List(metav1.ListOptions{})
	if err != nil || len(clusterList.Items) != 1 {
		return "", fmt.Errorf("invalid cluster list: %v", err)
	}
	cluster := clusterList.Items[0]
	boshCluster := ClusterProviderConfig{}
	if err := json.Unmarshal([]byte(cluster.Spec.ProviderConfig), &boshCluster); err != nil {
		return "", fmt.Errorf("unmarshalling ClusterProviderConfig: %v", err)
	}

	// Build a list of IPs not to use
	// HACK: reserved IPs must be single IP addresses (no ranges)
	reservedIps := boshCluster.Network.Reserved
	list, err := b.machineClient.List(metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("listing machine records: %v", err)
	}

	for _, m := range list.Items {
		vmState := VMState{}

		err := json.Unmarshal([]byte(m.Status.ProviderState), &vmState)
		if err != nil {
			return "", fmt.Errorf("unmarshalling ProviderState: %v", err)
		}

		reservedIps = append(reservedIps, vmState.IP)
	}

	ip, ipnet, err := net.ParseCIDR(boshCluster.Network.Range)
	if err != nil {
		return "", fmt.Errorf("parsing BoshCluster.Range: %v", err)
	}
	for isReserved(reservedIps, ip) || !ipnet.Contains(ip) {
		incrementIp(ip)
	}

	return fmt.Sprintf("%v", ip), nil
}
