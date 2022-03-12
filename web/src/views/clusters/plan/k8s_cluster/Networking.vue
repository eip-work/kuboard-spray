<i18n>
en:
  label: Network
  description: Network parameters to kubernetes.
zh:
  label: 网络参数
  description: Kubernetes 集群的网络参数
</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled anti-freeze>
    <template #more>
      请确保节点的 IP 地址不包含在服务子网和容器组子网当中。
    </template>
    <FieldString :holder="vars" :prop="prop" fieldName="kube_service_addresses" required :rules="service_address_rules"></FieldString>
    <FieldString :holder="vars" :prop="prop" fieldName="kube_pods_subnet" required :rules="service_address_rules"></FieldString>
    <FieldNumber :holder="vars" :prop="prop" fieldName="kube_network_node_prefix" required :rules="node_prefix_rules"></FieldNumber>
    <div class="desc" style="margin-bottom: 5px;">每个节点最多可以有 {{ node_ip_size }} 个容器组 IP 地址</div>
    <div class="desc">最多可以添加 {{ node_max_count }} 个节点</div>
    <FieldNumber :holder="vars" :prop="prop" fieldName="kubelet_max_pods" :rules="kubelet_max_pods_rules"></FieldNumber>
  </ConfigSection>
</template>

<script>
import { Netmask } from 'netmask'

function isInternalIP(value) {
  let flag = false
  let block = new Netmask('10.0.0.0/8')
  if (block.contains(value)) {
    flag = true
  }
  block = new Netmask('172.16.0.0/12')
  if (block.contains(value)) {
    flag = true
  }
  block = new Netmask('192.168.0.0/16') 
  if (block.contains(value)) {
    flag = true
  }
  return flag
}

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      service_address_rules: [
        {
          validator: (rule, value, callback) => {
            if (!isInternalIP(value)) {
              return callback('必须使用内网 IP 段')
            }
            let block = new Netmask(this.vars.kube_service_addresses)
            if (block.contains(this.vars.kube_pods_subnet)) {
              return callback(this.$t('field.kube_service_addresses') + ' 不能包含 ' + this.$t('field.kube_pods_subnet'))
            }
            block = new Netmask(this.vars.kube_pods_subnet)
            if (block.contains(this.vars.kube_service_addresses)) {
              return callback(this.$t('field.kube_pods_subnet') + ' 不能包含 ' + this.$t('field.kube_service_addresses'))
            }
            return callback()
          },
          trigger: 'blur'
        }
      ],
      node_prefix_rules: [
        {
          validator: (rule, value, callback) => {
            if (value >= 26) {
              return callback('不能大于 25，请减小' + this.$t('field.kube_pods_subnet'))
            }
            if (this.node_max_count <= 8) {
              return callback('可以创建的节点数量过少，请增大' + this.$t('field.kube_pods_subnet'))
            }
            return callback()
          },
          trigger: 'blur'
        }
      ],
      kubelet_max_pods_rules: [
        {
          validator: (rule, value, callback) => {
            if (value >= this.node_ip_size / 2) {
              return callback('不能大于或等于节点上可用的容器组 IP 地址数量的一半 ' + this.node_ip_size / 2)
            }
            return callback()
          },
          trigger: 'blur'
        }
      ]
    }
  },
  computed: {
    node_ip_size () {
      let mask = '22'
      if (this.vars.kube_network_node_prefix) {
        mask = this.vars.kube_network_node_prefix
      }
      try {
        let block = new Netmask('10.99.99.99/' + mask)
        return block.size
      } catch(e) {
        return 1
      }
    },
    node_max_count () {
      let cird = '10.234.0.0/16'
      if (this.vars.kube_pods_subnet) {
        cird = this.vars.kube_pods_subnet
      }
      try {
        let block = new Netmask(cird)
        return block.size / this.node_ip_size
      } catch(e) {
        return 0
      }
    },
    prop () {
      return 'all.children.target.children.k8s_cluster.vars'
    },
    vars: {
      get () { return this.cluster.inventory.all.children.target.children.k8s_cluster.vars },
      set () {}
    },
    enabled: {
      get () {return true},
      set () {},
    },
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="css">
.desc {
  margin-left: 120px;
  margin-top: -10px;
  margin-bottom: 10px;
  color: var(--el-text-color-secondary);
}
</style>
