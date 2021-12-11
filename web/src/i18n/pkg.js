const pkg = {
  en: {
    ingress_nginx_enabled: 'Includes Ingress Nginx',
  },
  zh: {
    kuboardspray_kubespray_version: 'Kubespray 版本',
    image_arch: 'CPU 架构',
    gcr_image_repo: 'gcr 源镜像仓库',
    kube_image_repo: 'Kubernetes 源镜像仓库',
    kube_version: 'Kubernetes 版本',
    etcd_version: 'ETCD 版本',
    crun_version: 'crun 版本',
    runc_version: 'runc 版本',
    calico_version: 'calico 版本',
    kuboardspray_docker_enabled: '支持 docker 作为容器引擎',
    kuboardspray_containerd_enabled: '支持 containerd 作为容器引擎',
    containerd_version: 'containerd 版本',
    helm_enabled: '包含 Helm',
    registry_enabled: '包含 registry',
    metrics_server_enabled: '包含 Metrics Server',
    local_path_provisioner_enabled: '包含 Local Path Provisioner',
    local_volume_provisioner_enabled: '包含 Local Volume Provisioner',
    cephfs_provisioner_enabled: '包含 CephFS Provisioner',
    rbd_provisioner_enabled: '包含 RBD Provisioner',
    ingress_nginx_enabled: '包含 Ingress Nginx',
    ingress_alb_enabled: '包含 Ingress ALB',
    cert_manager_enabled: '包含 Cert Manager',
    argocd_enabled: '包含 argocd',
    krew_enabled: '包含 Krew',
    metallb_enabled: '包含 Meta LB',
  }
}

export default pkg