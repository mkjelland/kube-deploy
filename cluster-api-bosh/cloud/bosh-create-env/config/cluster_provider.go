package config

type ClusterProvider struct {
	Cloud ClusterProviderCloud `json:"cloud"`
	VM    ClusterProviderVM    `json:"vm"`

	// Cluster level properties
	MasterIp       string                 `json:"master_ip"`
	DeploymentVars map[string]interface{} `json:"deployment_vars"`
}

type ClusterProviderVM struct {
	CloudProperties map[string]interface{}  `json:"cloud_properties"`
	Network         ClusterProviderNetwork  `json:"network"`
	Stemcell        ClusterProviderStemcell `json:"stemcell"`
}

type ClusterProviderStemcell struct {
	URL  string `json:"url"`
	Sha1 string `json:"sha1"`
}

type ClusterProviderNetwork struct {
	Netmask         string                 `json:"netmask"`
	Range           string                 `json:"range"`
	Gateway         string                 `json:"gateway"`
	DNS             []string               `json:"dns"`
	CloudProperties map[string]interface{} `json:"cloud_properties"`
	Reserved        []string               `json:"reserved"`
}

type ClusterProviderCloud struct {
	Type    string                      `json:"type"`
	Release ClusterProviderCloudRelease `json:"release"`
	Options map[string]interface{}      `json:"options"`
}

type ClusterProviderCloudRelease struct {
	Name       string                 `json:"name"`
	Job        string                 `json:"job"`
	URL        string                 `json:"url"`
	Sha1       string                 `json:"sha1"`
	Version    string                 `json:"version"`
	Properties map[string]interface{} `json:"properties"`
}
