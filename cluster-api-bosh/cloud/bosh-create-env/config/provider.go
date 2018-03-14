package config

type ClusterProvider struct {
	Network ClusterProviderNetwork `json:"network"`

	// Cluster level properties
	MasterIp       string                 `json:"master_ip"`
	CloudProvider  map[string]string      `json:"cloud_provider"`
	DeploymentVars map[string]interface{} `json:"deployment_vars"`
}

type ClusterProviderNetwork struct {
	// Subnet
	Netmask         string                      `json:"netmask"`
	Range           string                      `json:"range"`
	Gateway         string                      `json:"gateway"`
	DNS             []string                    `json:"dns"`
	CloudProperties map[interface{}]interface{} `json:"cloud_properties"`
	Reserved        []string                    `json:"reserved"`
}
