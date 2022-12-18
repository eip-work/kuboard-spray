<i18n>
en:
  label: DNS
  description: Network parameters to kubernetes.
  nodelocaldnsLabel: Node local DNS
  nodelocaldnsHelpStr: "Node local DNS virtual IP, default value configed into pods' /etc/resolv.conf"
  second: second
  apply: Apply
  applyDns: Apply DNS params to cluster.

  verbose: Include task params
  verbose_: 
  verbose_v: May include sensitive data in the trace, e.g. path to files, user name, password.
  verbose_vvv: includes more information in log, only used in development.
  v_: Normal
  v_v: Details
  v_vvv: More
  fork: ansible fork
  fork_more: Max number of nodes can be operated in the installation.

zh:
  label: DNS
  description: Kubernetes 集群的 DNS 参数
  nodelocaldnsLabel: 节点 DNS
  nodelocaldnsHelpStr: 节点 DNS 的虚拟 IP 地址，将默认配置到各容器组的 /etc/resolv.conf 文件
  upstream_dns_servers_placeholder: 上游 DNS 服务地址
  zone_placeholder: "例如： mycompany.com"
  cache_placeholder: "缓存时长，单位：秒"
  second: "秒"
  apply: 应 用
  applyDns: 将 DNS 设置参数应用到集群

  verbose: 显示任务参数
  verbose_: 正常输出的日志，通常选用此选项
  verbose_v: 日志中会包含部分敏感信息，例如：文件路径、用户名密码等
  verbose_vvv: 日志中会包含最详细的信息，通常只在开发阶段使用
  v_: 正常
  v_v: 详细
  v_vvv: 更多
  fork: 并发数量
  fork_more: 安装过程中可以同时操作的目标节点的最大数量。ansible fork.

</i18n>

<template>
  <ConfigSection v-model:enabled="enabled" :label="$t('label')" :description="$t('description')" disabled anti-freeze helpLink="https://kuboard.cn/guide/options/dns.html">
    <template #more>
      集群 DNS 设置
    </template>
    <FieldArray :holder="vars" :prop="prop" newItemTemplate="" fieldName="upstream_dns_servers" :itemRules="upstreamDnsServerItemRules" anti-freeze>
      <template v-slot:editItem="scope">
        <el-input v-model.trim="vars.upstream_dns_servers[scope.index]" :placeholder="$t('upstream_dns_servers_placeholder')"></el-input>
      </template>
      <template #help>
        <div style="color: var(--el-text-color-secondary)">
          <template v-if="vars.upstream_dns_servers == undefined || vars.upstream_dns_servers.length == 0">
            不配置上游 DNS 服务地址时，对于集群外的 DNS 地址，将使用节点上 /etc/resolv.conf 的配置进行域名解析。
          </template>
          <template v-else>配置了上游 DNS 服务地址，忽略节点上 /etc/resolv.conf</template>
        </div>
      </template>
      <template #helpView>
        <div style="color: var(--el-text-color-secondary)">
          <template v-if="vars.upstream_dns_servers == undefined || vars.upstream_dns_servers.length == 0">
            未配置上游 DNS 服务地址，对于集群外的 DNS 地址，将使用节点上 /etc/resolv.conf 的配置进行域名解析。
          </template>
          <template v-else>配置了上游 DNS 服务地址，忽略节点上 /etc/resolv.conf</template>
        </div>
      </template>
    </FieldArray>
    <div v-if="isClusterOnline && editMode === 'view'" style="text-align: right; margin: -10px 0 10px;">
      <ExecuteTask :history="cluster.history" placement="left"
          :title="$t('apply')" :startTask="applyDnsConfig" @refresh="$emit('refresh')" @visibleChange="onVisibleChange">
          <template #reference>
            <el-button icon="el-icon-download" type="warning" plain>{{$t('apply')}}</el-button>
          </template>
          <div>
            <el-form ref="form" :model="form" @submit.prevent.stop label-position="left" label-width="120px">
              <div class="app_block_title">{{$t('label')}}</div>
              <div style="margin-bottom: 20px; color: var(--el-text-color-secondary);">{{ $t('applyDns') }}</div>
              <el-form-item :label="$t('verbose')">
                <el-radio-group v-model="form.verbose">
                  <el-radio-button label="">{{$t('v_')}}</el-radio-button>
                  <el-radio-button label="v">{{$t('v_v')}}</el-radio-button>
                  <el-radio-button label="vvv">{{$t('v_vvv')}}</el-radio-button>
                </el-radio-group>
                <div style="width: 350px; margin-left: 0; color: #aaa; font-size: 12px;">{{$t('verbose_' + form.verbose)}}</div>
              </el-form-item>
              <el-form-item :label="$t('fork')" style="margin-top: 10px;">
                <el-input-number v-model="form.fork" :step="2" style="width: 166px;"></el-input-number>
                <div style="width: 350px; margin-left: 0; color: #aaa; font-size: 12px;">{{$t('fork_more')}}</div>
              </el-form-item>
              <el-form-item :label="$t('offlineNodes')" class="app_margin_top" v-if="offlineNodes.length > 0">
                <div class="form_description">{{ $t('offlineNodesDesc') }}</div>
                <template v-for="node in offlineNodes" :key="'exclude' + node">
                  <el-tooltip class="box-item" effect="dark" :content="pingpong[node].message" placement="top-end">
                    <el-tag type="danger" effect="dark" style="margin: 0 10px 10px 0;">
                      <span class="app_text_mono" style="font-size: 14px; margin-right: 10px;">{{ node }}</span>
                      <el-icon :size="14" style="width: 14px; height: 14px; vertical-align: top;">
                        <el-icon-question-filled></el-icon-question-filled>
                      </el-icon>
                    </el-tag>
                  </el-tooltip>
                </template>
              </el-form-item>
            </el-form>
          </div>
      </ExecuteTask>
    </div>
    <!-- <FieldArray :holder="vars" :prop="prop" :newItemTemplate="{zones:[], nameservers: [], cache: 5}" fieldName="coredns_external_zones" anti-freeze isBlockItem>
      <template v-slot:editItem="scope">
        <FieldArray :holder="vars.coredns_external_zones[scope.index]" fieldName="zones" anti-freeze label-width="90px">
          <template v-slot:editItem="zone">
            <el-input v-model.trim="vars.coredns_external_zones[scope.index].zones[zone.index]" :placeholder="$t('zone_placeholder')"></el-input>
          </template>
        </FieldArray>
        <FieldArray :holder="vars.coredns_external_zones[scope.index]" fieldName="nameservers" anti-freeze label-width="90px" style="margin-top: 15px;">
          <template v-slot:editItem="nameserver">
            <el-input v-model.trim="vars.coredns_external_zones[scope.index].nameservers[nameserver.index]" :placeholder="$t('zone_placeholder')"></el-input>
          </template>
        </FieldArray>
        <FieldNumber :holder="vars.coredns_external_zones[scope.index]" fieldName="cache" anti-freeze label-width="90px" style="margin-top: 15px;" :placeholder="$t('cache_placeholder')">
          <template #append>{{$t('second')}}</template>
        </FieldNumber>
      </template>
      <template #help>
        <div style="color: var(--el-text-color-secondary);">
          对于 kubernetes 集群外的域名，优先匹配『按域名后缀指定的集群外 DNS』进行解析，没有在此指定集群外 DNS 的域名后缀，通过『上游DNS服务』进行解析。
        </div>
      </template>
    </FieldArray> -->
    <ConfigSection v-model:enabled="nodeLocalDnsEnabled" :label="$t('nodelocaldnsLabel')" labelWidth="120px">
      <template #more>
        在每个节点上设置一个 DNS 服务，对 DNS 解析结果进行缓存，减少对 CoreDNS 的访问次数。
      </template>
      <FieldString :holder="vars" fieldName="nodelocaldns_ip" :prop="prop" :rules="nodelocaldns_ip_rules" :helpString="$t('nodelocaldnsHelpStr')"></FieldString>
    </ConfigSection>
  </ConfigSection>
</template>

<script>
import { Netmask } from 'netmask'
import ExecuteTask from '../../../common/task/ExecuteTask.vue'

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
      form: {
        verbose: '',
        fork: 5,
        action: 'install_addon',
      },
      pingpong: {},
      pingpong_loading: true,
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
      ],
      nodelocaldns_ip_rules: [
        { 
          validator (rule, value, callback) {
            if (!/^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/.test(value)) {
              return callback('requires an ip address')
            }
            return callback()
          },
          trigger: 'blur'
        }
      ],
      upstreamDnsServerItemRules: [
        {
          validator: this.$validators.ipv4,
          trigger: 'blur',
        }
      ]
    }
  },
  inject: ['editMode', 'isClusterOnline', 'isClusterInstalled'],
  computed: {
    offlineNodes () {
      let result = []
      for (let key in this.cluster.inventory.all.hosts) {
        if (this.pingpong[key] && this.pingpong[key].ping !== 'pong') {
          result.push(key)
        }
      }
      return result
    },
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
    nodeLocalDnsEnabled: {
      get () {
        return this.vars.enable_nodelocaldns
      },
      set (v) {
        this.vars.enable_nodelocaldns = v
      }
    },
    coredns_external_zones: {
      get () { return this.vars.coredns_external_zones },
      set () {},
    }
  },
  components: { ExecuteTask },
  watch: {
    coredns_external_zones: {
      deep: true,
      handler(newValue) {
        this.vars.nodelocaldns_external_zones = newValue
      }
    }
  },
  mounted () {
    if (this.vars.coredns_external_zones == undefined) {
      this.vars.coredns_external_zones = []
    }
  },
  methods: {
    testPingPong (nodes) {
      this.pingpong = {}
      this.pingpong_loading = true
      let req = { nodes: nodes }
      this.kuboardSprayApi.post(`/clusters/${this.cluster.name}/state/ping`, req).then(resp => {
        this.pingpong = resp.data.data.items
        this.pingpong_loading = false
      }).catch(e => {
        if (e.response && e.response.data) {
          this.$message.error('不能测试节点是否在线: ' + e.response.data.message)
        } else {
          this.$message.error('不能测试节点是否在线: ' + e)
        }
        this.pingpong_loading = false
      })
    },
    applyDnsConfig () {
      console.log('abc')
      return new Promise((resolve, reject) => {
        let req = {
          verbose: this.form.verbose,
          fork: this.form.fork,
          addon_name: 'nodelocaldns'
        }
        this.kuboardSprayApi.post(`/clusters/${this.cluster.name}/install_addon`, req).then(resp => {
          let pid = resp.data.data.pid
          resolve(pid)
        }).catch(e => {
          reject(e.response.data.message || e)
        })
      })
    },
    onVisibleChange (v) {
      if (v && this.cluster.state && this.addonState) {
        this.testPingPong('target')
      }
    },
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
