package config

// Machine.Status.ProviderState Schema
type VMState struct {
	// the state file create-env generates
	State []byte // boshcfg.DeploymentState
	IP    string
	// generated values for the VM; only keeping for debugging (unless drain is needed)
	Vars map[string]interface{}
}
