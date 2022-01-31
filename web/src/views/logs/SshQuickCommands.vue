<i18n>
en:
  kube_control_plane_cmds: Kube Control Plane
  get_nodes: List nodes
  get_namespaces: List namespaces
  get_all_pods: List all pods
  get_kube_system_pods: List pods in kube-system
  view_config: View cluster config
  cluster_version: View cluster version

  etcd_cmds: ETCD
  etcd_member_list: List etcd members
  etcd_member_status: List ETCD Endpoint status

  kube_node_cmds: Kube Node
  kubelet_version: View kubelet version
  kubelet_status: View kubelet status
  kubelet_journal: View kubelet logs

  crictl_ps: List containers
  crictl_images: List images
zh:
  kube_control_plane_cmds: 控制节点的常用命令
  get_nodes: 获取节点列表
  get_namespaces: 获取名称空间列表
  get_all_pods: 获取所有容器组列表
  get_kube_system_pods: 获取 kube-system 容器组列表
  view_config: 查看集群配置
  cluster_version: 查看集群版本

  etcd_cmds: ETCD节点的常用命令
  etcd_member_list: ETCD成员列表
  etcd_member_status: ETCD Endpoint状态

  kube_node_cmds: 工作节点的常用命令
  kubelet_version: 查看 kubelet 版本
  kubelet_status: 查看 kubelet 状态
  kubelet_journal: 查看 kubelet 日志

  crictl_ps: 查看容器组列表
  crictl_images: 查看容器镜像列表

  calico_node_list: Calico 节点列表
  calico_node_status: Calico 节点状态
</i18n>

<template>
  <div>
    <el-popover v-if="roles.kube_control_plane" v-model:visible="control_plane_commands_visible" placement="top-start" width="50vw" trigger="click">
      <template #reference>
        <el-button style="margin-left: 10px" type="primary">{{ $t('kube_control_plane_cmds') }}</el-button>
      </template>
      <div v-for="(cmds, type) in control_plane_commands" :key="type">
        <div class="app_text_mono commands_title">{{type}}</div>
        <SshQuickCommand v-model:visible="control_plane_commands_visible" v-for="(command, index) in cmds" :key="type + index" :command="command" :socket="socket"></SshQuickCommand>
      </div>
    </el-popover>
    <el-popover v-if="roles.etcd" v-model:visible="etcd_commands_visible" placement="top-start" width="50vw" trigger="click">
      <template #reference>
        <el-button style="margin-left: 10px" type="warning">{{ $t('etcd_cmds') }}</el-button>
      </template>
      <div v-for="(cmds, type) in etcd_commands" :key="type">
        <div class="app_text_mono commands_title">{{type}}</div>
        <SshQuickCommand v-model:visible="etcd_commands_visible" v-for="(command, index) in cmds" :key="type + index" :command="command" :socket="socket"></SshQuickCommand>
      </div>
    </el-popover>
    <el-popover v-if="roles.kube_node" v-model:visible="kube_node_commands_visible" placement="top-start" width="50vw" trigger="click">
      <template #reference>
        <el-button style="margin-left: 10px" type="success">{{ $t('kube_node_cmds') }}</el-button>
      </template>
      <div v-for="(cmds, type) in kube_node_commands" :key="type">
        <div class="app_text_mono commands_title">{{type}}</div>
        <SshQuickCommand v-model:visible="kube_node_commands_visible" v-for="(command, index) in cmds" :key="type + index" :command="command" :socket="socket"></SshQuickCommand>
      </div>
    </el-popover>
  </div>
</template>

<script>
import SshQuickCommand from './SshQuickCommand.vue'

export default {
  props: {
    roles: { type: Object, required: false, default: () => {return {}} },
    socket: { type: Object, required: false, default: () => {return {}} },
    cluster: { type: Object, required: false, default: undefined },
  },
  data() {
    return {
      control_plane_commands_visible: false,
      etcd_commands_visible: false,
      kube_node_commands_visible: false,
    }
  },
  computed: {
    control_plane_commands() {
      let temp = {
        kubectl: [
          { name: this.$t('view_config'), cmd: 'kubectl config view' },
          { name: this.$t('cluster_version'), cmd: 'kubectl version' },
          { name: this.$t('get_nodes'), cmd: 'kubectl get nodes -o wide' },
          { name: this.$t('get_namespaces'), cmd: 'kubectl get namespaces -o wide' },
          { name: this.$t('get_all_pods'), cmd: 'kubectl get pods --all-namespaces -o wide' },
          { name: this.$t('get_kube_system_pods'), cmd: 'kubectl get pods -n kube-system -o wide' },
        ],
        kubelet: [
          { name: this.$t('kubelet_version'), cmd: 'kubelet version' },
          { name: this.$t('kubelet_status'), cmd: 'systemctl status kubelet' },
        ],
        crictl: [
          { name: this.$t('crictl_ps'), cmd: 'crictl ps'},
          { name: this.$t('crictl_images'), cmd: 'crictl images'},
        ]
      }
      temp = Object.assign(temp, this.kube_node_commands)
      return temp
    } ,
    etcd_commands () {
      let statusCmd = 'etcdctl endpoint health'
      if (this.cluster && this.cluster.inventory) {
        let inventory = this.cluster.inventory
        statusCmd += ' --endpoints='
        for (let key in inventory.all.children.target.children.etcd.hosts) {
          statusCmd += inventory.all.hosts[key].ip || inventory.all.hosts[key].ansible_host
          statusCmd += ':2379,'
        }
        statusCmd = statusCmd.slice(0, statusCmd.length - 1)
      }
      return {
        etcdctl: [
          { name: this.$t('etcd_member_list'), cmd: 'etcdctl member list' },
          { name: this.$t('etcd_member_status'), cmd: statusCmd},
        ],
      }
    },
    kube_node_commands() {
      let temp = {
        kubelet: [
          { name: this.$t('kubelet_version'), cmd: 'kubelet version' },
          { name: this.$t('kubelet_status'), cmd: 'systemctl status kubelet' },
          { name: this.$t('kubelet_journal'), cmd: 'journalctl -u kubelet.service'}
        ],
        crictl: [
          { name: this.$t('crictl_ps'), cmd: 'crictl ps'},
          { name: this.$t('crictl_images'), cmd: 'crictl images'},
        ]
      }
      if (this.cluster.inventory.all.children.target.children.k8s_cluster.vars.kube_network_plugin === 'calico') {
        temp.calicoctl = [
          { name: this.$t('calico_node_list'), cmd: 'calicoctl get nodes' },
          { name: this.$t('calico_node_status'), cmd: 'calicoctl node status' },
        ]
      }
      return temp
    }
  },
  components: { SshQuickCommand },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">
.commands_title {
  padding-bottom: 10px;
  color: var(--el-color-primary);
  font-weight: bolder;
  font-size: 15px;
}
</style>
