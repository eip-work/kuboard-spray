const field = {
  en: {
    kuboardspray_resource_package: 'Resource Package',
    kuboardspray_resource_package_placeholder: 'Please select Resource Package.',
    
    http_proxy: 'HTTP_PROXY',
    https_proxy: 'HTTPS_PROXY',
    no_proxy: 'NO_PROXY',


    ansible_connection: 'ansible_connection',
    'ansible_connection-local': 'local',
    'ansible_connection-ssh': 'ssh',
  },
  zh: {
    kuboardspray_resource_package: '资源包',
    kuboardspray_resource_package_placeholder: '请选择资源包',
    
    http_proxy: 'HTTP_PROXY',
    https_proxy: 'HTTPS_PROXY',
    no_proxy: 'NO_PROXY',
    
    ansible_connection: 'ansible 连接类型',
    'ansible_connection-local': 'local',
    'ansible_connection-ssh': 'ssh',

    ansible_host: '主机',
    ansible_host_placeholder: 'KuboardSpray 连接该主机时所使用的主机名或 IP 地址',
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
    event_ttl_duration: '事件保留时长',
    auto_renew_certificates: '自动更新证书',

  }
}

export default field