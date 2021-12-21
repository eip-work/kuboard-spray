<i18n>
en:
  please: 'Please select '
  isRequiredField: ' is required.'
zh:
  please: 请选择
  isRequiredField: 为必填字段
</i18n>

<template>
  <el-form-item :rules="computedRules" :prop="prop ? prop + '.' + fieldName : undefined">
    <template #label>
      {{ label || $t('field.' + fieldName) }}
    </template>
    <div style="display: flex;" v-if="editMode !== 'view'">
      <el-select v-model.trim="obj[fieldName]" style="flex-grow: 1;" clearable :disabled="disabled"
        :placeholder="compute_placeholder" @visible-change="load($event)">
        <el-option v-for="(item, index) in options" :key="'i' + index" :value="item.value" :label="item.label">
          {{item.label}}
        </el-option>
      </el-select>
      <slot></slot>
    </div>
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
    loadOptions: { type: Function, required: true },
    disabled: { type: Boolean, required: false, default: false },
    placeholder: { type: String, required: false, default: undefined },
    label: { type: String, required: false, default: undefined },
  },
  data () {
    return {
      options: []
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
      console.log(temp, temp == ('field.' + this.fieldName + '_placeholder'))
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
          message: message,
          validator: (rule, value, callback) => {
            if (value !== undefined && value !== '') {
              return callback()
            }
            return callback(message)
          },
          trigger: 'change'
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
    this.load(true)
  },
  methods: {
    load (flag) {
      if (flag) {
        this.loadOptions.call(this.$parent).then(ops => {
          this.options = ops
        }).catch(e => {
          console.log(e)
        })
      }
    }
  }
}
</script>

<style scoped lang="scss">
.field_placeholder {
  color: $--color-text-placeholder;
  font-size: 12px;
}
</style>
