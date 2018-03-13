package boshcreatenv

// Machine.Status.ProviderState Schema
type VMState struct {
	// the state file create-env generates
	State []byte // boshcfg.DeploymentState
	IP    string
	// generated values for the VM; only keeping for debugging (unless drain is needed)
	Vars map[string]interface{}
}

type NetworkConfig struct {
	// Subnet
	Netmask         string                      `json:"netmask"`
	Range           string                      `json:"range"`
	Gateway         string                      `json:"gateway"`
	DNS             []string                    `json:"dns"`
	CloudProperties map[interface{}]interface{} `json:"cloud_properties"`
	Reserved        []string                    `json:"reserved"`
}
type ClusterProviderConfig struct {
	Network NetworkConfig `json:"network"`

	// Cluster level properties
	Certificate   map[string]string `json:"certificate"`
	ApiToken      string            `json:"api_token"`
	MasterIp      string            `json:"master_ip"`
	CloudProvider map[string]string `json:"cloud_provider"`
}
