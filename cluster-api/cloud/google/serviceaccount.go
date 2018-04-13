/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package google

import (
	"fmt"
	"os/exec"

	"github.com/golang/glog"
	clusterv1 "k8s.io/kube-deploy/cluster-api/pkg/apis/cluster/v1alpha1"
	"k8s.io/kube-deploy/cluster-api/util"
)

const (
	MasterNodeServiceAccountPrefix        = "k8s-master"
	WorkerNodeServiceAccountPrefix        = "k8s-worker"
	MachineControllerServiceAccountPrefix = "k8s-machine-controller"
	MachineControllerSecret               = "machine-controller-credential"
	ClusterAnnotationPrefix               = "service-account-"
)

var (
	MasterNodeRoles = []string{
		"compute.instanceAdmin",
		"compute.networkAdmin",
		"compute.securityAdmin",
		"compute.viewer",
		"iam.serviceAccountUser",
		"storage.admin",
		"storage.objectViewer",
	}
	WorkerNodeRoles = []string{
		"compute.viewer",
	}
	MachineControllerRoles = []string{
		"compute.instanceAdmin.v1",
		"iam.serviceAccountActor",
	}
)

// Returns the email address of the service account that should be used
// as the default service account for this machine
func (gce *GCEClient) GetDefaultServiceAccountForMachine(cluster *clusterv1.Cluster, machine *clusterv1.Machine) string {
	if util.IsMaster(machine) {
		return cluster.ObjectMeta.Annotations[ClusterAnnotationPrefix + MasterNodeServiceAccountPrefix]
	} else {
		return cluster.ObjectMeta.Annotations[ClusterAnnotationPrefix + WorkerNodeServiceAccountPrefix]
	}
}

// Creates a GCP service account for the master node, granted permissions
// that allow the control plane to provision disks and networking resources
func (gce *GCEClient) CreateMasterNodeServiceAccount(cluster *clusterv1.Cluster, initialMachines []*clusterv1.Machine) error {
	_, _, err := gce.createServiceAccount(MasterNodeServiceAccountPrefix, MasterNodeRoles, cluster, initialMachines)

	return err
}

// Creates a GCP service account for the worker node
func (gce *GCEClient) CreateWorkerNodeServiceAccount(cluster *clusterv1.Cluster, initialMachines []*clusterv1.Machine) error {
	_, _, err := gce.createServiceAccount(WorkerNodeServiceAccountPrefix, WorkerNodeRoles, cluster, initialMachines)

	return err
}

// Creates a GCP service account for the machine controller, granted the
// permissions to manage compute instances, and stores its credentials as a
// Kubernetes secret.
func (gce *GCEClient) CreateMachineControllerServiceAccount(cluster *clusterv1.Cluster, initialMachines []*clusterv1.Machine) error {
	accountId, project, err := gce.createServiceAccount(MachineControllerServiceAccountPrefix, MachineControllerRoles, cluster, initialMachines);
	if err != nil {
		return err
	}

	email := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project)

	localFile := accountId + "-key.json"
	err = run("gcloud", "--project", project, "iam", "service-accounts", "keys", "create", localFile, "--iam-account", email)
	if err != nil {
		return fmt.Errorf("couldn't create service account key: %v", err)
	}

	err = run("kubectl", "create", "secret", "generic", MachineControllerSecret, "--from-file=service-account.json="+localFile)
	if err != nil {
		return fmt.Errorf("couldn't import service account key as credential: %v", err)
	}

	if err := run("rm", localFile); err != nil {
		glog.Error(err)
	}

	return nil
}

func (gce *GCEClient) createServiceAccount(serviceAccountPrefix string, roles []string, cluster *clusterv1.Cluster, initialMachines []*clusterv1.Machine) (string, string, error) {
	if len(initialMachines) == 0 {
		return "", "", fmt.Errorf("machine count is zero, cannot create service a/c")
	}

	// TODO: use real go bindings
	// Figure out what projects the service account needs permission to.
	projects, err := gce.getProjects(initialMachines)
	if err != nil {
		return "", "", err
	}

	// The service account needs to be created in a single project, so just
	// use the first one, but grant permission to all projects in the list.
	project := projects[0]
	accountId := serviceAccountPrefix + "-" + util.RandomString(5)

	err = run("gcloud", "--project", project, "iam", "service-accounts", "create", "--display-name=" + serviceAccountPrefix + " service account", accountId)
	if err != nil {
		return "", "", fmt.Errorf("couldn't create service account: %v", err)
	}

	email := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project)

	for _, project := range projects {
		for _, role := range roles {
			err = run("gcloud", "projects", "add-iam-policy-binding", project, "--member=serviceAccount:"+email, "--role=roles/" + role)
			if err != nil {
				return "", "", fmt.Errorf("couldn't grant permissions to service account: %v", err)
			}
		}
	}

	if cluster.ObjectMeta.Annotations == nil {
		cluster.ObjectMeta.Annotations = make(map[string]string)
	}
	cluster.ObjectMeta.Annotations[ClusterAnnotationPrefix + serviceAccountPrefix] = email

	return accountId, project, nil
}

func (gce *GCEClient) DeleteMasterNodeServiceAccount(cluster *clusterv1.Cluster, machines []*clusterv1.Machine) error {
	return gce.deleteServiceAccount(MasterNodeServiceAccountPrefix, MasterNodeRoles, cluster, machines)
}

func (gce *GCEClient) DeleteWorkerNodeServiceAccount(cluster *clusterv1.Cluster, machines []*clusterv1.Machine) error {
	return gce.deleteServiceAccount(WorkerNodeServiceAccountPrefix, WorkerNodeRoles, cluster, machines)
}

func (gce *GCEClient) DeleteMachineControllerServiceAccount(cluster *clusterv1.Cluster, machines []*clusterv1.Machine) error {
	return gce.deleteServiceAccount(MachineControllerServiceAccountPrefix, MachineControllerRoles, cluster, machines)
}

func (gce *GCEClient) deleteServiceAccount(serviceAccountPrefix string, roles []string, cluster *clusterv1.Cluster, machines []*clusterv1.Machine) error {
	if len(machines) == 0 {
		glog.Info("machine count is zero, cannot determine project for service a/c deletion")
		return nil
	}

	projects, err := gce.getProjects(machines)
	if err != nil {
		return err
	}
	project := projects[0]
	var email string
	if cluster.ObjectMeta.Annotations != nil {
		email = cluster.ObjectMeta.Annotations[ClusterAnnotationPrefix + serviceAccountPrefix]
	}

	if email == "" {
		glog.Info("No service a/c found in cluster.")
		return nil
	}

	for _, role := range roles {
		err = run("gcloud", "projects", "remove-iam-policy-binding", project, "--member=serviceAccount:"+email, "--role=roles/" + role)
	}

	if err != nil {
		return fmt.Errorf("couldn't remove permissions to service account: %v", err)
	}

	err = run("gcloud", "--project", project, "iam", "service-accounts", "delete", email)
	if err != nil {
		return fmt.Errorf("couldn't delete service account: %v", err)
	}
	return nil
}

func (gce *GCEClient) getProjects(machines []*clusterv1.Machine) ([]string, error) {
	// Figure out what projects the service account needs permission to.
	var projects []string
	for _, machine := range machines {
		config, err := gce.providerconfig(machine.Spec.ProviderConfig)
		if err != nil {
			return nil, err
		}

		projects = append(projects, config.Project)
	}
	return projects, nil
}

func run(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	if out, err := c.CombinedOutput(); err != nil {
		return fmt.Errorf("error: %v, output: %s", err, string(out))
	}
	return nil
}
