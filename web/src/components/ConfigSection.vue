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
            <el-button size="default" type="info" plain @click="confirm(false)" icon="el-icon-close">取 消</el-button>
            <el-button type="primary" size="default" @click="confirm(true)" icon="el-icon-check">确 定</el-button>
          </div>
          <template #reference>
            <div></div>
          </template>
        </el-popover>
      </template>
    </template>
    <div style="flex-grow: 1;">
      <div :class="expandedRef ? 'config_section_header expanded noselect' : 'config_section_header noselect'" @click="expandedRef = !expandedRef" style="display: flex;">
        <div>
          <el-icon style="vertical-align: middle;" v-if="expandedRef"><el-icon-arrow-down-bold /></el-icon>
          <el-icon style="vertical-align: middle;" v-else><el-icon-arrow-up-bold /></el-icon>
        </div>
        <div style="flex-grow: 1; margin-left: 5px">
          <slot name="header">
            <span>
              {{ description || label }}
            </span>
          </slot>
        </div>
        <div v-if="helpLink">
          <el-link :href="helpLink" target="blank" type="warning" icon="el-icon-link" class="app_white_link" @click.stop>
            {{$t('msg.help')}}
          </el-link>
        </div>
      </div>
      <el-collapse-transition>
        <div v-if="expandedRef">
          <div v-if="$slots.more" class="config_section_more">
            <slot name="more"></slot>
          </div>
          <div class="config_section_content">
            <slot></slot>
          </div>
        </div>
      </el-collapse-transition>
    </div>
  </el-form-item>
</template>

<script>
export default {
  props: {
    enabled: { type: Boolean, required: true },
    label: { type: String, required: true },
    description: { type: String, required: false, default: undefined },
    disabled: { type: Boolean, required: false, default: false },
    antiFreeze: { type: Boolean, required: false, default: false },
    helpLink: { type: String, required: false, default: undefined },
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
  components: {},
  mounted () {
  },
  methods: {
    expand (flag) {
      this.expanded = flag
    },
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

<style scoped lang="css">
.enableButton {
  font-size: 13px;
}
.config_section_header {
  background-color: var(--el-color-primary-light-5);
  border: solid 1px var(--el-color-primary-light-5);
  color: white;
  height: 26px;
  padding: 0 15px;
  cursor: pointer;
  font-size: 13px;
}
.config_section_header.expanded {
  background-color: var(--el-color-primary-light-2);
}
.config_section_content {
  padding: 10px 15px 0 15px;
  margin-bottom: 10px;
  background-color: white;
  border: solid 1px var(--el-color-primary-light-9);
  border-top: none;
}
.viewLabel {
  font-weight: bold;
}
.config_section_more {
  font-size: 13px;
  padding: 2px 15px;
  border-left: solid 1px var(--el-color-primary-light-9);
  border-right: solid 1px var(--el-color-primary-light-9);
  background-color: var(--el-background-color-base);
  color: var(--el-color-info);
}
</style>
