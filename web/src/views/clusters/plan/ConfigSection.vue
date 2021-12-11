<i18n>
en:
  hi: kuboard
zh:
  hi: kuboard
</i18n>

<template>
  <el-form-item label-width="80px" class="config_section_section">
    <template #label>
      <el-checkbox v-model="enabledRef" :disabled="disabled">
        <span class="enableButton">{{label}}</span>
      </el-checkbox>
    </template>
    <div>
      <div :class="expandedRef ? 'config_section_header expanded' : 'config_section_header'" @click="expandedRef = !expandedRef">
        <el-icon style="vertical-align: middle;" v-if="expandedRef"><arrow-down-bold /></el-icon>
        <el-icon style="vertical-align: middle;" v-else><arrow-up-bold /></el-icon>
        <span style="margin-left: 5px">
          {{ description || label }}
        </span>
      </div>
      <el-collapse-transition>
        <div v-if="expandedRef" class="config_section_content">
          <slot></slot>
        </div>
      </el-collapse-transition>
    </div>
  </el-form-item>
</template>

<script>
import { ArrowDownBold, ArrowUpBold } from '@element-plus/icons'

export default {
  props: {
    enabled: { type: Boolean, required: true },
    label: { type: String, required: true },
    description: { type: String, required: false, default: undefined },
    disabled: { type: Boolean, required: false, default: false },
  },
  data () {
    return {
      expanded: true
    }
  },
  computed: {
    enabledRef: {
      get () {
        return this.enabled
      },
      set (v) {
        this.expanded = v
        this.$emit('update:enabled', v)
      }
    },
    expandedRef: {
      get () {
        return this.enabledRef && this.expanded
      },
      set (v) {
        this.expanded = v
      }
    }
  },
  components: { ArrowDownBold, ArrowUpBold },
  mounted () {
  },
  methods: {

  }
}
</script>

<style>
.config_section_section .el-form-item__label {
  font-size: 13px;
}
</style>

<style scoped lang="scss">
.enableButton {
  font-size: 13px;
}
.config_section_header {
  background-color: $--color-primary-light-5;
  border: solid 1px $--color-primary-light-5;
  color: white;
  height: 26px;
  padding: 0 15px;
  cursor: pointer;
  font-size: 13px;
}
.config_section_header.expanded {
  background-color: $--color-primary-light-2;
}
.config_section_content {
  padding: 10px 15px 0 15px;
  margin-bottom: 10px;
  background-color: white;
  border: solid 1px $--color-primary-light-9;
  border-top: none;
}
</style>
