<template>
  <FieldCommon :fieldName="fieldName" :holder="holder" :prop="prop" :rules="rules" :required="required" :label="label" :placeholder="placeholder">
    <template #edit>
      <div style="display: flex;">
        <el-select v-model.trim="obj[fieldName]" style="flex-grow: 1;" :clearable="clearable" :disabled="disabled"
          :placeholder="compute_placeholder" @visible-change="load($event)">
          <el-option v-for="(item, index) in options" :key="'i' + index" :value="item.value" :label="item.label" :disabled="item.disabled">
            {{item.label}}
          </el-option>
        </el-select>
        <slot></slot>
      </div>
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
    loadOptions: { type: Function, required: true },
    disabled: { type: Boolean, required: false, default: false },
    placeholder: { type: String, required: false, default: undefined },
    label: { type: String, required: false, default: undefined },
    clearable: { type: Boolean, required: false, default: false },
  },
  data () {
    return {
      options: []
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
