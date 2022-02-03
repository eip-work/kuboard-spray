const mixin = {
  computed: {
    tableData () {
      let result = []
      let inv = this.cluster.inventory
      let res = this.cluster.resourcePackage.data
      let kube_version = res.kubernetes.kube_version
      let k8s = {
        name: 'kubernetes',
        children: [
          { name: 'kube-apiserver', version: kube_version },
          { name: 'kubectl', version: kube_version },
          { name: 'kubelet', version: kube_version },
          { name: 'kubeproxy', version: kube_version },
          { name: 'kubeadm', version: kube_version },
        ]
      }
      result.push(k8s)

      let containerEngine = { name: 'container_engine', version: undefined, children: [] }
      for (let ce of res.container_engine) {
        if (ce.container_manager == inv.all.children.target.vars.container_manager) {
          if (ce.container_manager === 'docker') {
            containerEngine.children.push({ name: 'docker', version: ce.params.docker_version })
            containerEngine.children.push({ name: 'containerd', version: ce.params.docker_containerd_version })
          } else {
            containerEngine.children.push({ name: 'containerd', version: ce.params.containerd_version })
          }
        }
      }
      result.push(containerEngine)

      result.push({ name: 'etcd_cluster', version: undefined, children: [
        { name: 'etcd', version: res.etcd.etcd_version }
      ] })

      let networkPlugin = { name: 'network_plugin', children: [] }
      for (let np of res.network_plugin) {
        if (np.name === inv.all.children.target.children.k8s_cluster.vars.kube_network_plugin) {
          for (let param in np.params) {
            if (param.indexOf('_version') > 0) {
              networkPlugin.children.push({ name: param.slice(0, param.length - 8), version: np.params[param]})
            }
          }
        }
      }
      result.push(networkPlugin)
      
      let dependency = { name: 'dependency', children: [] }
      for (let dep of res.dependency) {
        dependency.children.push({ name: dep.name, version: dep.version })
      }
      result.push(dependency)

      let addon = { name: 'addon', children: [] }
      for (let add of res.addon) {
        if (!inv.all.children.target.children.k8s_cluster.vars[add.target]) {
          continue
        }
        for (let key in add.params) {
          if (key.indexOf('_version') > 0) {
            let v = { name: add.name }
            v.version = add.params[key]
            addon.children.push(v)
          }
        }
      }
      result.push(addon)

      return result
    },
  },
}

export default mixin
