package cluster

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/eip-work/kuboard-spray/common"
	"github.com/eip-work/kuboard-spray/constants"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type CreateClusterReq struct {
	Name            string `json:"name" binding:"required"`
	ResourcePackage string `json:"resource_package" binding:"required"`
}

func CreateCluster(c *gin.Context) {
	var req CreateClusterReq
	e := c.ShouldBindJSON(&req)

	if e != nil {
		common.HandleError(c, http.StatusBadRequest, "request should include field name and resourcePackage", e)
		return
	}

	clusterDir := constants.GET_DATA_CLUSTER_DIR() + "/" + req.Name

	_, err := ioutil.ReadDir(clusterDir)
	if err == nil {
		common.HandleError(c, http.StatusConflict, "conflict with a existing cluster", err)
		return
	}
	err = common.CreateDirIfNotExists(clusterDir)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to create folder", err)
		return
	}

	inventoryFilePath := clusterDir + "/inventory.yaml"
	inventoryFile, err := os.OpenFile(inventoryFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to create inventory file "+inventoryFilePath, err)
		return
	}

	defer inventoryFile.Close()

	template := getInventoryTemplate()
	template = strings.ReplaceAll(template, "KUBOARDSPRAY_RESOURCE_PACKAGE", req.ResourcePackage)

	inventory := map[string]interface{}{}
	yaml.Unmarshal([]byte(template), inventory)

	resourcePackagePath := constants.GET_DATA_RESOURCE_DIR() + "/" + req.ResourcePackage + "/content"
	resourcePackage, err := common.ParseYamlFile(resourcePackagePath + "/package.yaml")
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot parse package.yaml", err)
		return
	}

	common.PopulateKuboardSprayVars(inventory, "cluster", req.Name)

	addons := common.MapGet(resourcePackage, "data.addon").([]interface{})
	for _, a := range addons {
		addon := a.(map[string]interface{})
		target := addon["target"].(string)
		if addon["lifecycle"] != nil {
			lifecycle := addon["lifecycle"].(map[string]interface{})
			if lifecycle["install_by_default"] == true {
				common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars."+target, true)
			} else {
				common.MapSet(inventory, "all.children.target.children.k8s_cluster.vars."+target, false)
			}
		}
	}

	// if runtime.GOARCH == "arm64" {
	// 	common.MapSet(inventory, "all.children.target.vars.container_manager", "docker")
	// }

	content, _ := yaml.Marshal(inventory)

	_, err = inventoryFile.WriteString(string(content))

	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "failed to write inventory file: "+inventoryFilePath, err)
		return
	}

	data := gin.H{}

	err = yaml.Unmarshal([]byte(template), &data)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "cannot parse file: "+inventoryFilePath, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    data,
	})
}

func getInventoryTemplate() string {
	template := `all:
  hosts:
    localhost:
      ansible_connection: local
      kuboardspray_resource_package: KUBOARDSPRAY_RESOURCE_PACKAGE
  children:
    target:
      children:
        etcd:
          hosts: {}
          vars:
            etcd_deployment_type: host
            etcd_data_dir: /var/lib/etcd
            etcd_kubeadm_enabled: false
        k8s_cluster:
          children:
            kube_control_plane:
              hosts: {}
            kube_node:
              hosts: {}
            calico_rr:
              hosts: {}
          vars:
            # addons.yaml
            helm_enabled: false
            registry_enabled: false
            metrics_server_enabled: true
            local_path_provisioner_enabled: false
            local_volume_provisioner_enabled: false
            cephfs_provisioner_enabled: false
            rbd_provisioner_enabled: false
            ingress_nginx_enabled: false
            ingress_publish_status_address: ""
            ingress_alb_enabled: false
            cert_manager_enabled: false
            metallb_enabled: false
            metallb_speaker_enabled: false
            argocd_enabled: false
            krew_enabled: false
            krew_root_dir: "/usr/local/krew"

            # k8s-cluster.yml
            kube_config_dir: /etc/kubernetes
            kube_script_dir: "{{ bin_dir }}/kubernetes-scripts"
            kube_manifest_dir: "{{ kube_config_dir }}/manifests"
            # This is where all the cert scripts and certs will be located
            kube_cert_dir: "{{ kube_config_dir }}/ssl"
            # This is where all of the bearer tokens will be stored
            kube_token_dir: "{{ kube_config_dir }}/tokens"
            kube_api_anonymous_auth: true
            ## Change this to use another Kubernetes version, e.g. a current beta release
            ## kube_version: v1.22.4
            # Where the binaries will be downloaded.
            # Note: ensure that you've enough disk space (about 1G)
            local_release_dir: "/tmp/releases"
            # Random shifts for retrying failed ops like pushing/downloading
            retry_stagger: 5
            # This is the group that the cert creation scripts chgrp the
            # cert files to. Not really changeable...
            kube_cert_group: kube-cert
            # Cluster Loglevel configuration
            kube_log_level: 2
            # Directory where credentials will be stored
            credentials_dir: "{{ inventory_dir }}/credentials"

            # Choose network plugin (cilium, calico, weave or flannel. Use cni for generic cni plugin)
            # Can also be set to 'cloud', which lets the cloud provider setup appropriate routing
            kube_network_plugin: calico

            # Setting multi_networking to true will install Multus: https://github.com/intel/multus-cni
            kube_network_plugin_multus: false

            # Kubernetes internal network for services, unused block of space.
            kube_service_addresses: 10.233.0.0/16

            # internal network. When used, it will assign IP
            # addresses from this range to individual pods.
            # This network must be unused in your network infrastructure!
            kube_pods_subnet: 10.234.0.0/16

            kube_network_node_prefix: 22

            # Configure Dual Stack networking (i.e. both IPv4 and IPv6)
            enable_dual_stack_networks: false
            # Kubernetes internal network for IPv6 services, unused block of space.
            # This is only used if enable_dual_stack_networks is set to true
            # This provides 4096 IPv6 IPs
            kube_service_addresses_ipv6: fd85:ee78:d8a6:8607::1000/116

            # Internal network. When used, it will assign IPv6 addresses from this range to individual pods.
            # This network must not already be in your network infrastructure!
            # This is only used if enable_dual_stack_networks is set to true.
            # This provides room for 256 nodes with 254 pods per node.
            kube_pods_subnet_ipv6: fd85:ee78:d8a6:8607::1:0000/112

            # IPv6 subnet size allocated to each for pods.
            # This is only used if enable_dual_stack_networks is set to true
            # This provides room for 254 pods per node.
            kube_network_node_prefix_ipv6: 120

            # The port the API Server will be listening on.
            kube_apiserver_ip: "{{ kube_service_addresses|ipaddr('net')|ipaddr(1)|ipaddr('address') }}"
            kube_apiserver_port: 6443  # (https)
            # Set to 0 to disable insecure port - Requires RBAC in authorization_modes and kube_api_anonymous_auth: true
            kube_apiserver_insecure_port: 0  # (disabled)

            # Kube-proxy proxyMode configuration.
            # Can be ipvs, iptables
            kube_proxy_mode: ipvs

            # configure arp_ignore and arp_announce to avoid answering ARP queries from kube-ipvs0 interface
            # must be set to true for MetalLB to work
            kube_proxy_strict_arp: false

            # A string slice of values which specify the addresses to use for NodePorts.
            # Values may be valid IP blocks (e.g. 1.2.3.0/24, 1.2.3.4/32).
            # The default empty string slice ([]) means to use all local addresses.
            # kube_proxy_nodeport_addresses_cidr is retained for legacy config
            kube_proxy_nodeport_addresses: >-
              {%- if kube_proxy_nodeport_addresses_cidr is defined -%}
              [{{ kube_proxy_nodeport_addresses_cidr }}]
              {%- else -%}
              []
              {%- endif -%}

            ## Encrypting Secret Data at Rest (experimental)
            kube_encrypt_secret_data: false

            # DNS configuration.
            # Kubernetes cluster name, also will be used as DNS domain
            cluster_name: cluster.local
            # Subdomains of DNS domain to be resolved via /etc/resolv.conf for hostnet pods
            ndots: 2
            # Can be coredns, coredns_dual, manual or none
            dns_mode: coredns
            # Set manual server if using a custom cluster DNS server
            # manual_dns_server: 10.x.x.x
            # Enable nodelocal dns cache
            enable_nodelocaldns: true
            enable_nodelocaldns_secondary: false
            nodelocaldns_ip: 169.254.25.10
            nodelocaldns_health_port: 9254
            nodelocaldns_second_health_port: 9256
            nodelocaldns_bind_metrics_host_ip: false
            nodelocaldns_secondary_skew_seconds: 5
            # Enable k8s_external plugin for CoreDNS
            enable_coredns_k8s_external: false
            coredns_k8s_external_zone: k8s_external.local
            # Enable endpoint_pod_names option for kubernetes plugin
            enable_coredns_k8s_endpoint_pod_names: false

            # Can be docker_dns, host_resolvconf or none
            resolvconf_mode: docker_dns
            # Deploy netchecker app to verify DNS resolve as an HTTP service
            deploy_netchecker: false
            # Ip address of the kubernetes skydns service
            skydns_server: "{{ kube_service_addresses|ipaddr('net')|ipaddr(3)|ipaddr('address') }}"
            skydns_server_secondary: "{{ kube_service_addresses|ipaddr('net')|ipaddr(4)|ipaddr('address') }}"
            dns_domain: "{{ cluster_name }}"

            # Additional container runtimes
            kata_containers_enabled: false

            kubeadm_certificate_key: "{{ lookup('password', credentials_dir + '/kubeadm_certificate_key.creds length=64 chars=hexdigits') | lower }}"

            # K8s image pull policy (imagePullPolicy)
            k8s_image_pull_policy: IfNotPresent

            # audit log for kubernetes
            kubernetes_audit: true

            # dynamic kubelet configuration
            # Note: Feature DynamicKubeletConfig is deprecated in 1.22 and will not move to GA.
            # It is planned to be removed from Kubernetes in the version 1.23.
            # Please use alternative ways to update kubelet configuration.
            dynamic_kubelet_configuration: false

            # define kubelet config dir for dynamic kubelet
            # kubelet_config_dir:
            default_kubelet_config_dir: "{{ kube_config_dir }}/dynamic_kubelet_dir"
            dynamic_kubelet_configuration_dir: "{{ kubelet_config_dir | default(default_kubelet_config_dir) }}"

            # pod security policy (RBAC must be enabled either by having 'RBAC' in authorization_modes or kubeadm enabled)
            podsecuritypolicy_enabled: false

            ## Running on top of openstack vms with cinder enabled may lead to unschedulable pods due to NoVolumeZoneConflict restriction in kube-scheduler.
            ## See https://github.com/kubernetes-sigs/kubespray/issues/2141
            ## Set this variable to true to get rid of this issue
            volume_cross_zone_attachment: false
            ## Add Persistent Volumes Storage Class for corresponding cloud provider (supported: in-tree OpenStack, Cinder CSI,
            ## AWS EBS CSI, Azure Disk CSI, GCP Persistent Disk CSI)
            persistent_volumes_enabled: false

            ## Amount of time to retain events. (default 1h0m0s)
            event_ttl_duration: "1h0m0s"

            ## Automatically renew K8S control plane certificates on first Monday of each month
            auto_renew_certificates: true
            # First Monday of each month
            # auto_renew_certificates_systemd_calendar: "Mon *-*-1,2,3,4,5,6,7 03:{{ groups['kube_control_plane'].index(inventory_hostname) }}0:00"
            kubelet_rotate_certificates: true
            kubelet_rotate_server_certificates: true

      vars:
        bin_dir: /usr/local/bin
        # loadbalancer_apiserver_port: 6443
        # loadbalancer_apiserver_healthcheck_port: 8081
        containerd_use_systemd_cgroup: true
        docker_orphan_clean_up: true
        ## Container runtime
        ## docker for docker, crio for cri-o and containerd for containerd.
        ## Default: containerd
        container_manager: containerd
        kuboardspray_repo_centos: 'AS_IS'
        kuboardspray_repo_docker_centos: AS_IS
        kuboardspray_repo_ubuntu: 'AS_IS'
        kuboardspray_repo_docker_ubuntu: AS_IS
        ansible_python_interpreter: auto
  vars:
    ansible_ssh_pipelining: true
`
	return template
}
