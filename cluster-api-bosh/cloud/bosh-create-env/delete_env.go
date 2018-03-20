package boshcreateenv

import (
	"fmt"
	"io/ioutil"

	"encoding/json"

	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cppforlife/go-patch/patch"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kube-deploy/cluster-api-bosh/cloud/bosh-create-env/config"
	clusterv1 "k8s.io/kube-deploy/cluster-api/api/cluster/v1alpha1"
)

func (b *BOSHClient) deleteWorker(machine *clusterv1.Machine) error {
	vmState := &config.VMState{}
	err := json.Unmarshal([]byte(machine.Status.ProviderState), vmState)
	if err != nil {
		return fmt.Errorf("Error unmarshalling provider state: %v", err)
	}
	createEnvState := vmState.State
	fmt.Printf("createEnvState: %v", string(createEnvState))

	logger := boshlog.NewLogger(boshlog.LevelDebug)
	ui := boshui.NewConfUI(logger)
	defer ui.Flush()

	deps := boshcmd.NewBasicDeps(ui, logger)
	envProvider := func(manifestPath string, statePath string, vars boshtpl.Variables, op patch.Op) boshcmd.DeploymentDeleter {
		return boshcmd.NewEnvFactory(deps, manifestPath, statePath, vars, op).Deleter()
	}

	file, err := ioutil.TempFile("", "bosh-state")
	if err != nil {
		return fmt.Errorf("creating temp state file: %v", err)
	}
	file.Write(createEnvState)
	file.Close()

	fmt.Printf("state file name: %v", file.Name())

	clusterList, err := b.clusterClient.List(metav1.ListOptions{})
	if err != nil || len(clusterList.Items) != 1 {
		return fmt.Errorf("invalid cluster list: %v", err)
	}
	manifest, err := b.generator.Generate(machine, &clusterList.Items[0], vmState.IP)
	if err != nil {
		return err
	}
	manifestFile, err := ioutil.TempFile("", "bosh-manifest")
	if err != nil {
		return fmt.Errorf("creating manifest file: %v", err)
	}
	manifestFile.Write([]byte(manifest))
	manifestFile.Close()

	varsStoreFile, err := ioutil.TempFile("", "vars-store")
	if err != nil {
		return fmt.Errorf("creating temp vars store file: %v", err)
	}
	varsStore := &boshcmd.VarsFSStore{}
	varsStore.FS = deps.FS
	varsStore.UnmarshalFlag(varsStoreFile.Name())

	opts := &boshcmd.DeleteEnvOpts{}
	opts.VarFlags.VarsFSStore = *varsStore
	opts.StatePath = file.Name()
	opts.Args.Manifest.FS = deps.FS
	opts.Args.Manifest.Path = manifestFile.Name()

	stage := boshui.NewStage(deps.UI, deps.Time, deps.Logger)
	err = boshcmd.NewDeleteCmd(deps.UI, envProvider).Run(stage, *opts)
	return err
}
