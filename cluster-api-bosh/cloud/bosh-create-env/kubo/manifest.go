package kubo

// TODO: configurable machine type for vsphere
// TODO: configurable stemcell
// TODO: service account for controller
const base_manifest = `
name: ((name))
releases: # appended by role/version ops
- name: ((cloud_release_name))
  version: ((cloud_release_version))
  url: ((cloud_release_url))
  sha1: ((cloud_release_sha1))
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
  cloud_properties: ((vm_cloud_properties))
  stemcell: ((vm_stemcell))
networks:
- name: default
  type: manual
  subnets:
  - range: ((vm_network_cidr))
    gateway: ((vm_network_gw))
    static:
    - ((vm_network_ip))
    dns: ((vm_network_dns))
    cloud_properties: ((vm_network_cloud_properties))
cloud_provider:
  mbus: https://mbus:((mbus_bootstrap_password))@((vm_network_ip)):6868
  cert: ((mbus_bootstrap_ssl))
  template:
    name: ((cloud_release_job))
    release: ((cloud_release_name))
  properties:
    agent: {mbus: "https://mbus:((mbus_bootstrap_password))@0.0.0.0:6868"}
    blobstore: {provider: local, path: /var/vcap/micro_bosh/data/cache}
    ntp:
    - 169.254.169.254
    # cheating
    google: ((cloud_release_properties))
    vcenter: ((cloud_release_properties))
instance_groups:
- name: ((name))
  networks:
  - name: default
    static_ips:
    - ((vm_network_ip))
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
      common_name: ((vm_network_ip))
      alternative_names: [((vm_network_ip))]

`
