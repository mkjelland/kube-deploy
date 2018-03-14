package boshcreateenv

import (
	"fmt"
	"io/ioutil"

	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cppforlife/go-patch/patch"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kube-deploy/cluster-api-bosh/cloud/bosh-create-env/config"
	clusterv1 "k8s.io/kube-deploy/cluster-api/api/cluster/v1alpha1"
)

func (b *BOSHClient) deployWorker(machine *clusterv1.Machine) (*config.VMState, error) {
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

	varsStoreFile, err := ioutil.TempFile("", "vars-store")
	if err != nil {
		return nil, fmt.Errorf("creating temp vars store file: %v", err)
	}

	// TODO remove file

	varsStore := &boshcmd.VarsFSStore{}
	varsStore.FS = deps.FS
	varsStore.UnmarshalFlag(varsStoreFile.Name())

	opts := boshcmd.CreateEnvOpts{}
	opts.VarFlags.VarsFSStore = *varsStore

	manifestFile, err := ioutil.TempFile("", "bosh-manifest")
	if err != nil {
		return nil, fmt.Errorf("creating manifest file: %v", err)
	}
	manifestFile.Write([]byte(manifest))
	manifestFile.Close()

	// TODO defer remove

	opts.Args.Manifest.FS = deps.FS
	opts.Args.Manifest.Path = manifestFile.Name()
	// opts.Args.Manifest.Bytes = []byte(manifest)

	file, err := ioutil.TempFile("", "bosh-state")
	if err != nil {
		return nil, fmt.Errorf("creating temp state file: %v", err)
	}
	file.Write([]byte("{}"))
	file.Close()

	fmt.Printf("bosh_state.json path: %v", file.Name())

	opts.StatePath = file.Name()
	stage := boshui.NewStage(deps.UI, deps.Time, deps.Logger)
	err = boshcmd.NewCreateEnvCmd(deps.UI, envProvider).Run(stage, opts)
	if err != nil {
		return nil, err
	}

	stateContents, err := ioutil.ReadFile(file.Name())

	vmState := &config.VMState{}
	vmState.IP = ip
	vmState.State = stateContents
	// vmState.Vars = map[string]interface{} // TODO not currently generating vars
	_, err = b.machineClient.Update(machine)
	if err != nil {
		return nil, err
	}

	return vmState, nil
}
