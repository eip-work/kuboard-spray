<template>
  <div :class="(enabled ? 'role noselect selected ' : 'role noselect ') + role" @click="click" :style="readOnly ? 'cursor: inherit;' : ''">
    <el-icon style="vertical-align: middle; margin-top: -3px; margin-right: 5px;">
      <circle-check v-if="enabled"/>
      <circle-close v-else/>
    </el-icon>
    <slot>
      {{$t('node.' + role)}}
    </slot>
  </div>
</template>

<script>
import { CircleCheck, CircleClose } from '@element-plus/icons'

export default {
  props: {
    enabled: { type: Boolean, required: true },
    role: { type: String, required: true },
    readOnly: { type: Boolean, required: false, default: false },
  },
  data() {
    return {

    }
  },
  inject: ['editMode'],
  computed: {
  },
  components: { CircleCheck, CircleClose },
  mounted () {
  },
  methods: {
    click () {
      if (this.editMode !== 'view') {
        this.$emit('clickTag')
      }
    }
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
  text-align: center;
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
