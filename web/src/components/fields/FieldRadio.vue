<template>
  <FieldCommon :fieldName="fieldName" :holder="holder" :prop="prop" :rules="rules" :required="required" :label="label" :placeholder="placeholder">
    <template #edit>
      <el-radio-group v-model="value" :disabled="disabled">
        <template v-for="(option, index) in options" :key="index">
          <el-radio-button :label="option">
            {{ optionDesc(option) }}
          </el-radio-button>
        </template>
      </el-radio-group>
    </template>
    <template #view>
      <template v-for="(option, index) in options" :key="index">
        <el-button v-if="value === option" :label="option" type="primary" disabled>
          {{ optionDesc(option) }}
        </el-button>
      </template>
    </template>
  </FieldCommon>
</template>

<script>
import compute_placeholder_mixin from './compute_placeholder_mixin.js'

export default {
  props: {
    holder: { type: Object, required: true },
    fieldName: { type: String, required: true },
    prop: { type: String, required: false },
    required: { type: Boolean, required: false, default: false },
    rules: { type: Array, required: false, default: () => ([])},
    options: { type: Array, required: true },
    label: { type: String, required: false, default: undefined },
    disabled: { type: Boolean, required: false, default: false },
    placeholder: { type: String, required: false, default: undefined },
  },
  data () {
    return {

    }
  },
  mixins: [ compute_placeholder_mixin ],
  computed: {
    obj: {
      get () {
        return this.holder
      },
      set (v) {
        console.log(v)
      }
    },
    value: {
      get () {
        return this.holder[this.fieldName]
      },
      set (v) {
        this.obj[this.fieldName] = v
      }
    },
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
