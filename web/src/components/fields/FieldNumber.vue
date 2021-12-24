<template>
  <el-form-item :rules="computedRules" :prop="prop ? prop + '.' + fieldName : fieldName">
    <template #label>
      {{ $t('field.' + fieldName) }}
    </template>
    <el-input v-if="editMode !== 'view'" v-model.number="obj[fieldName]" :placeholder="compute_placeholder">
      <template #append><slot name="append"></slot></template>
    </el-input>
    <div v-else class="app_text_mono">
      <span v-if="obj[fieldName]">{{ obj[fieldName] }}</span>
      <span v-else>{{ compute_placeholder }}</span>
    </div>
    <slot></slot>
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
    placeholder: { type: String, required: false, default: undefined },
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
          trigger: 'blur'
        })
      }
      result.push(... this.rules)
      return result
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

</style>
