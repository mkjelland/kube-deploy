package kubo

// TODO: configurable machine type for vsphere
// TODO: configurable stemcell
// TODO: service account for controller
const base_manifest = `
name: ((name))
releases: # appended by role/version ops
- name: bosh-google-cpi
  version: 27.0.0
  url: https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-google-cpi-release?v=27.0.0
  sha1: cbbf73c102b1f27d3db15d95bc971b2b4995c78e
- name: docker
  version: 30.1.4
  url: https://bosh.io/d/github.com/cf-platform-eng/docker-boshrelease?v=30.1.4
  sha1: 90dc20d01a0c8a07242d9e846371a03e1a850073
- name: bosh-dns
  sha1: e38d3e5edd16ce2dca778440b636ab1ea61f892f
  version: 0.0.11
  url: https://bosh.io/d/github.com/cloudfoundry/bosh-dns-release?v=0.0.11
resource_pools:
- name: default
  network: default
  env:
    bosh:
      password: '*'
      mbus:
        cert: ((mbus_bootstrap_ssl))
  cloud_properties:
    service_account: ((worker_service_account))
    machine_type: n1-standard-2
    root_disk_size_gb: 100
    root_disk_type: pd-ssd
    tags:
    - no-ip
    - internal
    zone: us-west1-a
  stemcell:
    url: https://bosh.io/d/stemcells/bosh-google-kvm-ubuntu-trusty-go_agent?v=3468.21
    sha1: 242b39fd71da44a352e19ced355736d6f313aab4
networks:
- name: default
  type: manual
  subnets:
  - range: ((network_cidr))
    gateway: ((network_gw))
    static:
    - ((network_ip))
    dns: ((network_dns))
    cloud_properties: ((network_cloud_properties))
cloud_provider:
  mbus: https://mbus:((mbus_bootstrap_password))@((network_ip)):6868
  cert: ((mbus_bootstrap_ssl))
  template:
    name: google_cpi
    release: bosh-google-cpi
  properties:
    agent: {mbus: "https://mbus:((mbus_bootstrap_password))@0.0.0.0:6868"}
    blobstore: {provider: local, path: /var/vcap/micro_bosh/data/cache}
    ntp:
    - 169.254.169.254
    google:
      # TODO provider config
      project: graphite-test-mkjelland
instance_groups:
- name: ((name))
  networks:
  - name: default
    static_ips:
    - ((network_ip))
  jobs: [] # applied by role/version ops
  instances: 1
  resource_pool: default
variables:
  - name: mbus_bootstrap_password
    type: password
  - name: default_ca
    type: certificate
    options:
      is_ca: true
      common_name: ca
  - name: mbus_bootstrap_ssl
    type: certificate
    options:
      ca: default_ca
      common_name: ((network_ip))
      alternative_names: [((network_ip))]

`
