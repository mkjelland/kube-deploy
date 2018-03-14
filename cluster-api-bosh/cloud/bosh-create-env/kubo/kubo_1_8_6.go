package kubo

const kubo_worker_1_8_6 = `
- path: /releases/-
  type: replace
  value:
    name: kubo-1.8.6
    url: "https://storage.googleapis.com/test-boku-kubo-releases/kubo-release-1.8.6.tgz"
    version: "0+dev.4"
    sha1: "28159994d66e26a7ba54fbebc48237f30dba858f"
- path: /instance_groups/0/jobs
  type: replace
  value:
  - name: kubo-dns-aliases
    release: kubo-1.8.6
  - name: secure-var-vcap
    release: kubo-1.8.6
  - name: flanneld
    release: kubo-1.8.6
  - name: docker
    properties:
      bip: 172.17.0.1/24
      default_ulimits:
      - nofile=65536
      env: {}
      flannel: true
      ip_masq: false
      iptables: false
      log_level: error
      storage_driver: overlay
      store_dir: /var/vcap/data
      tls_cacert: ((tls-docker.ca))
      tls_cert: ((tls-docker.certificate))
      tls_key: ((tls-docker.private_key))
    release: docker
    # - name: cloud-provider
    #   properties:
    #     cloud-provider:
    #       type: gce
    #   provides:
    #     cloud-provider:
    #       as: worker
    #   release: kubo-1.8.6
  - name: kubelet
    properties:
      api-token: ((kubelet-password))
      tls:
        kubelet: ((tls-kubelet))
        kubernetes: ((tls-kubernetes))
    release: kubo-1.8.6
    consumes:
      cloud-provider:
        properties: ((cloud_provider))
  - name: kube-proxy
    properties:
      api-token: ((kube-proxy-password))
      tls:
        kubernetes: ((tls-kubernetes))
    release: kubo-1.8.6

- path: /variables/-
  type: replace
  value:
    name: kubo-admin-password
    type: password
- path: /variables/-
  type: replace
  value:
    name: kubelet-password
    type: password
- path: /variables/-
  type: replace
  value:
    name: kube-proxy-password
    type: password
- path: /variables/-
  type: replace
  value:
    name: kube-controller-manager-password
    type: password
- path: /variables/-
  type: replace
  value:
    name: kube-scheduler-password
    type: password
- path: /variables/-
  type: replace
  value:
    name: route-sync-password
    type: password
- path: /variables/-
  type: replace
  value:
    name: kubo_ca
    options:
      common_name: ca
      is_ca: true
    type: certificate
- path: /variables/-
  type: replace
  value:
    name: tls-kubelet
    options:
      alternative_names: []
      ca: kubo_ca
      common_name: kubelet.cfcr.internal
      # organization: system:nodes #
    type: certificate
- path: /variables/-
  type: replace
  value:
    name: tls-kubernetes
    options:
      alternative_names:
      - 10.100.200.1
      - kubernetes
      - kubernetes.default
      - kubernetes.default.svc
      - kubernetes.default.svc.cluster.local
      - master.cfcr.internal
      ca: kubo_ca
      # organization: system:masters
    type: certificate
- path: /variables/-
  type: replace
  value:
    name: tls-docker
    options:
      ca: kubo_ca
      common_name: docker.cfcr.internal
    type: certificate
- path: /variables/-
  type: replace
  value:
    name: tls-etcd-server
    options:
      alternative_names:
      - etcd.cfcr.internal
      - '*.etcd.cfcr.internal'
      ca: kubo_ca
      common_name: etcd.cfcr.internal
    type: certificate
- path: /variables/-
  type: replace
  value:
    name: tls-etcd-client
    options:
      ca: kubo_ca
      common_name: etcdClient
    type: certificate
- path: /variables/-
  type: replace
  value:
    name: tls-etcd-peer
    options:
      alternative_names:
      - '*.etcd.cfcr.internal'
      ca: kubo_ca
      common_name: etcd.cfcr.internal
    type: certificate
- path: /variables/-
  type: replace
  value:
    name: kubernetes-dashboard-ca
    options:
      common_name: ca
      is_ca: true
    type: certificate
- path: /variables/-
  type: replace
  value:
    name: tls-kubernetes-dashboard
    options:
      alternative_names: []
      ca: kubernetes-dashboard-ca
      common_name: kubernetesdashboard.cfcr.internal
    type: certificate

`
