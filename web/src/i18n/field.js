const field = {
  en: {
    please_input: 'Please input ',
    is_required_field: ' is required.',

    kuboardspray_resource_package: 'Resource Package',
    kuboardspray_os_mirror_type: 'OS type',
    kuboardspray_os_mirror_name: 'OS Mirror Name',
    kuboardspray_os_mirror_kind: 'Create Method',
    'kuboardspray_os_mirror_kind-existing': 'Use existing OS Source',
    'kuboardspray_os_mirror_kind-provision': 'Create with KuboardSpray',
    kuboardspray_os_mirror_url: 'Mirror URL',
    
    http_proxy: 'HTTP_PROXY',
    https_proxy: 'HTTPS_PROXY',
    additional_no_proxy: 'NO_PROXY',


    ansible_connection: 'ansible_connection',
    'ansible_connection-local': 'local',
    'ansible_connection-ssh': 'ssh',

  },
  zh: {
    please_input: '请输入',
    is_required_field: '为必填字段',

    kuboardspray_resource_package: '资源包',
    kuboardspray_resource_package_placeholder: '请选择资源包',

    kuboardspray_os_mirror_type: '操作系统类型',
    kuboardspray_os_mirror_name: '镜像源名称',
    kuboardspray_os_mirror_kind: '创建方式',
    'kuboardspray_os_mirror_kind-existing': '使用已有的软件源',
    'kuboardspray_os_mirror_kind-provision': '通过 KuboardSpray 创建',
    kuboardspray_os_mirror_url: '镜像访问地址',

    http_proxy: 'HTTP_PROXY',
    http_proxy_placeholder: 'HTTP Proxy',
    https_proxy: 'HTTPS_PROXY',
    https_proxy_placeholder: 'HTTPS Proxy',
    additional_no_proxy: 'NO_PROXY',
    additional_no_proxy_placeholder: 'NO_PROXY',
    
    ansible_connection: 'ansible 连接类型',
    'ansible_connection-local': 'local',
    'ansible_connection-ssh': 'ssh',

    ansible_host: '主机',
    ansible_port: 'SSH 端口',
    ansible_port_placeholder: 'KuboardSpray 连接该主机时所使用的 SHH 端口',
    
    ansible_user: '用户名',
    ansible_user_placeholder: 'KuboardSpray 连接该主机时所使用的用户名',
    ansible_password: '密码',
    ansible_password_placeholder: 'KuboardSpray 连接该主机时所使用的密码',
    ansible_ssh_private_key_file: '私钥文件',
    ansible_ssh_private_key_file_placeholder: 'KuboardSpray 连接该主机时所使用的私钥文件',
    ansible_become: '切换身份',
    ansible_become_placeholder: 'KuboardSpray 登录到该主机后，是否使用 su 命令切换身份',
    ansible_become_user: '切换到用户',
    ansible_become_user_placeholder: 'KuboardSpray 登录该主机后，使用 su 命令切换到用户名',
    ansible_become_password: '切换密码',
    ansible_become_password_placeholder: '切换密码',

    ansible_python_interpreter: 'Python 路径',

    container_manager: '容器引擎',
    containerd_insecure_registries: 'HTTP 镜像仓库 insecure registries',
    containerd_use_systemd_cgroup: 'systemd_cgroup',

    docker_orphan_clean_up: '清理孤儿容器',
    docker_insecure_registries: 'HTTP 镜像仓库 insecure registries',
    docker_registry_mirrors: '镜像仓库 mirror',

    
    etcd_deployment_type: 'ETCD 部署方式',
    'etcd_deployment_type-docker': '容器化部署',
    'etcd_deployment_type-host': '二进制部署',
    etcd_member_name: 'ETCD 成员名称',
    etcd_member_name_placeholder: 'etcd_member_name',
    etcd_data_dir: 'ETCD 数据目录',

    kube_network_plugin: '网络插件',
    cluster_name: '集群名称',
    cluster_name_placeholder: 'cluster domain name',
    event_ttl_duration: '事件保留时长',
    auto_renew_certificates: '自动更新证书',

    apt_mirror_dir: '数据目录',
    apt_mirror_ubuntu_mirror: '上游镜像源地址',
    apt_mirror_ubuntu_mirror_placeholder: '将从此地址获取镜像下载到本地镜像库',
    apt_mirror_schedule_updates: '是否定期从上游同步',
    'special_time-monthly': '每月',
    'special_time-weekly': '每周',
    'special_time-daily': '每天',
    apt_mirror_populate_repos: '完成后立刻同步',
    apt_mirror_enable_limit_rate: '限制同步速度',
    apt_mirror_limit_rate: '下载速度',
    apt_mirror_nthreads: '下载线程数',

    // Addons
    netchecker_port: '节点端口',
    netchecker_port_placeholder: '默认值 31081，网络检查应用所使用的节点端口',
    agent_report_interval: '检查间隔',
    agent_report_interval_placeholder: '默认值 15s，Agent 发起网络连通性测试的间隔时间',
    netcheck_namespace: '目标名称空间',
    netcheck_namespace_placeholder: '默认值 default，将网络检查组件部署到该名称空间',

    metrics_server_metric_resolution: '采样间隔',
    metrics_server_metric_resolution_placeholder: '默认值 15s，metrics_server 的采样间隔时间',

    // k8s-cluster.yml
    kube_api_anonymous_auth: '允许匿名用户',
    kube_log_level: 'kubelet 日志级别',
    kube_log_level_placeholder: '默认值 2，kubelet 日志级别',
    'kube_log_level-0': 'INFO',
    'kube_log_level-1': 'WARNING',
    'kube_log_level-2': 'ERROR',
    'kube_log_level-3': 'FATAL',

    kube_service_addresses: '服务子网',
    kube_service_addresses_placeholder: '默认值 0.233.0.0/18，服务子网',
    kube_pods_subnet: '容器组子网',
    kube_pods_subnet_placeholder: '默认值 10.233.64.0/18，容器组子网',
    kube_network_node_prefix: '节点子网掩码',
    kube_network_node_prefix_placeholder: '默认值 24，节点子网掩码位数',
    kubelet_max_pods: '节点容器组数',
    kubelet_max_pods_placeholder: '默认值 110，节点最大容器组数',

    nodelocaldns_ip: 'nodelocaldns_ip',
    nodelocaldns_ip_placeholder: '默认值 169.254.25.10，缓存 DNS 容器的虚拟 IP'
  }
}

export default field