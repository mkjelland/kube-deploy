package config

type ClusterProvider struct {
	Network ClusterProviderNetwork `json:"network"`

	// Cluster level properties
	MasterIp             string                 `json:"master_ip"`
	CloudProvider        map[string]interface{} `json:"cloud_provider"`
	WorkerServiceAccount string                 `json:"worker_service_account"`
	DeploymentVars       map[string]interface{} `json:"deployment_vars"`
}

type ClusterProviderNetwork struct {
	// Subnet
	Netmask         string                 `json:"netmask"`
	Range           string                 `json:"range"`
	Gateway         string                 `json:"gateway"`
	DNS             []string               `json:"dns"`
	CloudProperties map[string]interface{} `json:"cloud_properties"`
	Reserved        []string               `json:"reserved"`
}
