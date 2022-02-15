<template>
  <FieldCommon :fieldName="fieldName" :holder="holder" :prop="prop" :rules="rules" :required="required" :label="label" :placeholder="placeholder">
    <template #edit>
      <div style="display: flex; flex-grow: 1;">
        <el-select v-model="value" style="flex-grow: 1;" :clearable="clearable" :disabled="disabled" :allowCreate="allowCreate" :filterable="filterable"
          :placeholder="compute_placeholder" @visible-change="load($event)" :loading="loading" :multiple="multiple">
          <el-option v-for="(item, index) in options" :key="'i' + index" :value="item.value" :label="item.label" :disabled="item.disabled">
            {{item.label}}
          </el-option>
        </el-select>
        <slot name="edit"></slot>
      </div>
    </template>
    <template #view>
      <span v-if="value !== undefined && value !== ''">
        {{ compute_display_value }}
      </span>
      <span v-else class="field_placeholder">{{ compute_placeholder }}</span>
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
    allowCreate: { type: Boolean, required: false, default: false },
    filterable: { type: Boolean, required: false, default: false },
    multiple: { type: Boolean, required: false, default: false },
  },
  data () {
    return {
      options: [],
      loading: false,
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
      get () { return this.obj[this.fieldName] },
      set (v) {
        if (v) {
          this.obj[this.fieldName] = v
        } else {
          delete this.obj[this.fieldName]
        }
      }
    },
    compute_display_value () {
      let temp = ''
      for (let i in this.options) {
        if (this.options[i].value === this.value) {
          return this.options[i].label
        }
      }
      temp = this.$t('field.' + this.fieldName + '-' + this.value)
      if (temp === 'field.' + this.fieldName + '-' + this.value) {
        return this.value
      } else {
        return temp
      }
    },
  },
  components: { },
  mounted () {
    if (this.loadOptions) {
      this.load(true)
    }
  },
  methods: {
    load (flag) {
      if (flag) {
        this.loading = true
        this.loadOptions.call(this.$parent, this.fieldName).then(ops => {
          this.options = ops
          this.loading = false
        }).catch(e => {
          console.log(e)
          this.loading = false
        })
      }
    }
  }
}
</script>

<style scoped lang="css">
.field_placeholder {
  color: var(--el-text-color-placeholder);
  font-size: 12px;
}
</style>
