function trimSlash(str) {
  if (str[str.length - 1] === '/') {
    return str.slice(0, str.length - 1)
  }
  return str
}

export function syncParams (params, type, url) {
  if (type.indexOf('docker_') === 0) {
    switch (type) {
      case 'docker_ubuntu':
        params.docker_ubuntu_repo_base_url = trimSlash(url)
        params.docker_ubuntu_repo_gpgkey = trimSlash(url) + '/gpg'
        params.docker_ubuntu_repo_repokey = '9DC858229FC7DD38854AE2D88D81803C0EBFCD88'
        break
      case 'docker_centos':
      case 'docker_rocky':
      case 'docker_oraclelinux':
      case 'docker_anolis':
        params[`${type}_repo_base_url`] = trimSlash(url) + '/{{ ansible_distribution_major_version }}/$basearch/stable'
        params[`${type}_repo_gpgkey`] = trimSlash(url) + '/gpg'
        break
      case 'docker_openeuler':
        params.docker_openeuler_repo_base_url = trimSlash(url) + '/8/$basearch/stable'
        params.docker_openeuler_repo_gpgkey = trimSlash(url) + '/gpg'
        break
      default:

    }
  } else {
    params[type+'_repo'] = trimSlash(url)
  }
}