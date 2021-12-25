const field = {
  en: {
    kuboardspray_resource_package: 'Resource Package',
    kuboardspray_os_mirror_type: 'OS type',
    kuboardspray_os_mirror_name: 'OS Mirror Name',
    kuboardspray_os_mirror_kind: 'Create Method',
    'kuboardspray_os_mirror_kind-existing': 'Use existing OS mirror',
    'kuboardspray_os_mirror_kind-provision': 'Create with KuboardSpray',
    kuboardspray_os_mirror_url: 'Mirror URL',
    
    http_proxy: 'HTTP_PROXY',
    https_proxy: 'HTTPS_PROXY',
    no_proxy: 'NO_PROXY',


    ansible_connection: 'ansible_connection',
    'ansible_connection-local': 'local',
    'ansible_connection-ssh': 'ssh',

    yum_repo: 'yum repository',
    ubuntu_repo: 'ubuntu repository',
  },
  zh: {
    kuboardspray_resource_package: '资源包',
    kuboardspray_resource_package_placeholder: '请选择资源包',

    kuboardspray_os_mirror_type: '操作系统类型',
    kuboardspray_os_mirror_name: '镜像源名称',
    kuboardspray_os_mirror_kind: '创建方式',
    'kuboardspray_os_mirror_kind-existing': '使用已有 OS 镜像',
    'kuboardspray_os_mirror_kind-provision': '通过 KuboardSpray 创建',
    kuboardspray_os_mirror_url: '镜像访问地址',

    http_proxy: 'HTTP_PROXY',
    http_proxy_placeholder: 'HTTP Proxy',
    https_proxy: 'HTTPS_PROXY',
    https_proxy_placeholder: 'HTTPS Proxy',
    no_proxy: 'NO_PROXY',
    no_proxy_placeholder: 'NO_PROXY',
    
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
    ansible_become_user: '切换到用户',
    ansible_become_user_placeholder: 'KuboardSpray 登录该主机后，使用 su 命令切换到用户名',
    ansible_become_password: '切换密码',
    ansible_become_password_placeholder: '切换密码',

    
    etcd_deployment_type: 'ETCD 部署方式',
    etcd_deployment_type_placeholder: 'ETCD deployment type',
    etcd_member_name: 'ETCD 成员名称',
    etcd_member_name_placeholder: 'etcd_member_name',

    kube_network_plugin: '网络插件',
    cluster_name: '集群名称',
    cluster_name_placeholder: 'cluster domain name',
    event_ttl_duration: '事件保留时长',
    auto_renew_certificates: '自动更新证书',

    yum_repo: 'yum repo',
    yum_repo_placeholder: 'yum 镜像源',
    ubuntu_repo: 'ubuntu repo',
    ubuntu_repo_placeholder: 'ubuntu 镜像源',


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
  }
}

export default field