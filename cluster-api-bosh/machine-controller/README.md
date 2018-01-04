# Machine Controller
## Build

```bash
$ cd $GOPATH/src/k8s.io/
$ git clone git@github.com:kubernetes/kube-deploy.git
$ cd kube-deploy/cluster-api/machine-controller
$ go build
```

## Run
### Locally
1) Spin up a cluster with at least a master that uses kubeadm
2) Get the kubeadm join token. Run `kubeadm token list` in master vm
3) Run `gcloud auth application-default login` to get default credentials
4) Run controller for google cloud cluster `./machine-controller --cloud google --kubeconfig ~/.kube/config --token {step 1 token}`

### As a Pod
1. `kubectl create secret generic kubeconfig --from-file=~/.kube/config`
1. `cd $GOPATH/src/k8s.io/cluster-api-bosh/machine-controller`
1. `docker build -t <Your Dockerhub ID>/machinecontroller:v1 .`
1. `docker push <Your Dockerhub ID>/machinecontroller:v1`
1. update deployment.yml with fields to allow the machine-controller to talk to BOSH, and your Dockerhub ID
1. `kubectl create -f deployment.yml`
