<i18n>
en:
  clear1: This is going to clear 
  clear2: ', do you confirm to continue?'
zh:
  clear1: 此操作将清除
  clear2: 的相关设置，是否继续？

</i18n>

<template>
  <el-form-item label-width="80px" class="config_section_section">
    <template #label>
      <span v-if="!compute_edit_mode" class="viewLabel">{{label}}</span>
      <template v-else>
        <el-checkbox v-model="enabledRef" :disabled="disabled"
          ><span class="enableButton">{{label}}</span>
        </el-checkbox>
        <el-popover placement="bottom-start" width="320" trigger="manual"
          v-model:visible="visiblePopover">
          <p>
            {{ $t('clear1') }}
            <i class="strong-text">{{label}}</i>
            {{ $t('clear2') }}
          </p>
          <div style="text-align: right; margin: 0">
            <el-button size="mini" type="info" plain @click="confirm(false)" icon="el-icon-close">取 消</el-button>
            <el-button type="primary" size="mini" @click="confirm(true)" icon="el-icon-check">确 定</el-button>
          </div>
          <template #reference>
            <div></div>
          </template>
        </el-popover>
      </template>
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
    antiFreeze: { type: Boolean, required: false, default: false },
  },
  data () {
    return {
      expanded: true,
      visiblePopover: false,
    }
  },
  inject: ['editMode'],
  computed: {
    compute_edit_mode () {
      if (this.editMode === 'view') {
        return false
      }
      if (this.antiFreeze) {
        return true
      }
      if (this.editMode === 'frozen') {
        return false
      }
      return true
    },
    enabledRef: {
      get () {
        return this.enabled
      },
      set (v) {
        if (v === false) {
          this.visiblePopover = true
        } else {
          this.expanded = v
          this.$emit('update:enabled', v)
        }
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
    confirm (flag) {
      if (flag) {
        this.expanded = false
        this.$emit('update:enabled', false)
      } else {
        this.expanded = true
      }
      this.visiblePopover = false
    }
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
.viewLabel {
  font-weight: bold;
}
</style>
