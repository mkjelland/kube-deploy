package kubo

// TODO: configurable machine type for vsphere
// TODO: configurable stemcell
// TODO: service account for controller
const base_manifest = `
name: ((name))
releases: [] # applied by role/version ops
resource_pools:
- name: default
  network: default
  cloud_properties:
    machine_type: n1-standard-2
    root_disk_size_gb: 100
    root_disk_type: pd-ssd
    tags:
    - no-ip
    - internal
	stemcell:
		url: https://bosh.io/d/stemcells/bosh-google-kvm-ubuntu-trusty-go_agent?v=3468.21
		sha1: 242b39fd71da44a352e19ced355736d6f313aab4
networks:
- name: default
	type: manual
	subnets:
	- range: ((network_cidr))
		gateway: ((network_gw))
		static: [((network_ip))]
		dns: ((network_dns))
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
			project: ((project_id))
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
	    common_name: ((internal_ip))
	    alternative_names: [((internal_ip))]

`
