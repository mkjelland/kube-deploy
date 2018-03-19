package kubo

const kubo_worker_1_8_6 = `
- path: /releases/-
  type: replace
  value:
    name: kubo
    url: "https://storage.googleapis.com/test-boku-kubo-releases/kubo-release-1.8.6-dev.tgz"
    version: "0.11.1+dev.11"
    sha1: "e02cfd63a18da223bd8d9783a54a6eb2f4c8e9b0"
- path: /releases/-
  type: replace
  value:
    name: os-conf
    version: 18
    url: https://bosh.io/d/github.com/cloudfoundry/os-conf-release?v=18
    sha1: 78d79f08ff5001cc2a24f572837c7a9c59a0e796
- path: /instance_groups/0/jobs
  type: replace
  value:
  - name: bosh-dns
    release: bosh-dns
    properties:
      records_file: /var/vcap/jobs/kubo-dns-aliases/dns/records.json
      cache:
        enabled: true
      health:
        enabled: true
        server:
          tls: ((/dns_healthcheck_server_tls))
        client:
          tls: ((/dns_healthcheck_client_tls))  
  - name: user_add
    release: os-conf
    properties:
      users:
      - name: jumpbox
        public_key: ((jumpbox_ssh.public_key))
  - name: kubo-dns-aliases
    release: kubo
    properties:
      master_ip: ((master_address))
      bosh-dns-records: ((bosh-dns-aliases))
    consumes:
      etcd:
        instances:
        - name: master
          index: 0
          address: ((master_address))
        properties:
          etcd:
            advertise_urls_dns_suffix: etcd.cfcr.internal
            ca_cert: ((tls-etcd-client.ca))
            client_cert: ((tls-etcd-client.certificate))
            client_key: ((tls-etcd-client.private_key))
  - name: secure-var-vcap
    release: kubo
    properties: {}
  - name: flanneld
    release: kubo
    properties: {}
    consumes:
      etcd:
        instances:
        - name: master
          index: 0
          address: etcd.cfcr.internal
        properties:
          name: master
          etcd:
            name: master
            advertise_urls_dns_suffix: etcd.cfcr.internal
            ca_cert: ((tls-etcd-client.ca))
            client_cert: ((tls-etcd-client.certificate))
            client_key: ((tls-etcd-client.private_key))
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
  - name: cloud-provider
    properties:
      cloud-provider: ((cloud_provider))
    release: kubo
  - name: kubelet
    properties:
      api-token: ((kubelet-password))
      tls:
        kubelet: ((tls-kubelet))
        kubernetes: ((tls-kubernetes))
    release: kubo
  - name: kube-proxy
    properties:
      api-token: ((kube-proxy-password))
      tls:
        kubernetes: ((tls-kubernetes))
    release: kubo

- path: /variables/-
  type: replace
  value:
    name: kubo-admin-password
    type: password
- path: /variables/-
  type: replace
  value:
    name: jumpbox_ssh
    type: ssh
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
- path: /variables/-
  type: replace
  value:
    name: /dns_healthcheck_tls_ca
    type: certificate
    options:
      is_ca: true
      common_name: dns-healthcheck-tls-ca
- path: /variables/-
  type: replace
  value:
    name: /dns_healthcheck_server_tls
    type: certificate
    options:
      ca: /dns_healthcheck_tls_ca
      common_name: health.bosh-dns
      extended_key_usage:
      - server_auth
- path: /variables/-
  type: replace
  value:
    name: /dns_healthcheck_client_tls
    type: certificate
    options:
      ca: /dns_healthcheck_tls_ca
      common_name: health.bosh-dns
      extended_key_usage:
      - client_auth
`
