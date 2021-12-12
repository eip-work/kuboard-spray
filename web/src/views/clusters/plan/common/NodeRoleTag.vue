<template>
  <div :class="(enabled ? 'role noselect selected ' : 'role noselect ') + role">
    <el-icon style="vertical-align: middle; margin-top: -3px;">
      <circle-check v-if="enabled"/>
      <circle-close v-else/>
    </el-icon>
    {{$t('node.' + role)}}
  </div>
</template>

<script>
import { CircleCheck, CircleClose } from '@element-plus/icons'

export default {
  props: {
    enabled: { type: Boolean, required: true },
    role: { type: String, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { CircleCheck, CircleClose },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">
.role {
  color: white;
  margin-right: 10px;
  padding: 0 10px;
  width: 100px;
  border-radius: 5px;
  cursor: pointer;
  color: $--color-text-secondary;
}
.role:hover {
  opacity: 0.8;
}
.selected {
  color: white;
  font-weight: bold;
}
.etcd {
  background-color: $--color-warning-lighter;
}
.etcd.selected {
  background-color: $--color-warning;
}
.kube_control_plane {
  background-color: $--color-primary-light-9;
}
.kube_control_plane.selected {
  background-color: $--color-primary;
}
.kube_node {
  background-color: $--color-success-lighter;
}
.kube_node.selected {
  background-color: $--color-success;
}
</style>
