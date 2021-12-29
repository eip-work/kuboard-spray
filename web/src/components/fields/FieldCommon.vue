<i18n>
en:
  please: 'Please input '
  isRequiredField: ' is required.'
zh:
  please: 请输入
  isRequiredField: 为必填字段
</i18n>

<template>
  <el-form-item :rules="computedRules" :prop="prop ? prop + '.' + fieldName : fieldName">
    <template #label>
      {{ compute_label }}
    </template>
    <slot v-if="compute_edit_mode" name="edit"></slot>
    <div v-else class="app_text_mono">
      <slot v-if="holder === undefined" name="view"></slot>
      <span v-else-if="value !== undefined && value !== ''">
        <slot name="view">{{ compute_display_value }}</slot>
      </span>
      <span v-else class="field_placeholder">{{ compute_placeholder }}</span>
    </div>
  </el-form-item>
</template>

<script>
export default {
  props: {
    holder: { type: Object, required: false, undefined },
    fieldName: { type: String, required: false, undefined },
    prop: { type: String, required: false },
    required: { type: Boolean, required: false, default: false },
    rules: { type: Array, required: false, default: () => ([])},
    label: { type: String, required: false, default: undefined },
    placeholder: { type: String, required: false, default: undefined },
    antiFreeze: { type: Boolean, required: false, default: false },
    readOnly: { type: Boolean, required: false, default: false },
  },
  data () {
    return {
    }
  },
  inject: ['editMode'],
  computed: {
    obj: {
      get () {
        return this.holder
      },
      set (v) {
        console.log(v)
      }
    },
    compute_label () {
      if (this.label) {
        return this.label
      }
      let temp = this.$t('field.' + this.fieldName)
      if (temp === 'field.' + this.fieldName) {
        return this.fieldName
      }
      return temp
    },
    compute_edit_mode () {
      if (this.readOnly) {
        return false
      }
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
    compute_display_value () {
      let temp = this.$t('field.' + this.fieldName + '-' + this.value)
      if (temp === 'field.' + this.fieldName + '-' + this.value) {
        return this.value
      } else {
        return temp
      }
    },
    compute_placeholder () {
      if (this.placeholder) {
        return this.placeholder
      }
      let temp = this.$t('field.' + this.fieldName + '_placeholder')
      if (temp == ('field.' + this.fieldName + '_placeholder')) {
        return this.$t('field.please_input') + this.$t('field.' + this.fieldName)
      } else {
        return temp
      }
    },
    computedRules () {
      let result = []
      if (this.required) {
        let message = this.$t('field.' + this.fieldName) + ' ' + this.$t('field.is_required_field')
        result.push({
          required: true,
          validator: (rule, value, callback) => {
            if (value !== undefined && value !== '') {
              return callback()
            }
            return callback(message)
          },
          message: message,
          trigger: 'blur'
        })
      }
      result.push(... this.rules)
      return result
    },
    value: {
      get () {
        if (this.holder) {
          return this.holder[this.fieldName]
        } else {
          return undefined
        }
      },
      set (v) {
        if (v) {
          this.obj[this.fieldName] = v
        } else {
          delete this.obj[this.fieldName]
        }
      }
    }
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">
.field_placeholder {
  color: $--color-text-placeholder;
  font-size: 12px;
}
</style>
