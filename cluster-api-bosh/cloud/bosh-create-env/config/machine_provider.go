package config

type MachineProvider struct {
	// TODO support machine-level cloud properties or network settings?
	CloudProperties map[string]interface{}
}
