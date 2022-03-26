let docker_centos = [
  'https://download.docker.com/linux/centos/',
  'https://mirrors.aliyun.com/docker-ce/linux/centos/',
  'https://mirrors.tuna.tsinghua.edu.cn/docker-ce/linux/centos/',
]
let temp = {
  ubuntu: [
      'https://repo.huaweicloud.com/ubuntu/',
      'https://mirrors.aliyun.com/ubuntu/',
      'https://mirrors.tuna.tsinghua.edu.cn/ubuntu/',
      'http://cn.archive.ubuntu.com/ubuntu/',
      'https://mirrors.cloud.tencent.com/ubuntu/',
      'http://ftp.sjtu.edu.cn/ubuntu',
      'http://mirrors.163.com/ubuntu/',
      'http://mirrors.nju.edu.cn/ubuntu/',
    ],
  docker_ubuntu: [
      'https://mirrors.tuna.tsinghua.edu.cn/docker-ce/linux/ubuntu/',
      'https://mirrors.cloud.tencent.com/docker-ce/linux/ubuntu/',
    ],
  redhat: [],
  docker_redhat: docker_centos,
  centos: [
    'http://mirrors.aliyun.com/repo/',
    'https://repo.huaweicloud.com/centos/',
    'https://mirrors.tuna.tsinghua.edu.cn/centos/',
    'https://mirrors.cloud.tencent.com/centos/',
  ],
  docker_centos: docker_centos,
  anolis: [
    'http://mirrors.openanolis.cn/anolis/',
    'https://mirrors.aliyun.com/anolis/',
  ],
  docker_anolis: docker_centos,
  rocky: [
    'http://mirrors.nju.edu.cn/rocky',
    'http://mirrors.sdu.edu.cn/rocky',
    'http://mirrors.aliyun.com/rockylinux/',
  ],
  docker_rocky: docker_centos,
  openeuler: [
    'https://repo.huaweicloud.com/repository/conf/openeuler_x86_64.repo',
  ],
  docker_openeuler: docker_centos,
  docker_oraclelinux: docker_centos,
  'docker_kylin linux advanced server': docker_centos,
  // 'opensuse leap': [
  //   'https://mirrors.tuna.tsinghua.edu.cn/opensuse',
  // ],
}
if (window.KuboardSpray.version.arch === 'arm64') {
  temp.openeuler = [
    'https://repo.huaweicloud.com/repository/conf/openeuler_aarch64.repo',
  ]
}

export default temp