function trimSlash(str) {
  if (str[str.length - 1] === '/') {
    return str.slice(0, str.length - 1)
  }
  return str
}

export function syncParams (params, type, url) {
  if (type === 'ubuntu') {
    params.ubuntu_repo = trimSlash(url)
  } else if (type === 'docker_ubuntu'){
    params.docker_ubuntu_repo_base_url = trimSlash(url)
    params.docker_ubuntu_repo_gpgkey = trimSlash(url) + '/gpg'
    params.docker_ubuntu_repo_repokey = '9DC858229FC7DD38854AE2D88D81803C0EBFCD88'
  } else if (type === 'centos') {
    params.centos_repo = trimSlash(url)
  } else if (type === 'docker_centos') {
    params.docker_rh_repo_base_url = trimSlash(url) + '/{{ ansible_distribution_major_version }}/$basearch/stable'
    params.docker_rh_repo_gpgkey = trimSlash(url) + '/gpg'
  }
}