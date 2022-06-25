<i18n>
en:
  
zh:
  ipip_mode: ipip_mode
  vxlan_mode: vxlan_mode
  backend: backend
  calico_mode: calico_mode
</i18n>


<template>
  <div>
    <FieldCommon :holder="vars" :prop="prop" fieldName="calico_mode" :placeholder="$t('calico_mode')" label-width="160px">
      <template #view>{{ calico_mode }}</template>
      <template #edit>
        <el-radio-group v-model="calico_mode">
          <el-radio-button label="IPIP">IPIP</el-radio-button>
          <el-radio-button label="VXLAN">VXLAN</el-radio-button>
        </el-radio-group>
      </template>
    </FieldCommon>
    <FieldCommon v-if="calico_mode === 'IPIP'" :holder="vars" :prop="prop" fieldName="calico_ipip_mode" :placeholder="$t('ipip_mode')" label-width="160px">
      <template #view>
        <span>{{ ipip_mode }}</span>
      </template>
      <template #edit>
        <el-radio-group v-model="ipip_mode">
          <el-radio label="Always">Always</el-radio>
          <el-radio label="CrossSubnet">CrossSubnet</el-radio>
        </el-radio-group>
      </template>
    </FieldCommon>
    <FieldCommon v-if="calico_mode === 'VXLAN'" :holder="vars" :prop="prop" fieldName="calico_vxlan_mode" :placeholder="$t('vxlan_mode')" label-width="160px">
      <template #view>
        <span>{{ vxlan_mode }}</span>
      </template>
      <template #edit>
        <el-radio-group v-model="vxlan_mode">
          <el-radio label="Always">Always</el-radio>
          <el-radio label="CrossSubnet">CrossSubnet</el-radio>
        </el-radio-group>
      </template>
    </FieldCommon>
    <FieldCommon :holder="vars" :prop="prop" fieldName="calico_network_backend" :placeholder="$t('backend')" label-width="160px">
      <template #view>
        <span>{{ calico_network_backend }}</span>
      </template>
      <template #edit>
        <el-radio-group v-model="calico_network_backend">
          <el-radio label="bird">bird</el-radio>
          <el-radio label="vxlan">vxlan</el-radio>
          <el-radio label="none">none</el-radio>
        </el-radio-group>
      </template>
    </FieldCommon>
  </div>
</template>

<script>
import compareVersions from 'compare-versions'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
    vars: {
      get () { return this.cluster.inventory.all.children.target.children.k8s_cluster.vars },
      set () {},
    },
    prop() { return 'all.children.target.children.k8s_cluster.vars' },
    calico_mode: {
      get () {
        if (this.ipip_mode === 'Never') {
          return 'VXLAN'
        }
        if (this.vxlan_mode === 'Never') {
          return 'IPIP'
        }
        return ''
      },
      set (v) {
        if (v === 'IPIP') {
          this.ipip_mode = 'Always'
          this.vxlan_mode = 'Never'
        } else if (v === 'VXLAN') {
          this.vxlan_mode = 'Always'
          this.ipip_mode = 'Never'
        }
      }
    },
    ipip_mode: {
      get () {
        if (this.vars.calico_ipip_mode) {
          return this.vars.calico_ipip_mode
        } else {
          if (compareVersions(this.cluster.resourcePackage.data.kubernetes.kube_version, 'v1.24.0') < 0) {
            return 'Always'
          } else {
            return 'Never'
          }
        }
      },
      set (v) {
        this.vars.calico_ipip_mode = v
      },
    },
    vxlan_mode: {
      get () {
        if (this.vars.calico_vxlan_mode) {
          return this.vars.calico_vxlan_mode
        } else {
          if (compareVersions(this.cluster.resourcePackage.data.kubernetes.kube_version, 'v1.24.0') < 0) {
            return 'Never'
          } else {
            return 'Always'
          }
        }
      },
      set (v) {
        this.vars.calico_vxlan_mode = v
      },
    },
    calico_network_backend: {
      get () {
        if (this.vars.calico_network_backend) {
          return this.vars.calico_network_backend
        } else {
          return 'bird'
        }
      },
      set (v) {
        this.vars.calico_network_backend = v
      }
    }
  },
  components: { },
  mounted () {
  },
  methods: {
    compareVersions,
  }
}
</script>

<style scoped lang="scss">

</style>
