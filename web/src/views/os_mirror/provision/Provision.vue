<i18n>
en:
  target_host: Target Host
zh:
  target_host: 目标主机

</i18n>

<template>
  <div>
    <el-form-item :label="$t('target_host')">
      <div class="ssh_content">
        <SshParams :holder="inventory.all.hosts.mirror_node" prop="inventory.all.hosts.mirror_node" :mirrorName="os_mirror.name"></SshParams>
        <NodeFact v-bind="inventory.all.hosts.mirror_node" class="provision_node_fact" ref="fact"
          node_name="mirror_node" node_owner_type="mirror" :node_owner="os_mirror.name"></NodeFact>
      </div>
    </el-form-item>
    <ProvisionUbuntu :os_mirror="os_mirror"></ProvisionUbuntu>
  </div>
</template>

<script>
import SshParams from './SshParams.vue'
import NodeFact from '../../common/fact/NodeFact.vue'
import ProvisionUbuntu from './ProvisionUbuntu.vue'

export default {
  props: {
    os_mirror: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
    inventory: {
      get () { return this.os_mirror.inventory },
      set () {},
    },
    status: {
      get () { return this.os_mirror.status },
      set () {},
    }
  },
  components: { SshParams, NodeFact, ProvisionUbuntu },
  watch: {
    'inventory.all.hosts.mirror_node.ansible_host': function() {
      this.updateUrl()
    },
    // 'inventory.all.children.target.vars.apache2_default_port': function () {
    //   this.updateUrl()
    // }
  },
  mounted () {
    if (this.inventory.all.hosts.mirror_node.ansible_host) {
      this.$refs.fact.loadFacts()
    }
  },
  methods: {
    updateUrl () {
      let temp =  'http://' + this.inventory.all.hosts.mirror_node.ansible_host
      // if (this.inventory.all.children.target.vars.apache2_default_port && this.inventory.all.children.target.vars.apache2_default_port !== 80) {
      //   temp += ':' + this.inventory.all.children.target.vars.apache2_default_port
      // }
      temp += '/ubuntu'
      this.status.url = temp
    }
  }
}
</script>

<style lang="css">
.provision_node_fact .package_title {
  padding: 0 20px;
}
.provision_node_fact .package_info {
  padding: 0 20px;
}
</style>

<style scoped lang="css">
.ssh_content {
  padding: 20px;
  background-color: var(--el-color-primary-light-9);
}
</style>
