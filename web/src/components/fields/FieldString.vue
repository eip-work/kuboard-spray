<i18n>
en:
  please: 'Please input '
  isRequiredField: ' is required.'
zh:
  please: 请输入
  isRequiredField: 为必填字段
</i18n>

<template>
  <el-form-item :rules="computedRules" :prop="prop ? prop + '.' + fieldName : undefined">
    <template #label>
      {{ label || $t('field.' + fieldName) }}
    </template>
    <el-input v-if="editMode !== 'view'" v-model.trim="value" :show-password="showPassword" :disabled="disabled"
      :placeholder="compute_placeholder"></el-input>
    <div v-else class="app_text_mono">
      <span v-if="value">{{ value }}</span>
      <span v-else class="field_placeholder">{{ compute_placeholder }}</span>
    </div>
  </el-form-item>
</template>

<script>

export default {
  props: {
    holder: { type: Object, required: true },
    fieldName: { type: String, required: true },
    prop: { type: String, required: false },
    required: { type: Boolean, required: false, default: false },
    rules: { type: Array, required: false, default: () => ([])},
    showPassword: { type: Boolean, required: false, default: false },
    placeholder: { type: String, required: false, default: undefined },
    label: { type: String, required: false, default: undefined },
    disabled: { type: Boolean, required: false, default: false },
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
    compute_placeholder () {
      if (this.placeholder) {
        return this.placeholder
      }
      let temp = this.$t('field.' + this.fieldName + '_placeholder')
      if (temp == ('field.' + this.fieldName + '_placeholder')) {
        return this.$t('please') + this.$t('field.' + this.fieldName)
      } else {
        return temp
      }
    },
    computedRules () {
      let result = []
      if (this.required) {
        let message = this.$t('field.' + this.fieldName) + ' ' + this.$t('isRequiredField')
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
        return this.holder[this.fieldName]
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
