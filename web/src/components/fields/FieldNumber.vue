<template>
  <FieldCommon :fieldName="fieldName" :holder="holder" :prop="prop" :rules="rules" :required="required" :label="label" :placeholder="placeholder">
    <template #edit>
      <el-input v-model.number="value" :placeholder="compute_placeholder">
        <template #append v-if="$slots.append"><slot name="append"></slot></template>
      </el-input>
      <div class="desc" v-if="$slots.edit_desc">
        <slot name="edit_desc"></slot>
      </div>
    </template>
    <template #view>
      {{ value }} <span class="app_text_mono" style="font-size: 13px;"><slot name="append"></slot></span>
      <div class="desc" v-if="$slots.view_desc">
        <slot name="view_desc"></slot>
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
    placeholder: { type: String, required: false, default: undefined },
    label: { type: String, required: false, default: undefined },
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
      get () { return this.holder[this.fieldName] },
      set (v) {
        if (v != '') {
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

<style scoped lang="css">
.field_placeholder {
  color: var(--el-text-color-placeholder);
  font-size: 12px;
}
.desc {
  color: var(--el-text-color-secondary);
  display: inline-block;
  vertical-align: middle;
  margin-left: 10px;
}
</style>
