package kubo

import (
	"fmt"

	"github.com/cppforlife/go-patch/patch"
	yaml "gopkg.in/yaml.v2"
)

var workerInstanceGroups map[string]patch.Ops
var masterInstanceGroups map[string]patch.Ops

func init() {
	{ // 1.9.2
		workerInstanceGroups = map[string]patch.Ops{}
		masterInstanceGroups = map[string]patch.Ops{}
		var unmarshal []patch.OpDefinition

		err := yaml.Unmarshal([]byte(kubo_worker_1_9_2), &unmarshal)
		if err != nil {
			panic(fmt.Errorf("Deserializing ops '%s': %v", "1.9.2", err))
		}

		workerInstanceGroups["1.9.2"], err = patch.NewOpsFromDefinitions(unmarshal)
		if err != nil {
			panic(fmt.Errorf("building ops: %v", err))
		}
	}

	{ // 1.8.6
		var unmarshal []patch.OpDefinition

		err := yaml.Unmarshal([]byte(kubo_worker_1_8_6), &unmarshal)
		if err != nil {
			panic(fmt.Errorf("Deserializing ops '%s': %v", "1.8.6", err))
		}

		workerInstanceGroups["1.8.6"], err = patch.NewOpsFromDefinitions(unmarshal)
		if err != nil {
			panic(fmt.Errorf("building ops: %v", err))
		}
	}
}
