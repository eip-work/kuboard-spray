<template>
  <el-form-item :rules="computedRules" :prop="prop ? prop + '.' + fieldName : fieldName">
    <template #label>
      {{ $t('field.' + fieldName) }}
    </template>
    <el-switch v-model.number="obj[fieldName]" :disabled="editMode === 'view'"></el-switch>
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

  }
}
</script>

<style scoped lang="scss">

</style>
