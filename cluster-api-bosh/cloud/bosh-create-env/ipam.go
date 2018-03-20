package boshcreateenv

import (
	"encoding/json"
	"fmt"
	"net"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kube-deploy/cluster-api-bosh/cloud/bosh-create-env/config"
)

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
	boshCluster := config.ClusterProvider{}
	if err := json.Unmarshal([]byte(cluster.Spec.ProviderConfig), &boshCluster); err != nil {
		return "", fmt.Errorf("unmarshalling ClusterProviderConfig: %v", err)
	}

	// Build a list of IPs not to use
	// HACK: reserved IPs must be single IP addresses (no ranges)
	// TODO reserve network and broadcast ip
	reservedIps := boshCluster.VM.Network.Reserved
	list, err := b.machineClient.List(metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("listing machine records: %v", err)
	}

	for _, m := range list.Items {
		vmState := config.VMState{}

		if m.Status.ProviderState == "" {
			continue
		}

		err := json.Unmarshal([]byte(m.Status.ProviderState), &vmState)
		if err != nil {
			return "", fmt.Errorf("unmarshalling ProviderState: %v", err)
		}

		reservedIps = append(reservedIps, vmState.IP)
	}

	ip, ipnet, err := net.ParseCIDR(boshCluster.VM.Network.Range)
	if err != nil {
		return "", fmt.Errorf("parsing BoshCluster.Range: %v", err)
	}
	for isReserved(reservedIps, ip) || !ipnet.Contains(ip) {
		incrementIp(ip)
	}

	return fmt.Sprintf("%v", ip), nil
}
