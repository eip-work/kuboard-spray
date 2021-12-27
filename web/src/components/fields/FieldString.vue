<template>
  <FieldCommon :fieldName="fieldName" :holder="holder" :prop="prop" :rules="rules" :required="required" :label="label" :placeholder="placeholder">
    <template #edit>
      <el-input v-model.trim="value" :show-password="showPassword" :disabled="disabled"
        :placeholder="compute_placeholder"></el-input>
    </template>
    <template #view>
      {{ showPassword && hidePassword ? value.replace(/.?/g, '*') : value }}
      <el-button v-if="showPassword" type="text" icon="el-icon-view" @click="hidePassword = !hidePassword"></el-button>
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
    showPassword: { type: Boolean, required: false, default: false },
    placeholder: { type: String, required: false, default: undefined },
    label: { type: String, required: false, default: undefined },
    disabled: { type: Boolean, required: false, default: false },
  },
  data () {
    return {
      hidePassword: true,
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
</style>
