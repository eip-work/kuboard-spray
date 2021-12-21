<template>
  <el-form-item :rules="computedRules" :prop="prop ? prop + '.' + fieldName : fieldName">
    <template #label>
      {{ $t('field.' + fieldName) }}
    </template>
    <el-radio-group v-model="obj[fieldName]">
      <template v-for="(option, index) in options" :key="index">
        <el-radio-button v-if="editMode !== 'view' || obj[fieldName] === option" :label="option">
          {{ optionDesc(option) }}
        </el-radio-button>
      </template>
    </el-radio-group>
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
    options: { type: Array, required: true },
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
    optionDesc (option) {
      let temp = this.$t('field.' + this.fieldName + '-' + option)
      if (temp === 'field.' + this.fieldName + '-' + option) {
        return option
      } else {
        return temp
      }
    },
  }
}
</script>

<style scoped lang="scss">

</style>
