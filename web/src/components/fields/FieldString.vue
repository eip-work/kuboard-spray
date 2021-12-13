<template>
  <el-form-item :rules="computedRules" :prop="prop ? prop + '.' + fieldName : undefined">
    <template #label>
      {{ $t('field.' + fieldName) }}
    </template>
    <el-input v-if="editMode !== 'view'" v-model.trim="obj[fieldName]" :show-password="showPassword"
      :placeholder="placeholder || $t('field.' + fieldName + '_placeholder')"></el-input>
    <div v-else class="app_text_mono">{{ obj[fieldName] }}</div>
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
