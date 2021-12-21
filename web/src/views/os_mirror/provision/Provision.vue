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
        <NodeFact v-bind="inventory.all.hosts.mirror_node" class="provision_node_fact"
          node_name="mirror_node" :node_owner_type="mirror" :node_owner="os_mirror.name"></NodeFact>
      </div>
    </el-form-item>
    <FieldString :holder="inventory.all.children.target.vars" prop="inventory.all.children.target.vars" fieldName="apt_mirror_apache_root_dir" required></FieldString>
  </div>
</template>

<script>
import SshParams from './SshParams.vue'
import NodeFact from '../../fact/NodeFact.vue'

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
    }
  },
  components: { SshParams, NodeFact },
  mounted () {
  },
  methods: {

  }
}
</script>

<style lang="scss">
.provision_node_fact .package_title {
  padding: 0 20px;
}
.provision_node_fact .package_info {
  padding: 0 20px;
}
</style>

<style scoped lang="scss">
.ssh_content {
  padding: 20px;
  background-color: $--color-primary-light-9;
}
</style>
